# saturn

Event Listener for Go

### Install

```bash
go get github.com/ramadani/saturn
```

### Usage

Create event

```go
type user struct {
	Name, Email string
}

type userRegistered struct {
	data user
}

func (e *userRegistered) Header() string {
	return "userRegistered"
}

func (e *userRegistered) Body() ([]byte, error) {
	b, err := json.Marshal(&e.data)
	if err != nil {
		return nil, err
	}
	return b, nil
}
```

Create the listeners

```go
type sendEmailVerificationListener struct{}

func (l *sendEmailVerificationListener) Handle(_ context.Context, value []byte) (err error) {
	data := &user{}
	_ = json.Unmarshal(value, data)

	time.Sleep(time.Second * 3)

	log.Println("send email to", data.Email)
	return
}

type logNewUserListener struct{}

func (l *logNewUserListener) Handle(_ context.Context, value []byte) (err error) {
	data := &user{}
	_ = json.Unmarshal(value, data)

	time.Sleep(time.Second * 1)

	log.Println("log new user", data.Name)
	return
}
```

Setup event emitter and listen the events

```go
event := saturn.NewLocalEvent(map[string][]saturn.Listener{
    "userRegistered": {
        &sendEmailVerificationListener{},
    },
})

emitter := saturn.NewLocalEmitter(event)
ctx := context.Background()

// add another listeners
_ = event.On(ctx, "userRegistered", []saturn.Listener{
    &logNewUserListener{},
})
```

Send the event

```go
_ = emitter.Emit(ctx, &userRegistered{data: user{Name: "Ramadani", Email: "dani@gmail.id"}})
```