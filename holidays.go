package gplan

import (
	"time"
)

type IHolidays interface {
	Base() *HolidaysBase
}

// HolidaysBase contiene información de un rango de fechas de vacaciones o días de fiesta
type HolidaysBase struct {
	// Fecha de vacaciones/fiesta desde
	From time.Time `json:"from"`
	// Fecha de vacaciones/fiesta Hasta
	To time.Time `json:"to"`
}

// NewHolidays crea un nuevo rango de fechas de vacaciones
func NewHolidaysBase(from time.Time, to time.Time) *HolidaysBase {
	return &HolidaysBase{
		From: from,
		To:   to,
	}
}
