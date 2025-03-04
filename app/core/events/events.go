package events

import (
	"context"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type EventManager struct {
	ctx context.Context
	mu  sync.Mutex
}

func NewEventManager() *EventManager {
	return &EventManager{}
}

func (e *EventManager) SetContext(ctx context.Context) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.ctx = ctx
}

func (e *EventManager) Emit(eventName string, data ...interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.ctx != nil {
		runtime.EventsEmit(e.ctx, eventName, data...)
	}
}
