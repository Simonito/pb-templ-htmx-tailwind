package services

import (
	"fmt"
	"time"
)

type CalendarDay struct {
	Date       int
	Month      int
	Year       int
	DateString string
	DayOfWeek  string
}

type CalendarWindowDayEntry struct {
	Day         *CalendarDay
	IsThisMonth bool
}

type CalendarService interface {
	CurrentMonthWindow() (dayEntries []*CalendarWindowDayEntry, err error)
}

type DefaultCalendarService struct {
	DateService DateService
}

func (calendarServ *DefaultCalendarService) CurrentMonthWindow() ([]*CalendarWindowDayEntry, error) {
	thisDate := calendarServ.DateService.GetTodayDate()
	return nil, nil
}

func newCalendarDay(date int, month int, year int) *CalendarDay {
	dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, date)
	tim, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(err)
	}
	return &CalendarDay{
		Date:       date,
		Month:      month,
		Year:       year,
		DateString: dateStr,
		DayOfWeek:  tim.Weekday().String(),
	}
}

func getCalendarStartDays(currentDate time.Time) []time.Time {
	startMonthDate := time.Date(
		currentDate.Year(),
		currentDate.Month(),
		1,
		0, 0, 0, 0, time.Now().Location(),
	)
	startMonthWeekDay := int(startMonthDate.Weekday())

	firstDay := startMonthDate.Add(time.Duration(startMonthWeekDay-1) * time.Hour * -24)
	if firstDay.Day() == 1 {
		// no days from the previous month to display
		return []time.Time{}
	}

	var days []time.Time
	for day := firstDay; day.Before(startMonthDate); day = day.Add(24 * time.Hour) {
		days = append(days, day)
	}

	return days
}
