package gplan

type ITaskDependency interface {
	Base() *TaskDependencyBase
}

// TaskDependencyBase Contiene información de una dependencia entre tareas
type TaskDependencyBase struct {
	TaskID TaskID `json:"taskId"`
}

// NewTaskDependency Crea un nuevo objeto TaskDependency
func NewTaskDependency(taskID TaskID) *TaskDependencyBase {
	return &TaskDependencyBase{TaskID: taskID}
}
