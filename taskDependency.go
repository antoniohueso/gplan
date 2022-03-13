package gplan

type ITaskDependency interface {
	GetTaskDependency() *TaskDependency
}

// TaskDependency Contiene información de una dependencia entre tareas
type TaskDependency struct {
	TaskID TaskID `json:"taskId"`
}

// NewTaskDependency Crea un nuevo objeto TaskDependency
func NewTaskDependency(taskID TaskID) *TaskDependency {
	return &TaskDependency{TaskID: taskID}
}

// GetTaskDependency implementa ITaskDependency
func (s *TaskDependency) GetTaskDependency() *TaskDependency {
	return s
}
