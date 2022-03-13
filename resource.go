package gplan

import (
	"time"
)

type IResource interface {
	GetResource() *Resource
}

// ResourceID alias del ID único de un recurso
type ResourceID string

// Resource Contiene información de un recurso
type Resource struct {
	// Nombre del recurso
	ID ResourceID `json:"name"`
	// Descripción del recurso
	Description string `json:"description"`
	// Clasifica el tipo de recurso, Deberá coincidir con el tipo de tarea a asignarle.
	Type string `json:"type"`
	// Fecha desde la que estará disponible para poder asignarle una tarea
	AvailableFrom time.Time `json:"availableFrom"`
	// Conjunto de días de vacaciones
	Holidays []IHolidays `json:"holidays"`
	// Siguiente fecha de disponibilidad. De uso interno para calcular la siguiente fecha de disponibilidad
	nextAvailableDate time.Time
}

// NewResource crea un nuevo recurso
func NewResource(id ResourceID, description string, resourceType string, availableFrom time.Time, holidays []IHolidays) *Resource {

	if holidays == nil {
		holidays = []IHolidays{}
	}

	return &Resource{
		ID:                id,
		Description:       description,
		Type:              resourceType,
		AvailableFrom:     availableFrom,
		nextAvailableDate: availableFrom,
		Holidays:          holidays,
	}
}

// GetResource() implementa IResource
func (s *Resource) GetResource() *Resource {
	return s
}
