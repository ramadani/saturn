package saturn

import (
	"context"
	"errors"
	"strings"
)

type eventListenersMap map[string][]Listener

type localEmitter struct {
	event Event
}

type localEvent struct {
	evLst eventListenersMap
}

type emitResult struct {
	Index int
	Err   error
}

func (e *localEmitter) Emit(ctx context.Context, dispatchable Dispatchable) (err error) {
	header := dispatchable.Header()
	listeners, err := e.event.Listeners(ctx, header)
	if err != nil {
		return
	}

	if len(listeners) == 0 {
		return
	}

	value, err := dispatchable.Body()
	if err != nil {
		return
	}

	total := len(listeners)
	messages := make([]string, total)
	resChan := make(chan emitResult)

	for i, lis := range listeners {
		go func(index int, listener Listener) {
			er := listener.Handle(ctx, value)
			resChan <- emitResult{
				Index: index,
				Err:   er,
			}
		}(i, lis)
	}

	for i := 0; i < total; i++ {
		res := <-resChan
		if res.Err != nil {
			messages[res.Index] = res.Err.Error()
		} else {
			messages[res.Index] = ""
		}
	}

	res := make([]string, 0)
	for _, msg := range messages {
		if msg != "" {
			res = append(res, msg)
		}
	}

	if msg := strings.Join(res, "; "); msg != "" {
		err = errors.New(msg)
	}
	return
}

func (e *localEvent) On(_ context.Context, header string, listeners []Listener) (err error) {
	if lts, ok := e.evLst[header]; ok {
		e.evLst[header] = append(lts, listeners...)
	} else {
		e.evLst[header] = listeners
	}
	return
}

func (e *localEvent) Listeners(_ context.Context, header string) (res []Listener, err error) {
	res = make([]Listener, 0)

	if lts, ok := e.evLst[header]; ok {
		res = lts
	}
	return
}

func NewLocalEmitter(ev Event) Emitter {
	return &localEmitter{event: ev}
}

func NewLocalEvent(evLts eventListenersMap) Event {
	return &localEvent{evLst: evLts}
}

func DefaultLocalEvent() Event {
	return NewLocalEvent(make(eventListenersMap))
}
