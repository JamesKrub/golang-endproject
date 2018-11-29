package post

import (
	"database/sql"
	"errors"
	"log"
	"strconv"

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
	Id         int
	EvId       int
	StartTime  string
	EndTime    string
	StartDate  string
	EndDate    string
	Limit      string
	Count      int
	ShowOption bool
}

type UpdateEvent struct {
	Main   Event
	Detail []EventDetail
}

type Register struct {
	Id     int
	FName  string
	LName  string
	Tel    string
	UserId string
	Event  string
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

func FindDetailByID(id int) (UpdateEvent, error) {
	var rs UpdateEvent
	var ev Event
	row := db.QueryRow("SELECT name, place, speaker, detail	FROM events WHERE id = $1", id)
	err := row.Scan(&ev.Name, &ev.Place, &ev.Speaker, &ev.Detail)
	if err != nil {
		return rs, err
	}

	rows, err := db.Query("SELECT id, start_date, end_date, limitation, start_time, end_time FROM event_detail WHERE evnt_id = $1 ORDER BY id ASC", id)
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
		tempID := strconv.Itoa(v.Id)
		count, err := GetAmountJoiner(tempID)
		if err != nil {
			return rs, err
		}
		v.Count = count

		limit, err := strconv.Atoi(v.Limit)
		if err != nil {
			return rs, err
		}
		if count < limit {
			v.ShowOption = true
		}
		// fmt.Println(v.ShowOption)
		evd = append(evd, v)
	}

	rs = UpdateEvent{ev, evd}
	return rs, nil
}

func GetLimit(id string) (int, error) {
	var limit int
	row := db.QueryRow("SELECT limitation FROM event_detail WHERE evnt_id = $1 ORDER BY id ASC", id)
	err := row.Scan(&limit)
	if err != nil {
		return limit, err
	}
	return limit, nil
}

func GetAmountJoiner(id string) (int, error) {
	var amount int
	row := db.QueryRow("SELECT count(1) as amount FROM Register WHERE event = $1", id)
	err := row.Scan(&amount)
	if err != nil {
		return amount, err
	}
	return amount, nil
}

func InsertRegister(reg *Register) error {
	var (
		amount int
		limit  int
	)

	limit, err := GetLimit(reg.Event)
	if err != nil {
		return err
	}

	amount, err = GetAmountJoiner(reg.Event)
	if err != nil {
		return err
	}

	if amount >= limit {
		return errors.New("amount of joiner exceed the limit")
	}
	_, err = db.Exec("INSERT INTO Register(fName, lName, tel, userId, event) VALUES ($1, $2, $3, $4, $5)", reg.FName, reg.LName, reg.Tel, reg.UserId, reg.Event)

	if err != nil {
		return err
	}
	return nil
}

func GetRegister(id int) ([]Register, error) {
	var reg []Register
	rows, err := db.Query("SELECT fname, lname, tel, userid FROM Register WHERE event = $1 ORDER BY event ASC", id)
	if err != nil {
		return reg, err
	}

	for rows.Next() {
		var r Register
		err := rows.Scan(&r.FName, &r.LName, &r.Tel, &r.UserId)
		if err != nil {
			return nil, err
		}
		reg = append(reg, r)
	}
	return reg, nil
}
