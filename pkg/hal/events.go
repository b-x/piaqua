package hal

type Event interface{}

type EventSource interface {
	Loop(quit <-chan struct{}, events chan<- Event)
}

type ButtonPressed struct {
	ID int
}

type TemperatureRead struct {
	ID   int
	Temp int
}

type TemperatureError struct {
	ID    int
	Error error
}
