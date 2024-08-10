package di

import (
    "sync"
    "github.com/Depado/pb-templ-htmx-tailwind/services"
)

type provider struct {
    DateProvider services.DateService;
    CalendarProvider services.CalendarService;
}

var instance *provider
var once sync.Once

func Instance() *provider {
	once.Do(func() {
		instance = &provider{}
        instance.initializeServices()
	})
	return instance
}

func (p *provider) initializeServices() {
    p.DateProvider = initDateService()
    p.CalendarProvider = initCalendarService(p.DateProvider)
}

func initDateService() services.DateService {
    return &services.ServerDateService{}
}

func initCalendarService(dateProv services.DateService) services.CalendarService {
    return &services.DefaultCalendarService{
        DateService: dateProv,
    }
}
