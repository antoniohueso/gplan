package gplan

import (
	"fmt"
	"time"

	"github.com/antoniohueso/gplan/dateutil"
)

// CalculateLaborableDate devuelve una fecha laborable a partir de una fecha sumando o restando los días que recibe como parámetro.
func CalculateLaborableDate(from time.Time, days int, holidays []Holidays) time.Time {

	var (
		increment int
		date      time.Time = from
	)

	if days < 0 {
		increment = 1
	} else {
		increment = -1
	}

	for days != 0 {
		// Si days es < 0 restará uno, si es > 0 debe sumar 1
		date = date.AddDate(0, 0, increment*(-1))
		if IsLaborableDay(date, holidays) {
			days += increment
		}
	}
	return date
}

// CalculateLaborableDays Devuelve los días laborables que hay entre dos fechas, incluidas ambas.
func CalculateLaborableDays(from time.Time, to time.Time, holidays []Holidays) uint {
	var days uint
	date := from
	for dateutil.IsLte(date, to) {
		if IsLaborableDay(date, holidays) {
			days++
		}
		date = date.AddDate(0, 0, 1)
	}
	return days
}

// IsLaborableDay devuelve True si el día que recibe como parámetro es laborable
func IsLaborableDay(day time.Time, feastDays []Holidays) bool {

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

// Version versión de la librería gplan
func Version() string {
	return "1.18_2"
}
