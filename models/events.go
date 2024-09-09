package models

import (
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Event struct {
	models.BaseModel

	// Fields
	Car          string         `db:"car"`
	Date         types.DateTime `db:"date"`
	Name         string         `db:"name"`
	Oraganizer   string         `db:"organizer"`
	Address      string         `db:"address"`
	ContactName  string         `db:"contact_name"`
	ContactPhone string         `db:"contact_phone"`
	TimeArrival  types.DateTime `db:"time_arrival"`
	TimeReady    types.DateTime `db:"time_ready"`
	TimeEnd      types.DateTime `db:"time_end"`
	Note         string         `db:"note"`

	// Relations
	IdResponsibleEmployee string         `db:"id_employee_responsible"`
	ResponsibleEmployee   *models.Record `db:"-"`
}

func (*Event) TableName() string {
	return "event"
}

func EventQuery(dao *daos.Dao) *dbx.SelectQuery {
	return dao.ModelQuery(&Event{})
}

func EventQueryByMonth(dao *daos.Dao, randomDayInMonth time.Time) *dbx.SelectQuery {
	firstOfMonth := time.Date(randomDayInMonth.Year(), randomDayInMonth.Month(), 1, 0, 0, 0, 0, randomDayInMonth.Location())
	lastOfMonth := firstOfMonth.AddDate(0, 1, 0).Add(time.Hour * -24)

	// fmt.Println("First Of Month: " + firstOfMonth.String())
	// fmt.Println("Last Of Month: " + lastOfMonth.String())

	return EventQueryBetweenDays(dao, firstOfMonth, lastOfMonth)
}

func EventQueryBetweenDays(dao *daos.Dao, firstDay time.Time, lastDay time.Time) *dbx.SelectQuery {
	filterExp := fmt.Sprintf("'%04d-%02d-01 00:00:00.000Z' <= date AND date <= '%04d-%02d-%02d 23:59:59.999Z'",
		firstDay.Year(),
		int(firstDay.Month()),
		lastDay.Year(),
		int(lastDay.Month()),
		lastDay.Day())

	// fmt.Println("Query filter Exp: " + filterExp)
	return dao.DB().
		Select("event.*").
		From("event").
		Where(dbx.NewExp(filterExp))

}

func GetEvents(dao *daos.Dao) ([]*Event, error) {
	events := []*Event{}

	err := EventQuery(dao).All(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func GetEventsOfMonth(dao *daos.Dao, dayOfMonth time.Time) ([]*Event, error) {
	events := []*Event{}

	err := EventQueryByMonth(dao, dayOfMonth).All(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func GetEventsBetweenDays(dao *daos.Dao, firstDay time.Time, lastDay time.Time) ([]*Event, error) {
	events := []*Event{}

	err := EventQueryBetweenDays(dao, firstDay, lastDay).All(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}
