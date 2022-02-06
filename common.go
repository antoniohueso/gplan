package gplan

import (
	"fmt"
	"time"

	"com.github.antoniohueso/gplan/dateutil"
)

// calculateLaborableDays Devuelve los días laborables que hay entre dos fechas, incluidas ambas.
func calculateLaborableDays(from time.Time, to time.Time, holidays []IHolidays) int {
	days := 0
	date := from
	for dateutil.IsLte(date, to) {
		if isLaborableDay(date, holidays) {
			days++
		}
		date = date.AddDate(0, 0, 1)
	}
	return days
}

// isLaborableDay devuelve True si el día que recibe como parámetro es laborable
func isLaborableDay(day time.Time, feastDays []IHolidays) bool {

	// Si el día que recibimos es sábado o domingo devuelve False
	if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
		return false
	}

	// Si el día que recibimos está entre los rangos de vacaciones o días de fiesta devuelve False
	for _, h := range feastDays {
		if dateutil.IsBetween(day, h.GetFrom(), h.GetTo()) {
			return false
		}
	}

	return true
}

// Crea un Message de tipo Error solo con un mensaje de texto
func newTextError(message string, params ...interface{}) *Error {
	return newError(message, nil, params...)
}

// newError Crea un objeto de tipo Error
func newError(message string, tasks []TaskID, params ...interface{}) *Error {
	return &Error{
		Message: fmt.Errorf(message, params...),
		Tasks:   tasks,
	}
}
