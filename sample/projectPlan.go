package sample

import (
	"sort"
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
	Tasks []*Task
	// Lista de recursos
	Resources []*Resource
	// Días de fiesta
	FeastDays []*Holidays
	// Porcentaje real completado
	Complete int
	// Porcentaje completado según lo planificado
	EstimatedComplete int
	// Avance o retraso real en días
	RealAdvancedOrDelayed float64
	// Indica si está archivado o no
	Archived bool
	// Fecha en la que el proyecto fue archivado
	ArchivedDate time.Time `bson:"archivedDate"`
}

// NewProjectPlan crea un nuevo plan de proyecto para poder ser planificado o revisado
func NewProjectPlan(id gplan.ProjectPlanID, tasks []*Task, resources []*Resource, feastDays []*Holidays) *ProjectPlan {
	return &ProjectPlan{
		ID:        id,
		Tasks:     tasks,
		Resources: resources,
		FeastDays: feastDays,
	}
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

// GetTasks Devuelve un array de Task
func (p ProjectPlan) GetTasks() []gplan.Task {
	newArr := make([]gplan.Task, len(p.Tasks))
	for i := range p.Tasks {
		newArr[i] = p.Tasks[i]
	}
	return newArr
}

// SortTasksByOrder Implementa SortByOrder
func (p ProjectPlan) SortTasksByOrder() {
	// Ordena las tareas por número de orden para poder planificarlas
	sort.Slice(p.Tasks, func(i, j int) bool {
		return p.Tasks[i].Order < p.Tasks[j].Order
	})
}

// GetResources Devuelve un nuevo array de Resources
func (p ProjectPlan) GetResources() []gplan.Resource {
	newArr := make([]gplan.Resource, len(p.Resources))
	for i := range p.Resources {
		newArr[i] = p.Resources[i]
	}
	return newArr
}

// GetFeastDays Devuelve un nuevo array de Holidays
func (p ProjectPlan) GetFeastDays() []gplan.Holidays {
	newArr := make([]gplan.Holidays, len(p.FeastDays))
	for i := range p.FeastDays {
		newArr[i] = p.FeastDays[i]
	}
	return newArr
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

// GetRealAdvancedOrDelayed Getter de RealAdvancedOrDelayed
func (p ProjectPlan) GetRealAdvancedOrDelayed() float64 {
	return p.RealAdvancedOrDelayed
}

// SetRealAdvancedOrDelayed Setter de RealAdvancedOrDelayed
func (p *ProjectPlan) SetRealAdvancedOrDelayed(n float64) {
	p.RealAdvancedOrDelayed = n
}

// IsArchived Getter de Archived
func (p ProjectPlan) IsArchived() bool {
	return p.Archived
}

// SetArchived Getter de Archived
func (p *ProjectPlan) SetArchived(archived bool) {
	p.Archived = archived
}

// GetArchivedDate Getter de ArchivedDate
func (p ProjectPlan) GetArchivedDate() time.Time {
	return p.ArchivedDate
}

// SetArchivedDate Setter de ArchivedDate
func (p *ProjectPlan) SetArchivedDate(date time.Time) {
	p.ArchivedDate = date
}
