package sample

import (
	"time"

	"github.com/antoniohueso/gplan"
)

type Holidays struct {
	*gplan.HolidaysBase
}

func (s *Holidays) Base() *gplan.HolidaysBase {
	return s.HolidaysBase
}

// NewHolidays crea un nuevo rango de fechas de vacaciones
func NewHolidays(from time.Time, to time.Time) *Holidays {
	return &Holidays{
		HolidaysBase: gplan.NewHolidaysBase(from, to),
	}
}
