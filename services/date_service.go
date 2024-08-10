package services

import "time"

type DateService interface {
    GetTodayDate() time.Time
}

type ServerDateService struct {
}

func (*ServerDateService) GetTodayDate() time.Time {
    return time.Now()
}
