package store

import (
	"github.com/pkg/errors"

	"github.com/shijuvar/go-distributed-sys/pb"
)

type EventStore struct{}

func (store EventStore) CreateEvent(event *pb.Event) error {
	// Insert two rows into the "accounts" table.
	// sql := fmt.Sprintf("INSERT INTO events (id, eventtype, aggregateid, aggregatetype, eventdata, channel)
	// VALUES ('%s','%s','%s','%s','%s','%s')", event.EventId, event.EventType, event.AggregateId, event.AggregateType, event.EventData, event.Channel)
	sql := `
INSERT INTO events (id, eventtype, aggregateid, aggregatetype, eventdata, channel) 
VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(sql, event.EventId, event.EventType, event.AggregateId, event.AggregateType, event.EventData, event.Channel)
	if err != nil {
		return errors.Wrap(err, "error on insert into events")
	}
	return nil
}

func (store EventStore) GetEvents(filter *pb.EventFilter) []*pb.Event {
	var events []*pb.Event
	return events
}
