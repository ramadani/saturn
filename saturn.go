package saturn

import (
	"context"
)

type Emitter interface {
	Emit(ctx context.Context, dispatchable Dispatchable) error
}

type Event interface {
	On(ctx context.Context, header string, listeners []Listener) error
	Listeners(ctx context.Context, header string) ([]Listener, error)
}

type EventListener interface {
	Event
	Listen(ctx context.Context) error
}

type Dispatchable interface {
	Header() string
	Body() ([]byte, error)
}

type Listener interface {
	Handle(ctx context.Context, value []byte) error
}
