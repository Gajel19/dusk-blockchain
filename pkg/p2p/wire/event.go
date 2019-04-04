package wire

import (
	"bytes"

	log "github.com/sirupsen/logrus"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire/topics"
)

// QuitTopic is the topic to make all components quit
const QuitTopic = "quit"

type (

	// The Event is an Entity that represents the Messages travelling on the EventBus. It would normally present always the same fields.
	Event interface {
		Sender() []byte
		Equal(Event) bool
	}

	// EventUnmarshaller unmarshals an Event from a buffer. Following Golang's way of defining interfaces, it exposes an Unmarshal method which allows for flexibility and reusability across all the different components that need to read the buffer coming from the EventBus into different structs
	EventUnmarshaller interface {
		Unmarshal(*bytes.Buffer, Event) error
	}

	// EventMarshaller is the specular operation of an EventUnmarshaller. Following Golang's way of defining interfaces, it exposes an Unmarshal method which allows for flexibility and reusability across all the different components that need to read the buffer coming from the EventBus into different structs
	EventMarshaller interface {
		Marshal(*bytes.Buffer, Event) error
	}

	// EventUnMarshaller is a convenient interface providing both Marshalling and Unmarshalling capabilities
	EventUnMarshaller interface {
		EventMarshaller
		EventUnmarshaller
	}

	// EventPrioritizer is used by the EventSelector to prioritize events (normally to return the best collected after a timespan)
	EventPrioritizer interface {
		Priority(Event, Event) Event
	}

	// EventVerifier is the interface to verify an Event
	EventVerifier interface {
		Verify(Event) error
	}

	// EventCollector is the interface for collecting Events. Pretty much processors involves some degree of Event collection (either until a Quorum is reached or until a Timeout). This Interface is typically implemented by a struct that will perform some Event unmarshalling.
	EventCollector interface {
		Collect(*bytes.Buffer) error
	}

	// EventSubscriber accepts events from the EventBus and takes care of reacting on quit Events. It delegates the business logic to the EventCollector which is supposed to handle the incoming events
	EventSubscriber struct {
		eventBus       *EventBus
		eventCollector EventCollector
		msgChan        <-chan *bytes.Buffer
		msgChanID      uint32
		quitChan       <-chan *bytes.Buffer
		quitChanID     uint32
		topic          string
	}
)

// NewEventSubscriber creates the EventSubscriber listening to a topic on the EventBus. The EventBus, EventCollector and Topic are injected
func NewEventSubscriber(eventBus *EventBus, collector EventCollector,
	topic string) *EventSubscriber {

	quitChan := make(chan *bytes.Buffer, 1)
	msgChan := make(chan *bytes.Buffer, 100)

	msgChanID := eventBus.Subscribe(topic, msgChan)
	quitChanID := eventBus.Subscribe(string(QuitTopic), quitChan)

	return &EventSubscriber{
		eventBus:       eventBus,
		msgChan:        msgChan,
		msgChanID:      msgChanID,
		quitChan:       quitChan,
		quitChanID:     quitChanID,
		topic:          topic,
		eventCollector: collector,
	}
}

// Accept incoming (mashalled) Events on the topic of interest and dispatch them to the EventCollector.Collect
func (n *EventSubscriber) Accept() {
	for {
		select {
		case <-n.quitChan:
			n.eventBus.Unsubscribe(n.topic, n.msgChanID)
			n.eventBus.Unsubscribe(string(QuitTopic), n.quitChanID)
			return
		case eventMsg := <-n.msgChan:
			if len(n.msgChan) > 1 {
				log.WithFields(log.Fields{
					"id":         n.msgChanID,
					"topic":      n.topic,
					"Unconsumed": len(n.msgChan),
				}).Warnln("Channel is accumulating messages")
			} else {
				log.WithFields(log.Fields{
					"id":    n.msgChanID,
					"topic": n.topic,
				}).Debug("Channel clean")
			}
			_ = n.eventCollector.Collect(eventMsg)
		}
	}
}

func AddTopic(m *bytes.Buffer, topic topics.Topic) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	topicBytes := topics.TopicToByteArray(topic)
	if _, err := buffer.Write(topicBytes[:]); err != nil {
		return nil, err
	}

	if _, err := buffer.Write(m.Bytes()); err != nil {
		return nil, err
	}

	return buffer, nil
}
