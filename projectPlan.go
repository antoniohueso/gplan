package gplan

import (
	"time"
)

type IProjectPlan interface {
	Base() *ProjectPlanBase
	GetTasks() []ITask
	GetResources() []IResource
	GetFeastDays() []IHolidays
	SortTasksByOrder()
}

// ProjectPlanID alias para el identificador de un plan de proyecto
type ProjectPlanID string

// ProjectPlanBase Contiene información de la planificación del proyecto
// Solo pongo etiquetas bson a aquellas por las que voy a buscar o de lo contrario no funciona
type ProjectPlanBase struct {
	// Identificador del plan de proyecto
	ID ProjectPlanID `json:"id" bson:"_id"`
	// Fecha planificada de comienzo del proyecto
	StartDate time.Time `json:"startDate"`
	// Fecha planificada de fin del proyecto
	EndDate time.Time `json:"endDate"`
	// Porcentaje real completado
	RealProgress uint `json:"realProgress"`
	// Porcentaje completado según lo planificado
	ExpectedProgress uint `json:"expectedProgress"`
	// Avance o retraso real en días
	RealProgressDays float64 `json:"realProgressDays"`
	// Fecha estimada de fin calculada en cada revisión en función de los días de avance o retraso y los días de fiesta
	EstimatedEndDate time.Time `json:"estimatedEndDate"`
	Workdays         uint      `json:"workDays"`
	// Indica si está archivado o no
	IsArchived bool `json:"archived"`
	// Fecha en la que el proyecto fue archivado
	ArchivedDate time.Time `json:"archivedDate"`
	// Tareas completadas
	CompleteTasks uint `json:"completeTasks" bson:",omitempty"`
	// Número Total de tareas
	TotalTasks uint `json:"totalTasks"`
}

// NewProjectPlan crea un nuevo plan de proyecto para poder ser planificado o revisado
func NewProjectPlan(id ProjectPlanID) *ProjectPlanBase {
	return &ProjectPlanBase{
		ID: id,
	}
}
