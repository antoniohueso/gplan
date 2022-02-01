package defmodel

import (
	"time"

	"com.github.antoniohueso/gplan"
)

// ProjectPlan Contiene información de la planificación del proyecto
// Solo pongo etiquetas bson a aquellas por las que voy a buscar o de lo contrario no funciona
type ProjectPlan struct {
	// Identificador del plan de proyecto
	ID gplan.ProjectPlanID `bson:"_id"`
	// Fecha planificada de comienzo del proyecto
	StartDate time.Time `bson:"startDate"`
	// Fecha planificada de fin del proyecto
	EndDate time.Time `bson:"endDate"`
	// Lista de tareas planificadas en el proyecto
	Tasks ArrayOfTasks
	// Lista de recursos
	Resources ArrayOfResources
	// Días de fiesta
	FeastDays ArrayOfHolidays
	// Porcentaje real completado
	Complete int
	// Porcentaje completado según lo planificado
	EstimatedComplete int
	// Avance o retraso real en días
	RealAdvancedOrDelayed float64
	// Fecha en la que el proyecto fue archivado
	ArchivedDate *time.Time `bson:"archivedDate"`
}

// GetID Getter de ID
func (p ProjectPlan) GetID() gplan.ProjectPlanID {
	return p.ID
}

// GetStartDate Getter de StartDate
func (p ProjectPlan) GetStartDate() time.Time {
	return p.StartDate
}

// SetStartDate Setter de StartDate
func (p *ProjectPlan) SetStartDate(date time.Time) {
	p.StartDate = date
}

// GetEndDate Getter de EndDate
func (p ProjectPlan) GetEndDate() time.Time {
	return p.EndDate
}

// SetEndDate Setter de EndDate
func (p *ProjectPlan) SetEndDate(date time.Time) {
	p.EndDate = date
}

// GetTasks Getter de Tasks
func (p ProjectPlan) GetTasks() gplan.ArrayOfTasks {
	return p.Tasks
}

// GetResources Getter de Resources
func (p ProjectPlan) GetResources() ArrayOfResources {
	return p.Resources
}

// GetFeastDays Getter de FeastDays
func (p ProjectPlan) GetFeastDays() ArrayOfHolidays {
	return p.FeastDays
}

// GetComplete Getter de Complete
func (p ProjectPlan) GetComplete() int {
	return p.Complete
}

// SetComplete Setter de Complete
func (p *ProjectPlan) SetComplete(n int) {
	p.Complete = n
}

// GetEstimatedComplete Getter de EstimatedComplete
func (p ProjectPlan) GetEstimatedComplete() int {
	return p.EstimatedComplete
}

// SetEstimatedComplete Setter de EstimatedComplete
func (p *ProjectPlan) SetEstimatedComplete(n int) {
	p.EstimatedComplete = n
}

// ProjectPlan interfaz para implementar el acceso a la información necesaria para poder crear o revisar una planificación
type ProjectPlan1 interface {
	// Avance o retraso real en días
	GetRealAdvancedOrDelayed() float64
	SetRealAdvancedOrDelayed(n float64) float64
	// Si está archivado o no
	IsArchived() bool
	SetArchived(archived bool)
	// Fecha en la que el proyecto fue archivado
	GetArchivedDate() time.Time
	SetArchivedDate(date time.Time)
}
