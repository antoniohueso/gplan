package sample

import "github.com/antoniohueso/gplan"

// TaskDependency Contiene información de una dependencia entre tareas
type TaskDependency struct {
	*gplan.TaskDependencyBase
	Summary string
	KeyURL  string
}

func (s *TaskDependency) Base() *gplan.TaskDependencyBase {
	return s.TaskDependencyBase
}

// NewTaskDependency Crea un nuevo objeto TaskDependency
func NewTaskDependency(taskID gplan.TaskID, summary string, keyUrl string) *TaskDependency {
	return &TaskDependency{
		TaskDependencyBase: &gplan.TaskDependencyBase{TaskID: taskID},
		Summary:            summary,
		KeyURL:             keyUrl}
}

func (s TaskDependency) GetTaskDependency() *gplan.TaskDependencyBase {
	return s.TaskDependencyBase
}
