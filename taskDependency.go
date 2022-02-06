package gplan

// TaskDependency Contiene información de una dependencia entre tareas
type TaskDependency struct {
	TaskID TaskID
}

// NewTaskDependency Crea un nuevo objeto TaskDependency
func NewTaskDependency(taskID TaskID) *TaskDependency {
	return &TaskDependency{TaskID: taskID}
}

// GetTaskID Getter de ID
func (td TaskDependency) GetTaskID() TaskID {
	return td.TaskID
}
