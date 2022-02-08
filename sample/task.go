package sample

import (
	"time"

	"github.com/antoniohueso/gplan"
)

// Task Contiene información de una tarea
type Task struct {
	// Clave por la que identificar una tarea
	ID gplan.TaskID
	// Descripción de la tarea
	Summary string
	// Tipo de recurso que podrá ser asignado a esta tarea
	ResourceType string
	// Número de orden de la tarea dentro de la lista de tareas
	Order int
	// Duración de la tarea en días
	Duration int
	// Lista de tareas a las que bloquea y no se pueden empezar hasta que estuviera completada
	BlocksTo []*TaskDependency
	// Lista de tareas que bloquean a esta tarea y no podría empezarse hasta que estuvieran completadas
	BlocksBy []*TaskDependency
	// Recurso asignado
	Resource gplan.Resource
	// Fecha planificada de comienzo de la tarea
	StartDate time.Time
	// Fecha planificada de fin de la tarea
	EndDate time.Time
	// Porcentaje real completado
	Complete int
	// Porcentaje completado según lo planificado
	EstimatedComplete int
	// Fecha real de finalización
	RealEndDate time.Time
}

// NewTask crea una nueva tarea sin bloqueos
func NewTask(id gplan.TaskID, summary string, resourceType string, order int, duration int) *Task {
	return NewTaskWithBlocks(id, summary, resourceType, order, duration, nil, nil)
}

// NewTaskWithBlocks crea una nueva tarea con bloqueos
func NewTaskWithBlocks(id gplan.TaskID, summary string, resourceType string, order int, duration int, blocksTo []*TaskDependency, blocksBy []*TaskDependency) *Task {

	if blocksTo == nil {
		blocksTo = []*TaskDependency{}
	}

	if blocksBy == nil {
		blocksBy = []*TaskDependency{}
	}

	return &Task{
		ID:           id,
		Summary:      summary,
		ResourceType: resourceType,
		Order:        order,
		Duration:     duration,
		BlocksTo:     blocksTo,
		BlocksBy:     blocksBy,
	}
}

// GetID Getter de ID
func (t Task) GetID() gplan.TaskID {
	return t.ID
}

// GetResourceType Getter de ResourceType
func (t Task) GetResourceType() string {
	return t.ResourceType
}

// GetOrder Getter de Order
func (t Task) GetOrder() int {
	return t.Order
}

// GetDuration Getter de Duration
func (t Task) GetDuration() int {
	return t.Duration
}

// GetBlocksTo devuelve un nuevo array de TaskDependency
func (t Task) GetBlocksTo() []gplan.TaskDependency {
	newArr := make([]gplan.TaskDependency, len(t.BlocksTo))
	for i := range t.BlocksTo {
		newArr[i] = t.BlocksTo[i]
	}
	return newArr
}

// GetBlocksBy devuelve un nuevo array de TaskDependency
func (t Task) GetBlocksBy() []gplan.TaskDependency {
	newArr := make([]gplan.TaskDependency, len(t.BlocksBy))
	for i := range t.BlocksBy {
		newArr[i] = t.BlocksBy[i]
	}
	return newArr
}

// GetResource Getter de Resource
func (t Task) GetResource() gplan.Resource {
	return t.Resource
}

// SetResource Setter de Resource
func (t *Task) SetResource(resource gplan.Resource) {
	t.Resource = resource
}

// GetStartDate Getter de StartDate
func (t Task) GetStartDate() time.Time {
	return t.StartDate
}

// SetStartDate Setter de StartDate
func (t *Task) SetStartDate(date time.Time) {
	t.StartDate = date
}

// GetEndDate Getter de EndDate
func (t Task) GetEndDate() time.Time {
	return t.EndDate
}

// SetEndDate Setter de EndDate
func (t *Task) SetEndDate(date time.Time) {
	t.EndDate = date
}

// GetComplete Getter de Complete
func (t Task) GetComplete() int {
	return t.Complete
}

// SetComplete Setter de Complete
func (t *Task) SetComplete(n int) {
	t.Complete = n
}

// GetEstimatedComplete Getter de EstimatedComplete
func (t Task) GetEstimatedComplete() int {
	return t.EstimatedComplete
}

// SetEstimatedComplete Setter de EstimatedComplete
func (t *Task) SetEstimatedComplete(n int) {
	t.EstimatedComplete = n
}

// GetRealEndDate Getter de RealEndDate
func (t Task) GetRealEndDate() time.Time {
	return t.RealEndDate
}

// SetRealEndDate Setter de RealEndDate
func (t *Task) SetRealEndDate(date time.Time) {
	t.RealEndDate = date
}

// BlocksTo Crea una referencia de una tarea que bloquea a otra tarea
func BlocksTo(taskBlocks *Task, taskBloked *Task) {
	if !hasTaskDependency(taskBloked.ID, taskBlocks.BlocksTo) {
		taskBlocks.BlocksTo = append(taskBlocks.BlocksTo, NewTaskDependency(taskBloked.ID))
	}
	if !hasTaskDependency(taskBlocks.ID, taskBloked.BlocksBy) {
		taskBloked.BlocksBy = append(taskBloked.BlocksBy, NewTaskDependency(taskBlocks.ID))
	}
}

// taskIDExists Devuelve true si el ID de una tarea existe en el arrays de IDs
func hasTaskDependency(taskID gplan.TaskID, taskDependencies []*TaskDependency) bool {
	for _, td := range taskDependencies {
		if td.TaskID == taskID {
			return true
		}
	}
	return false
}
