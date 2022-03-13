package gplan

import (
	"time"
)

// TaskID Alias para los ID de las tareas
type TaskID string

type ITask interface {
	Base() *TaskBase
	GetResource() IResource
	SetResource(resource IResource)
	GetBlocksTo() []ITaskDependency
	GetBlocksBy() []ITaskDependency
}

// TaskBase Contiene información de una tarea
type TaskBase struct {
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

// NewTaskBase crea una nueva tarea
func NewTaskBase(id TaskID, summary string, resourceType string, order int, duration uint) *TaskBase {
	return &TaskBase{
		ID:           id,
		Summary:      summary,
		ResourceType: resourceType,
		Order:        order,
		Duration:     duration,
	}
}

/*
// BlocksTo Crea una referencia de una tarea que bloquea a otra tarea
func BlocksTo(taskBlocks *TaskBase, taskBloked *TaskBase) {
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
*/
