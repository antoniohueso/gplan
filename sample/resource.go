package sample

import (
	"time"

	"com.github.antoniohueso/gplan"
)

// Resource Contiene información de un recurso
type Resource struct {
	// Nombre del recurso
	ID gplan.ResourceID `json:"name"`
	// Descripción del recurso
	Description string `json:"description"`
	// Clasifica el tipo de recurso, Deberá coincidir con el tipo de tarea a asignarle.
	Type string `json:"type"`
	// Fecha desde la que estará disponible para poder asignarle una tarea
	AvailableFrom time.Time `json:"availableFrom"`
	// Conjunto de días de vacaciones
	Holidays []*Holidays `json:"hollidays"`
	// Siguiente fecha de disponibilidad. De uso interno para calcular la siguiente fecha de disponibilidad
	nextAvailableDate time.Time
}

// NewResource crea un nuevo recurso
func NewResource(id gplan.ResourceID, description string, resourceType string, availableFrom time.Time, holidays []*Holidays) *Resource {
	return &Resource{
		ID:                id,
		Description:       description,
		Type:              resourceType,
		AvailableFrom:     availableFrom,
		nextAvailableDate: availableFrom,
		Holidays:          holidays,
	}
}

// GetID Getter de ID
func (r Resource) GetID() gplan.ResourceID {
	return r.ID
}

// GetType Getter de Type
func (r Resource) GetType() string {
	return r.Type
}

// GetAvailableFrom Getter de AvailableFrom
func (r Resource) GetAvailableFrom() time.Time {
	return r.AvailableFrom
}

// GetHolidays Devuelve un nuevo array de gplan.Holidays
func (r Resource) GetHolidays() []gplan.IHolidays {
	newArr := make([]gplan.IHolidays, len(r.Holidays))
	for i := range r.Holidays {
		newArr[i] = r.Holidays[i]
	}
	return newArr
}

// GetNextAvailableDate Getter de nextAvailableDate
func (r Resource) GetNextAvailableDate() time.Time {
	return r.nextAvailableDate
}

// SetNextAvailableDate Setter de nextAvailableDate
func (r *Resource) SetNextAvailableDate(date time.Time) {
	r.nextAvailableDate = date
}
