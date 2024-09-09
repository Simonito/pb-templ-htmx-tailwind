package router

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	pbmodels "github.com/pocketbase/pocketbase/models"

	"github.com/Depado/pb-templ-htmx-tailwind/components"
	"github.com/Depado/pb-templ-htmx-tailwind/components/calendar"
	"github.com/Depado/pb-templ-htmx-tailwind/components/shared"
	"github.com/Depado/pb-templ-htmx-tailwind/htmx"
	"github.com/Depado/pb-templ-htmx-tailwind/models"
	"github.com/Depado/pb-templ-htmx-tailwind/services/di"
)

func organizeEventsByDate(events []*models.Event) map[string][]*models.Event {
	eventMap := make(map[string][]*models.Event)

	for _, event := range events {
		dateStr := di.Instance().CalendarProvider.Format(event.Date.Time())

		eventMap[dateStr] = append(eventMap[dateStr], event)
	}
	return eventMap
}

func (ar *AppRouter) GetHome(c echo.Context) error {
	rec := c.Get(apis.ContextAuthRecordKey)
	if rec == nil {
		return components.Render(c, http.StatusOK, components.Home(components.HomeContext{}, false))
	}

	user := c.Get(apis.ContextAuthRecordKey).(*pbmodels.Record)
	lists, err := models.FindUserLists(ar.App.Dao(), user.Id)
	if err != nil {
		ar.App.Logger().Error("unable to get lists for user", "error", err, "id", user.Id)
		return htmx.Error(c, "Unable to get lists")
	}

	currMonthWindow := di.Instance().CalendarProvider.CurrentMonthWindow()
	firstOfWindow, err := di.Instance().CalendarProvider.ParseDate(currMonthWindow[0].Day.DateString)
	lastOfWindow, err := di.Instance().CalendarProvider.ParseDate(currMonthWindow[len(currMonthWindow)-1].Day.DateString)
	if err != nil {
		ar.App.Logger().Error("Error parsing date string provided by CalendarProvider", "error", err)
	}
	events, err := models.GetEventsBetweenDays(ar.App.Dao(), firstOfWindow, lastOfWindow)
	if err != nil {
		ar.App.Logger().Error("Unable to get events", "error", err)
		return htmx.Error(c, "Unable to get events")
	}

	return components.Render(
		c,
		http.StatusOK,
		components.Home(
			components.HomeContext{
				BaseContext: shared.Context{
					User:   user,
					Lists:  lists,
					Events: events,
				},
				CalendarContext: calendar.Context{
					DayEntries:       di.Instance().CalendarProvider.CurrentMonthWindow(),
					FirstOfCurrMonth: di.Instance().DateProvider.GetTodayDate(),
					EventsByDate:     organizeEventsByDate(events),
				},
			},
			false,
		),
	)
}
