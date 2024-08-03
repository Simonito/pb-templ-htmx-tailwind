package models

import (
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

func GetEvents(dao *daos.Dao) ([]*Event, error) {
	events := []*Event{}

	err := EventQuery(dao).All(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}
