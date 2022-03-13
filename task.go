package gplan

import (
	"time"
)

// TaskID Alias para los ID de las tareas
type TaskID string

type ITask interface {
	GetTask() *Task
}

// Task Contiene información de una tarea
type Task struct {
	// Clave por la que identificar una tarea
	ID TaskID `json:"id"`
	// Descripción de la tarea
	Summary string `json:"summary"`
	// Tipo de recurso que podrá ser asignado a esta tarea
	ResourceType string `json:"label"`
	// Número de orden de la tarea dentro de la lista de tareas
	Order int `json:"order"`
	// Duración de la tarea en días
	Duration uint `json:"duration"`
	// Lista de tareas a las que bloquea y no se pueden empezar hasta que estuviera completada
	BlocksTo []ITaskDependency `json:"blocksTo"`
	// Lista de tareas que bloquean a esta tarea y no podría empezarse hasta que estuvieran completadas
	BlocksBy []ITaskDependency `json:"blocksBy"`
	// Recurso asignado - No se almacena en BD
	Resource IResource `json:"resource" bson:",omitempty"`
	// Fecha planificada de comienzo de la tarea
	StartDate time.Time `json:"startDate"`
	// Fecha planificada de fin de la tarea
	EndDate time.Time `json:"endDate"`
	// Porcentaje real completado
	RealProgress uint `json:"realProgress"`
	// Porcentaje completado según lo planificado
	ExpectedProgress uint `json:"expectedProgress"`
	// Fecha real de finalización
	RealEndDate time.Time `json:"realEndDate"`
}

// NewTask crea una nueva tarea sin bloqueos
func NewTask(id TaskID, issueID uint64, summary string, resourceType string, order int, duration uint) *Task {
	return NewTaskWithBlocks(id, issueID, summary, resourceType, order, duration, nil, nil)
}

// NewTaskWithBlocks crea una nueva tarea con bloqueos
func NewTaskWithBlocks(id TaskID, issueID uint64, summary string, resourceType string, order int, duration uint, blocksTo []ITaskDependency, blocksBy []ITaskDependency) *Task {

	if blocksTo == nil {
		blocksTo = []ITaskDependency{}
	}

	if blocksBy == nil {
		blocksBy = []ITaskDependency{}
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

// GetTask implementa ITask
func (s *Task) GetTask() *Task {
	return s
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

// hasTaskDependency Devuelve true si el ID de una tarea existe en el arrays de dependencias
func hasTaskDependency(taskID TaskID, taskDependencies []ITaskDependency) bool {
	for _, td := range taskDependencies {
		if td.GetTaskDependency().TaskID == taskID {
			return true
		}
	}
	return false
}
