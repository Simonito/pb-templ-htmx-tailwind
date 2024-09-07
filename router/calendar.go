package router

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"

	"github.com/Depado/pb-templ-htmx-tailwind/components"
	"github.com/Depado/pb-templ-htmx-tailwind/components/calendar"
	"github.com/Depado/pb-templ-htmx-tailwind/htmx"
	"github.com/Depado/pb-templ-htmx-tailwind/services/di"
)

func (ar *AppRouter) GetPrevMonthCalendar(c echo.Context) error {
	return getMonthAddedCalendar(c, -1)
}

func (ar *AppRouter) GetNextMonthCalendar(c echo.Context) error {
	return getMonthAddedCalendar(c, +1)
}

func getMonthAddedCalendar(c echo.Context, monthsToAdd int) error {
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

	return getCalendar(c, firstOfThatMonth)
}

func getCalendar(c echo.Context, date time.Time) error {
	calCtx := calendar.Context{
		DayEntries:       di.Instance().CalendarProvider.CalendarMonthWindow(date),
		FirstOfCurrMonth: date,
	}
	return components.Render(c, http.StatusOK, calendar.Calendar(calCtx))
}
