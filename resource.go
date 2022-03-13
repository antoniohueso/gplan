package gplan

import (
	"time"
)

type IResource interface {
	Base() *ResourceBase
	GetHolidays() []IHolidays
}

// ResourceID alias del ID único de un recurso
type ResourceID string

// ResourceBase Contiene información de un recurso
type ResourceBase struct {
	// Nombre del recurso
	ID ResourceID `json:"name"`
	// Descripción del recurso
	Description string `json:"description"`
	// Clasifica el tipo de recurso, Deberá coincidir con el tipo de tarea a asignarle.
	Type string `json:"type"`
	// Fecha desde la que estará disponible para poder asignarle una tarea
	AvailableFrom time.Time `json:"availableFrom"`
	// Siguiente fecha de disponibilidad. De uso interno para calcular la siguiente fecha de disponibilidad
	nextAvailableDate time.Time
}

// NewResource crea un nuevo recurso
func NewResourceBase(id ResourceID, description string, resourceType string, availableFrom time.Time) *ResourceBase {

	return &ResourceBase{
		ID:                id,
		Description:       description,
		Type:              resourceType,
		AvailableFrom:     availableFrom,
		nextAvailableDate: availableFrom,
	}
}
