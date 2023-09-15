package events

import (
	"errors"
	"sync"
)

var ErrHandlerAlreadyRegistered = errors.New("handler already registered")

type EventDispatcher struct {
	handles map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handles: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handles[eventName]; ok {
		for _, h := range ed.handles[eventName] {
			if h == handler {
				return ErrHandlerAlreadyRegistered
			}
		}
	}
	ed.handles[eventName] = append(ed.handles[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Dispatch(event EventInterface) error {
	if handles, ok := ed.handles[event.GetName()]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handles {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, ok := ed.handles[eventName]; ok {
		for _, h := range ed.handles[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, ok := ed.handles[eventName]; ok {
		for i, h := range ed.handles[eventName] {
			if h == handler {
				ed.handles[eventName] = append(ed.handles[eventName][:i], ed.handles[eventName][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (ed *EventDispatcher) Clear() error {
	ed.handles = make(map[string][]EventHandlerInterface)
	return nil
}
