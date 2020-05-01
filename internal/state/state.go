package state

import "time"

var GlobalState State

func MakeGlobalState() *State {
	return &State{
		Subscribers: make(map[string]Subscriber),
	}
}

type State struct {
	Subscribers map[string]Subscriber
}

type Subscriber struct {
	Number   string
	LastSent time.Time
}

func (s *State) AddSubscriber(number string) {
	s.Subscribers[number] = Subscriber{Number: number}
}
func (s *State) RemoveSubscriber(number string) {
	delete(s.Subscribers, number)
}
