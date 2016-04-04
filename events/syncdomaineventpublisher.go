package events

type synchronousEventPublisher struct {
	subscribers []EventSubscriber
}

func (p synchronousEventPublisher) Publish(e interface{}) {
	for _, subscriber := range p.subscribers {
		subscriber.Receive(e)
	}
}

func (p synchronousEventPublisher) Add(s EventSubscriber) {
	p.subscribers = append(p.subscribers, s)
}

//NewSynchEventPublisher returns a simple, synchronous EventPublisher
func NewSynchEventPublisher() EventPublisher {
	return synchronousEventPublisher{
		subscribers: []EventSubscriber{},
	}
}
