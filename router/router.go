package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/pocketbase/pocketbase/core"

	"github.com/Depado/pb-templ-htmx-tailwind/assets"
	"github.com/Depado/pb-templ-htmx-tailwind/components"
	"github.com/Depado/pb-templ-htmx-tailwind/components/test"
	"github.com/Depado/pb-templ-htmx-tailwind/htmx"
	"github.com/Depado/pb-templ-htmx-tailwind/models"
)

type AppRouter struct {
	App    core.App
	Router *echo.Echo
}

func NewAppRouter(e *core.ServeEvent) *AppRouter {
	return &AppRouter{
		App:    e.App,
		Router: e.Router,
	}
}

func (ar *AppRouter) SetupRoutes(live bool) error {
	ar.Router.Use(middleware.Logger())
	ar.Router.HTTPErrorHandler = htmx.WrapDefaultErrorHandler(ar.Router.HTTPErrorHandler)
	ar.Router.GET("/static/*", assets.AssetsHandler(ar.App.Logger(), live), middleware.Gzip())

	ar.Router.Use(ar.LoadAuthContextFromCookie())
	ar.Router.GET("/", ar.GetHome)
	ar.Router.GET("/login", ar.GetLogin)
	ar.Router.POST("/login", ar.PostLogin)
	ar.Router.POST("/register", ar.PostRegister)
	ar.Router.POST("/logout", ar.PostLogout)
	ar.Router.GET("/error", ar.GetError)

	ar.Router.PATCH("/task/:id", ar.ToggleTask)
	ar.Router.POST("/list/:id/task", ar.CreateTask)
	ar.Router.PATCH("/list/:id/archive", ar.ToggleArchive)

	err := ar.setupHtmxRoutes()

	ar.Router.GET("/test", ar.TestPath)
	return err
}

func (ar *AppRouter) setupHtmxRoutes() error {
	ar.Router.GET("/calendar/nextMonth/:date", ar.GetNextMonthCalendar)
	ar.Router.GET("/calendar/prevMonth/:date", ar.GetPrevMonthCalendar)

	return nil
}

func (ar *AppRouter) TestPath(c echo.Context) error {
	records, err := models.GetEventsOfMonth(ar.App.Dao(), time.Date(2024, time.August, 1, 1, 1, 1, 1, time.Now().Location()))
	if err != nil {
		return htmx.Error(c, err.Error())
	}
	return components.Render(c, http.StatusOK, test.Test(strconv.Itoa(len(records))))

}
