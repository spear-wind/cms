package events

type EventPublisher interface {
	Publish(event interface{})
	Add(s EventSubscriber)
}

type EventSubscriber interface {
	Receive(event interface{})
}
