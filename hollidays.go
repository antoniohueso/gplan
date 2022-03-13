package gplan

import (
	"time"
)

type IHolidays interface {
	GetHolidays() *Holidays
}

// HolidaysBase contiene información de un rango de fechas de vacaciones o días de fiesta
type Holidays struct {
	// Fecha de vacaciones/fiesta desde
	From time.Time `json:"from"`
	// Fecha de vacaciones/fiesta Hasta
	To time.Time `json:"to"`
}

// NewHolidays crea un nuevo rango de fechas de vacaciones
func NewHolidays(from time.Time, to time.Time) *Holidays {
	return &Holidays{
		From: from,
		To:   to,
	}
}

func (s *Holidays) GetHolidays() *Holidays {
	return s
}
