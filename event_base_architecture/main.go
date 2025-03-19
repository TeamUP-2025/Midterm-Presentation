package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// Message represents a chat message
type Message struct {
	ID        string
	Sender    string
	Content   string
	Timestamp time.Time
}

// Event represents an event in the system
type Event struct {
	Type    string
	Payload Message
}

// EventBus handles event distribution
type EventBus struct {
	subscribers map[string][]chan Event
	mutex       sync.RWMutex
}

// NewEventBus creates a new event bus
func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]chan Event),
	}
}

// Subscribe registers a subscriber for a specific event type
func (eb *EventBus) Subscribe(eventType string, ch chan Event) {
	eb.mutex.Lock()
	defer eb.mutex.Unlock()
	eb.subscribers[eventType] = append(eb.subscribers[eventType], ch)
}

// Publish sends an event to all subscribers of that event type
func (eb *EventBus) Publish(event Event) {
	eb.mutex.RLock()
	defer eb.mutex.RUnlock()
	subscribers, exists := eb.subscribers[event.Type]
	if !exists {
		return
	}

	for _, ch := range subscribers {
		go func(c chan Event) {
			c <- event
		}(ch)
	}
}

// Client component that sends messages
type Client struct {
	ID       string
	eventBus *EventBus
}

// SendMessage sends a message from the client
func (c *Client) SendMessage(content string) {
	msg := Message{
		ID:        fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		Sender:    c.ID,
		Content:   content,
		Timestamp: time.Now(),
	}

	c.eventBus.Publish(Event{
		Type:    "messageSent",
		Payload: msg,
	})
}

// MessageReceiver handles incoming messages
type MessageReceiver struct {
	eventBus *EventBus
}

// Start begins listening for incoming messages
func (mr *MessageReceiver) Start() {
	ch := make(chan Event)
	mr.eventBus.Subscribe("messageSent", ch)

	go func() {
		for event := range ch {
			log.Printf("Message received: %s from %s",
				event.Payload.Content, event.Payload.Sender)

			// Create message event for other components
			mr.eventBus.Publish(Event{
				Type:    "messageCreate",
				Payload: event.Payload,
			})
		}
	}()
}

// MessageSaver saves messages to storage
type MessageSaver struct {
	eventBus *EventBus
	messages []Message
	mutex    sync.RWMutex
}

// Start begins listening for messages to save
func (ms *MessageSaver) Start() {
	ch := make(chan Event)
	ms.eventBus.Subscribe("messageCreate", ch)

	go func() {
		for event := range ch {
			ms.mutex.Lock()
			ms.messages = append(ms.messages, event.Payload)
			ms.mutex.Unlock()
			log.Printf("Message saved: %s", event.Payload.ID)
		}
	}()
}

// GetMessages retrieves saved messages
func (ms *MessageSaver) GetMessages() []Message {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()
	result := make([]Message, len(ms.messages))
	copy(result, ms.messages)
	return result
}

// MessagePublisher publishes messages to subscribers
type MessagePublisher struct {
	eventBus *EventBus
}

// Start begins listening for messages to publish
func (mp *MessagePublisher) Start() {
	ch := make(chan Event)
	mp.eventBus.Subscribe("messageCreate", ch)

	go func() {
		for event := range ch {
			log.Printf("Publishing message: %s", event.Payload.ID)

			// Publish message event for notifications
			mp.eventBus.Publish(Event{
				Type:    "messagePublish",
				Payload: event.Payload,
			})
		}
	}()
}

// MessageNotifier notifies clients about new messages
type MessageNotifier struct {
	eventBus *EventBus
	clients  map[string]chan Message
	mutex    sync.RWMutex
}

// NewMessageNotifier creates a new message notifier
func NewMessageNotifier(eventBus *EventBus) *MessageNotifier {
	return &MessageNotifier{
		eventBus: eventBus,
		clients:  make(map[string]chan Message),
	}
}

// RegisterClient registers a client to receive notifications
func (mn *MessageNotifier) RegisterClient(clientID string) chan Message {
	mn.mutex.Lock()
	defer mn.mutex.Unlock()
	ch := make(chan Message, 10)
	mn.clients[clientID] = ch
	return ch
}

// Start begins listening for messages to notify about
func (mn *MessageNotifier) Start() {
	ch := make(chan Event)
	mn.eventBus.Subscribe("messagePublish", ch)

	go func() {
		for event := range ch {
			log.Printf("Notification for message: %s", event.Payload.ID)

			// Notify all clients except the sender
			mn.mutex.RLock()
			for clientID, clientCh := range mn.clients {
				if clientID != event.Payload.Sender {
					select {
					case clientCh <- event.Payload:
						// Message sent to client
					default:
						// Client buffer is full, skip notification
					}
				}
			}
			mn.mutex.RUnlock()
		}
	}()
}

func main() {
	// Initialize the event bus
	eventBus := NewEventBus()

	// Initialize components
	messageReceiver := &MessageReceiver{eventBus: eventBus}
	messageSaver := &MessageSaver{eventBus: eventBus, messages: []Message{}}
	messagePublisher := &MessagePublisher{eventBus: eventBus}
	messageNotifier := NewMessageNotifier(eventBus)

	// Start all components
	messageReceiver.Start()
	messageSaver.Start()
	messagePublisher.Start()
	messageNotifier.Start()

	// Create clients
	alice := &Client{ID: "alice", eventBus: eventBus}
	bob := &Client{ID: "bob", eventBus: eventBus}

	// Register clients for notifications
	aliceChannel := messageNotifier.RegisterClient("alice")
	bobChannel := messageNotifier.RegisterClient("bob")

	// Start listening for notifications in separate goroutines
	go func() {
		for msg := range aliceChannel {
			fmt.Printf("Alice received: %s from %s\n", msg.Content, msg.Sender)
		}
	}()

	go func() {
		for msg := range bobChannel {
			fmt.Printf("Bob received: %s from %s\n", msg.Content, msg.Sender)
		}
	}()

	// Send some messages
	alice.SendMessage("Hello, everyone!")
	time.Sleep(100 * time.Millisecond) // Give time for processing

	bob.SendMessage("Hi Alice, how are you?")
	time.Sleep(100 * time.Millisecond)

	alice.SendMessage("I'm good, thanks!")

	// Wait to see the results
	time.Sleep(500 * time.Millisecond)

	// Print saved messages
	fmt.Println("\nSaved Messages:")
	for _, msg := range messageSaver.GetMessages() {
		fmt.Printf("[%s] %s: %s\n",
			msg.Timestamp.Format("15:04:05"), msg.Sender, msg.Content)
	}
}
