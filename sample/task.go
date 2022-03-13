package sample

import (
	"github.com/antoniohueso/gplan"
)

// Task Contiene información de una tarea
type Task struct {
	*gplan.TaskBase
	// Recurso asignado - No se almacena en BD
	Resource *Resource `json:"resource" bson:",omitempty"`
	// Lista de tareas a las que bloquea y no se pueden empezar hasta que estuviera completada
	BlocksTo []*TaskDependency
	// Lista de tareas que bloquean a esta tarea y no podría empezarse hasta que estuvieran completadas
	BlocksBy []*TaskDependency
}

func (s *Task) Base() *gplan.TaskBase {
	return s.TaskBase
}

func (s *Task) GetResource() gplan.IResource {
	return s.Resource
}

func (s *Task) SetResource(resource gplan.IResource) {
	s.Resource = resource.(*Resource)
}

// GetBlocksTo devuelve un nuevo array de TaskDependency
func (s *Task) GetBlocksTo() []gplan.ITaskDependency {
	newArr := make([]gplan.ITaskDependency, len(s.BlocksTo))
	for i := range s.BlocksTo {
		newArr[i] = s.BlocksTo[i]
	}
	return newArr
}

// GetBlocksBy devuelve un nuevo array de TaskDependency
func (s *Task) GetBlocksBy() []gplan.ITaskDependency {
	newArr := make([]gplan.ITaskDependency, len(s.BlocksBy))
	for i := range s.BlocksBy {
		newArr[i] = s.BlocksBy[i]
	}
	return newArr
}

// NewTask crea una nueva tarea sin bloqueos
func NewTask(id gplan.TaskID, summary string, resourceType string, order int, duration uint) *Task {
	return NewTaskWithBlocks(id, summary, resourceType, order, duration, nil, nil)
}

// NewTaskWithBlocks crea una nueva tarea con bloqueos
func NewTaskWithBlocks(id gplan.TaskID, summary string, resourceType string, order int, duration uint, blocksTo []*TaskDependency, blocksBy []*TaskDependency) *Task {

	if blocksTo == nil {
		blocksTo = []*TaskDependency{}
	}

	if blocksBy == nil {
		blocksBy = []*TaskDependency{}
	}

	return &Task{
		TaskBase: gplan.NewTaskBase(id, summary, resourceType, order, duration),
		BlocksTo: blocksTo,
		BlocksBy: blocksBy,
	}
}

// BlocksTo Crea una referencia de una tarea que bloquea a otra tarea
func BlocksTo(taskBlocks *Task, taskBloked *Task) {
	if !hasTaskDependency(taskBloked.ID, taskBlocks.BlocksTo) {
		taskBlocks.BlocksTo = append(taskBlocks.BlocksTo, NewTaskDependency(taskBloked.ID, taskBloked.Summary, "KeyURL"))
	}
	if !hasTaskDependency(taskBlocks.ID, taskBloked.BlocksBy) {
		taskBloked.BlocksBy = append(taskBloked.BlocksBy, NewTaskDependency(taskBlocks.ID, taskBlocks.Summary, "KeyURL"))
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
