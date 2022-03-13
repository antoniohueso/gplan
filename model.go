package gplan

import "time"

// ScheduledTaskInfo datos de la planificación de una tarea
type scheduledTaskInfo struct {
	// Recurso asignado
	Resource IResource
	// Fecha planificada de comienzo de la tarea
	StartDate time.Time
	// Fecha planificada de fin de la tarea
	EndDate time.Time
}

// Error Contiene información de un Message que se haya podido producir al crear o revisar la planificación
type Error struct {
	Message error
	Tasks   []TaskID
}
