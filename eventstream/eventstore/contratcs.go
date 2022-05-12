package eventstore

type Repository interface {
	CreateEvent(event *Event) error
	GetEvents(filter *GetEventsRequest) ([]*Event, error)
}
