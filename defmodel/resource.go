package defmodel

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
	Holidays ArrayOfHolidays `json:"hollidays"`
	// Siguiente fecha de disponibilidad. De uso interno para calcular la siguiente fecha de disponibilidad
	nextAvailableDate time.Time
}

// NewResource crea un nuevo recurso
func NewResource(id gplan.ResourceID, description string, resourceType string, availableFrom time.Time, holidays ArrayOfHolidays) *Resource {
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

// GetHolidays Getter de Holidays
func (r Resource) GetHolidays() gplan.ArrayOfHolidays {
	return r.Holidays
}

// GetNextAvailableDate Getter de nextAvailableDate
func (r Resource) GetNextAvailableDate() time.Time {
	return r.nextAvailableDate
}

// SetNextAvailableDate Setter de nextAvailableDate
func (r *Resource) SetNextAvailableDate(date time.Time) {
	r.nextAvailableDate = date
}

// ArrayOfResources implementa un gplan.ArrayOfResources
type ArrayOfResources []*Resource

// Iterable Implementa Iterable
func (a ArrayOfResources) Iterable() []gplan.Resource {
	newArr := make([]gplan.Resource, len(a))
	for i := range a {
		newArr[i] = a[i]
	}
	return newArr
}
