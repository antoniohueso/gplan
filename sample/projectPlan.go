package sample

import (
	"sort"

	"github.com/antoniohueso/gplan"
)

// ProjectPlan Contiene información de la planificación del proyecto
// Solo pongo etiquetas bson a aquellas por las que voy a buscar o de lo contrario no funciona
type ProjectPlan struct {
	*gplan.ProjectPlanBase
	// Lista de tareas planificadas en el proyecto
	Tasks []*Task
	// Lista de recursos
	Resources []*Resource
	// Días de fiesta
	FeastDays []*Holidays
	// Filtro SQL
	filter string
}

func (s *ProjectPlan) Base() *gplan.ProjectPlanBase {
	return s.ProjectPlanBase
}

// GetTasks devuelve un nuevo array de ITask
func (s *ProjectPlan) GetTasks() []gplan.ITask {
	newArr := make([]gplan.ITask, len(s.Tasks))
	for i := range s.Tasks {
		newArr[i] = s.Tasks[i]
	}
	return newArr
}

// GetResources devuelve un nuevo array de IResource
func (s *ProjectPlan) GetResources() []gplan.IResource {
	newArr := make([]gplan.IResource, len(s.Resources))
	for i := range s.Resources {
		newArr[i] = s.Resources[i]
	}
	return newArr
}

// GetFeastDays devuelve un nuevo array de IHolidays
func (s *ProjectPlan) GetFeastDays() []gplan.IHolidays {
	newArr := make([]gplan.IHolidays, len(s.FeastDays))
	for i := range s.FeastDays {
		newArr[i] = s.FeastDays[i]
	}
	return newArr
}

// SortTasksByOrder Implementa SortTasksByOrder
func (s ProjectPlan) SortTasksByOrder() {
	// Ordena las tareas por número de orden para poder planificarlas
	sort.Slice(s.Tasks, func(i, j int) bool {
		return s.Tasks[i].Order < s.Tasks[j].Order
	})
}

// NewProjectPlan crea un nuevo plan de proyecto para poder ser planificado o revisado
func NewProjectPlan(id gplan.ProjectPlanID, filter string, tasks []*Task, resources []*Resource, feastDays []*Holidays) *ProjectPlan {
	return &ProjectPlan{
		ProjectPlanBase: gplan.NewProjectPlan(id),
		Tasks:           tasks,
		Resources:       resources,
		FeastDays:       feastDays,
	}
}
