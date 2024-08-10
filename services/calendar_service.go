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
	CurrentMonthWindow() []CalendarWindowDayEntry
}

type DefaultCalendarService struct {
	DateService DateService
}

const dateFormat = "2006-01-02"

// Creates and returnes a so called window of a current month for calendar
// The window is basically all the dates that should get displayed in a calendar
//
// imagine this month
// May 2024:
// | MON | TUE | WED | THU | FRI | SAT | SUN |
// | _29 | _30 | 1   | 2   | 3   | 4   | 5   |
// | .... .the other day of that month ..... |
// | 27  | 28  | 29  | 30  | 31  | _1  | _2  |
//
//	where the underscored numbers represent days of different month.
func (calendarServ *DefaultCalendarService) CurrentMonthWindow() []CalendarWindowDayEntry {
	currentDate := calendarServ.DateService.GetTodayDate()

	// calculate the dates that need to be "prepended" from prev month
	startDays := getCalendarStartDays(currentDate)

	var calendarDayEntries []CalendarWindowDayEntry
	for _, day := range startDays {
		calEntry, err := newCalendarDay(day.Day(), int(day.Month()), day.Year())
		if err != nil {
			return nil
		}
		calendarDayEntries = append(calendarDayEntries, CalendarWindowDayEntry{
			Day:         calEntry,
			IsThisMonth: false,
		})
	}

	// append the all the dates of the current month
	currentYear, currentMonth, _ := currentDate.Date()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentDate.Location())
	firstOfNextMonth := firstOfMonth.AddDate(0, 1, 0)
	lastOfMonth := firstOfNextMonth.Add(-time.Hour * 24)

	for day := 1; day <= lastOfMonth.Day(); day++ {
		calEntry, err := newCalendarDay(day, int(firstOfMonth.Month()), firstOfMonth.Year())
		if err != nil {
			return nil
		}
		calendarDayEntries = append(calendarDayEntries, CalendarWindowDayEntry{
			Day:         calEntry,
			IsThisMonth: true,
		})
	}

	// finish off by adding the calculated days of next month to complete the window
	endDays := getCalendarEndDays(lastOfMonth)
	for _, day := range endDays {
		calEntry, err := newCalendarDay(day.Day(), int(day.Month()), day.Year())
		if err != nil {
			return nil
		}
		calendarDayEntries = append(calendarDayEntries, CalendarWindowDayEntry{
			Day:         calEntry,
			IsThisMonth: false,
		})
	}
	return calendarDayEntries
}

func newCalendarDay(date int, month int, year int) (*CalendarDay, error) {
	dateStr := fmt.Sprintf("%04d-%02d-%02d", year, month, date)
	tim, err := time.Parse(dateFormat, dateStr)
	if err != nil {
		return nil, err
	}
	return &CalendarDay{
		Date:       date,
		Month:      month,
		Year:       year,
		DateString: dateStr,
		DayOfWeek:  tim.Weekday().String(),
	}, nil
}

func dayOfWeek(date time.Time) int {
	dow := int(date.Weekday())
	if dow == 0 {
		dow = 7
	}

	return dow
}

func getCalendarStartDays(currentDate time.Time) []time.Time {
	startMonthDate := time.Date(
		currentDate.Year(),
		currentDate.Month(),
		1,
		0, 0, 0, 0, time.Now().Location(),
	)
	startMonthWeekDay := dayOfWeek(startMonthDate)

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

func getCalendarEndDays(lastOfMonth time.Time) []time.Time {
	lastDowOfMonth := dayOfWeek(lastOfMonth)

	var days []time.Time
	remain := 7 - lastDowOfMonth
	for i := 1; i <= remain; i++ {
		days = append(days, lastOfMonth.Add(time.Duration(i)*time.Hour*24))
	}

	return days
}
