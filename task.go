package gplan

// Dependency Contiene la información de una dependencia entre tareas
type Dependency struct {
	TaskID TaskID
}

// GetTaskID Getter TaskID
func (d *Dependency) GetTaskID() TaskID {
	return d.TaskID
}
