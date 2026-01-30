package models

import (
	"time"

	"example.com/event-booking/db"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

type EventParticipants struct {
	ID          int64
	Name        string
	Description string
	Location    string
	Email       string
}

func (e *Event) Save() error {
	query := `
		INSERT INTO events(name, description, location, date_time, user_id)
		VALUES (?, ?, ?, ?, ?)
	`
	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)

		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"

	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (event Event) Update() error {
	query := `
		UPDATE events
		SET
			name = ?,
			description = ?,
			location = ?,
			date_time = ?
		WHERE
			id = ?
	`

	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := `
		DELETE FROM events WHERE id = ?
	`

	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.ID)
	return err
}

func (e Event) Register(userId int64) error {
	query := `
		INSERT INTO registrations(event_id, user_id)
		VALUES (?, ?)
	`
	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(e.ID, userId)

	return err
}

func (e Event) CancelRegister(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	statement, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(e.ID, userId)

	return err
}

func GetEventParticipants(id int64) ([]EventParticipants, error) {
	query := `
		SELECT 
			e.id,
			e.name,
			e.description,
			e.location,
			u.email AS participant
		FROM 
			registrations r 
			JOIN events e ON r.event_id = e.id
			JOIN users u ON r.user_id = u.id
		WHERE 
			e.id = ?`

	rows, err := db.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]EventParticipants, 0)

	for rows.Next() {
		var row EventParticipants
		err := rows.Scan(
			&row.ID,
			&row.Name,
			&row.Description,
			&row.Location,
			&row.Email,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, row)
	}

	return result, nil
}
