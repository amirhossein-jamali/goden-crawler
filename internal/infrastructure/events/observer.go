// File: internal/infrastructure/events/observer.go

package events

import (
	"sync"

	"github.com/amirhossein-jamali/goden-crawler/pkg/logger"
)

// EventType represents the type of an event
type EventType string

// Event types
const (
	// WordFetchStarted is triggered when a word fetch starts
	WordFetchStarted EventType = "word_fetch_started"
	// WordFetchCompleted is triggered when a word fetch completes
	WordFetchCompleted EventType = "word_fetch_completed"
	// WordFetchFailed is triggered when a word fetch fails
	WordFetchFailed     EventType = "word_fetch_failed"
	ExtractionStarted   EventType = "extraction_started"
	ExtractionCompleted EventType = "extraction_completed"
	ExtractionFailed    EventType = "extraction_failed"
)

// Event represents an event in the system
type Event struct {
	Type    EventType
	Payload interface{}
}

// Observer defines the interface for event observers
type Observer interface {
	// OnEvent is called when an event occurs
	OnEvent(event Event)

	// GetID returns the ID of the observer
	GetID() string
}

// Subject defines the interface for event subjects
type Subject interface {
	// RegisterObserver registers an observer
	RegisterObserver(observer Observer)

	// UnregisterObserver unregisters an observer
	UnregisterObserver(observerID string)

	// NotifyObservers notifies all observers of an event
	NotifyObservers(event Event)
}

// EventManager manages event observers and notifications
type EventManager struct {
	observers map[string]Observer
	mutex     sync.RWMutex
}

// NewEventManager creates a new event manager
func NewEventManager() *EventManager {
	return &EventManager{
		observers: make(map[string]Observer),
	}
}

// RegisterObserver registers an observer
func (m *EventManager) RegisterObserver(observer Observer) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.observers[observer.GetID()] = observer
	logger.Debug("Registered observer", logger.F("id", observer.GetID()))
}

// UnregisterObserver unregisters an observer
func (m *EventManager) UnregisterObserver(observerID string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.observers, observerID)
	logger.Debug("Unregistered observer", logger.F("id", observerID))
}

// NotifyObservers notifies all observers of an event
func (m *EventManager) NotifyObservers(event Event) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	logger.Debug("Notifying observers of event", logger.F("type", string(event.Type)))
	for _, observer := range m.observers {
		go observer.OnEvent(event)
	}
}

// Global event manager
var globalEventManager *EventManager
var once sync.Once

// GetEventManager returns the global event manager
func GetEventManager() *EventManager {
	once.Do(func() {
		globalEventManager = NewEventManager()
	})
	return globalEventManager
}

// Global event manager functions

// RegisterObserver registers an observer with the global event manager
func RegisterObserver(observer Observer) {
	GetEventManager().RegisterObserver(observer)
}

// UnregisterObserver unregisters an observer from the global event manager
func UnregisterObserver(observerID string) {
	GetEventManager().UnregisterObserver(observerID)
}

// NotifyObservers notifies all observers of an event using the global event manager
func NotifyObservers(event Event) {
	GetEventManager().NotifyObservers(event)
}

// BaseObserver provides a base implementation of the Observer interface
type BaseObserver struct {
	ID string
}

// NewBaseObserver creates a new base observer
func NewBaseObserver(id string) *BaseObserver {
	return &BaseObserver{
		ID: id,
	}
}

// GetID returns the ID of the observer
func (o *BaseObserver) GetID() string {
	return o.ID
}

// OnEvent is called when an event occurs
func (o *BaseObserver) OnEvent(event Event) {
	// Base implementation does nothing
}
