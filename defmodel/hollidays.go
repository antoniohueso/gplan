package defmodel

import "time"

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
func (h *Holidays) GetFrom() time.Time {
	return h.From
}

// GetTo Getter de To
func (h *Holidays) GetTo() time.Time {
	return h.To
}
