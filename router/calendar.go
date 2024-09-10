package router

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/Depado/pb-templ-htmx-tailwind/components"
	"github.com/Depado/pb-templ-htmx-tailwind/components/calendar"
	"github.com/Depado/pb-templ-htmx-tailwind/htmx"
	"github.com/Depado/pb-templ-htmx-tailwind/models"
	"github.com/Depado/pb-templ-htmx-tailwind/services"
	"github.com/Depado/pb-templ-htmx-tailwind/services/di"
)

func (ar *AppRouter) GetPrevMonthCalendar(c echo.Context) error {
	return ar.getMonthAddedCalendar(c, -1)
}

func (ar *AppRouter) GetNextMonthCalendar(c echo.Context) error {
	return ar.getMonthAddedCalendar(c, +1)
}

func (ar *AppRouter) getMonthAddedCalendar(c echo.Context, monthsToAdd int) error {
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec == nil {
		return htmx.Redirect(c, "/")
	}

	dateStr := c.PathParam("date")
	date, err := di.Instance().CalendarProvider.ParseDate(dateStr)
	if err != nil {
		return htmx.Error(c, err.Error())
	}
	givenYear, givenMonth, _ := date.Date()
	firstOfMonth := time.Date(givenYear, givenMonth, 1, 0, 0, 0, 0, date.Location())
	firstOfThatMonth := firstOfMonth.AddDate(0, monthsToAdd, 0)
	// fmt.Printf("first of the requested month: %s\n", firstOfThatMonth.String())

	thatMonthWindow := di.Instance().CalendarProvider.CalendarMonthWindow(firstOfThatMonth)
	firstOfWindow, err := di.Instance().CalendarProvider.ParseDate(thatMonthWindow[0].Day.DateString)
	lastOfWindow, err := di.Instance().CalendarProvider.ParseDate(thatMonthWindow[len(thatMonthWindow)-1].Day.DateString)
	if err != nil {
		return htmx.Error(c, err.Error())
	}

	// fmt.Printf("first: %s | last: %s\n", firstOfWindow.String(), lastOfWindow.String())
	events, err := models.GetEventsBetweenDays(ar.App.Dao(), firstOfWindow, lastOfWindow)
	if err != nil {
		ar.App.Logger().Error("Unable to get events", "error", err)
		return htmx.Error(c, "Unable to get events")
	}

	return getCalendar(c, thatMonthWindow, events, firstOfThatMonth)
}

func getCalendar(c echo.Context,
	monthWindow []services.CalendarWindowDayEntry,
	events []*models.Event,
	date time.Time) error {
	calCtx := calendar.Context{
		DayEntries:       monthWindow,
		FirstOfCurrMonth: date,
		EventsByDate:     organizeEventsByDate(events),
	}
	return components.Render(c, http.StatusOK, calendar.Calendar(calCtx))
}
