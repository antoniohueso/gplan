package defmodel

import (
	"time"

	"com.github.antoniohueso/gplan"
)

// Holidays contiene información de un rango de fechas de vacaciones o días de fiesta
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

// GetFrom Getter de from
func (h Holidays) GetFrom() time.Time {
	return h.From
}

// GetTo Getter de To
func (h Holidays) GetTo() time.Time {
	return h.To
}

// ArrayOfHolidays Implementa gplan.ArrayOfHolidays
type ArrayOfHolidays []*Holidays

// Iterable Implementa Iterable de gplan.ArrayOfHolidays
func (a ArrayOfHolidays) Iterable() []gplan.Holidays {
	newArr := make([]gplan.Holidays, len(a))
	for i := range a {
		newArr[i] = a[i]
	}
	return newArr
}

// Get Implementa Get de gplan.ArrayOfHolidays
func (a ArrayOfHolidays) Get(i int) gplan.Holidays {
	return a[i]
}

// Len Implementa Get de gplan.ArrayOfHolidays
func (a ArrayOfHolidays) Len() int {
	return len(a)
}
