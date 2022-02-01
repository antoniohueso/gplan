package defmodel

import "com.github.antoniohueso/gplan"

// TaskDependency Contiene información de una dependencia entre tareas
type TaskDependency struct {
	TaskID gplan.TaskID
}

// NewTaskDependency Crea un nuevo objeto TaskDependency
func NewTaskDependency(taskID gplan.TaskID) *TaskDependency {
	return &TaskDependency{TaskID: taskID}
}

// GetTaskID Getter de ID
func (td TaskDependency) GetTaskID() gplan.TaskID {
	return td.TaskID
}

// ArrayOfTaskDependencies implementa gplan.ArrayOfTaskDependencies
type ArrayOfTaskDependencies []*TaskDependency

// Iterable Implementa Iterable
func (a ArrayOfTaskDependencies) Iterable() []gplan.TaskDependency {
	newArr := make([]gplan.TaskDependency, len(a))
	for i := range a {
		newArr[i] = a[i]
	}
	return newArr
}
