package eventstorerepository

import (
	"database/sql"
	"fmt"

	"github.com/shijuvar/go-distributed-sys/eventstream/eventstore"
)

type repository struct {
	db *sql.DB
}

// New returns a concrete repository backed by CockroachDB
func New(db *sql.DB) (eventstore.Repository, error) {
	// return  repository
	return &repository{
		db: db,
	}, nil
}

func (repo repository) CreateEvent(event *eventstore.Event) error {
	// Insert two rows into the "events" table.
	// sql := fmt.Sprintf("INSERT INTO events (id, eventtype, aggregateid, aggregatetype, eventdata, channel)
	// VALUES ('%s','%s','%s','%s','%s','%s')", event.EventId, event.EventType, event.AggregateId, event.AggregateType, event.EventData, event.Channel)
	sql := `
INSERT INTO events (id, eventtype, aggregateid, aggregatetype, eventdata, stream) 
VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := repo.db.Exec(sql, event.EventId, event.EventType, event.AggregateId, event.AggregateType, event.EventData, event.Stream)
	if err != nil {
		return fmt.Errorf("error on insert into events: %w", err)
	}
	return nil
}

// To Do: GetEvents
func (repo repository) GetEvents(filter *eventstore.GetEventsRequest) ([]*eventstore.Event, error) {
	var events []*eventstore.Event
	return events, nil
}
