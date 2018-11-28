package post

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDB() {
	var err error
	connStr := "postgres://ndpluaxa:po9QvIZ_VqpCicoDvIWjY43vH597Pd1c@elmer.db.elephantsql.com:5432/ndpluaxa"
	// const connStr = "user=postgres dbname=blog-db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

type Event struct {
	Id           int
	Name         string
	StartTime    string
	EndTime      string
	Place        string
	Speaker      string
	Detail       string
	StartDate    string
	EndDate      string
	Limit        string
	EventId      int
	EvntDetailId int
	Count        int
}

type EventDetail struct {
	Id        int
	EvId      int
	StartTime string
	EndTime   string
	StartDate string
	EndDate   string
	Limit     string
}

type updateEvent struct {
	Main   Event
	Detail []EventDetail
}

func Insert(e *Event) error {
	r := db.QueryRow("INSERT INTO events(name,place,speaker,detail) VALUES ($1,$2,$3,$4) RETURNING id", e.Name, e.Place, e.Speaker, e.Detail)
	err := r.Scan(&e.Id)
	if err != nil {
		return err
	}
	return nil
}

func InsertEventDetails(e *Event) {
	db.QueryRow("INSERT INTO event_detail(start_date,end_date,limitation,evnt_id, start_time, end_time) VALUES ($1,$2,$3,$4,$5,$6)", e.StartDate, e.EndDate, e.Limit, e.EventId, e.StartTime, e.EndTime)
}

func All() ([]Event, error) {
	var rs []Event
	rows, err := db.Query(`SELECT id, name, place, speaker, detail, 
				(SELECT count(1) FROM event_detail WHERE evnt_id = events.id) as count 
				FROM events`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r Event
		err := rows.Scan(&r.Id, &r.Name, &r.Place, &r.Speaker, &r.Detail, &r.Count)
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
	}
	return rs, nil
}

func Delete(id int) error {
	_, err := db.Exec("DELETE FROM event_detail WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func FindEventByID(id int) (Event, error) {
	var ev Event
	ev.Id = id
	row := db.QueryRow("SELECT name, place, speaker, detail	FROM events WHERE id = $1", id)
	err := row.Scan(&ev.Name, &ev.Place, &ev.Speaker, &ev.Detail)
	if err != nil {
		return ev, err
	}

	return ev, nil
}

func FindDetailByID(id int) (updateEvent, error) {
	var rs updateEvent
	var ev Event
	row := db.QueryRow("SELECT name, place, speaker, detail	FROM events WHERE id = $1", id)
	err := row.Scan(&ev.Name, &ev.Place, &ev.Speaker, &ev.Detail)
	if err != nil {
		return rs, err
	}

	rows, err := db.Query("SELECT id, start_date, end_date, limitation, start_time, end_time FROM event_detail WHERE evnt_id = $1", id)
	if err != nil {
		return rs, err
	}
	var evd []EventDetail
	for rows.Next() {
		var v EventDetail
		err := rows.Scan(&v.Id, &v.StartDate, &v.EndDate, &v.Limit, &v.StartTime, &v.EndTime)
		if err != nil {
			return rs, err
		}
		evd = append(evd, v)
	}

	rs = updateEvent{ev, evd}
	return rs, nil
}

// func Save(p *Post) error {
// 	_, err := db.Exec("UPDATE posts SET title = $1, body = $2 WHERE id = $3", p.Title, p.Body, p.ID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func AddComment(p *Post, c *Comment) error {
// 	r := db.QueryRow("INSERT INTO comments(body, post_id) VALUES ($1,$2) RETURNING id", c.Body, p.ID)
// 	err := r.Scan(&c.ID)
// 	if err != nil {
// 		return err
// 	}
// 	c.PostID = p.ID
// 	p.Comments = append(p.Comments, *c)
// 	return nil
// }
