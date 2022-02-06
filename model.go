package gplan

import "time"

// IHolidays interfaz para implementar el acceso a la información necesaria sobre rango de fechas de vacaciones o días de fiesta
type IHolidays interface {
	GetFrom() time.Time
	GetTo() time.Time
}

// ResourceID alias del ID único de un recurso
type ResourceID string

// IResource interfaz para implementar el acceso a la información necesaria de un recurso para poder planificar una tarea
type IResource interface {
	// Identificador único del recurso
	GetID() ResourceID
	// Clasifica el tipo de recurso, Deberá coincidir con el tipo de tarea a asignarle.
	GetType() string
	// Fecha desde la que estará disponible para poder asignarle una tarea
	GetAvailableFrom() time.Time
	// Conjunto de días de vacaciones
	GetHolidays() []IHolidays
	// Siguiente fecha de disponibilidad. De uso interno para calcular la siguiente fecha de disponibilidad
	GetNextAvailableDate() time.Time
	SetNextAvailableDate(date time.Time)
}

// TaskID Alias para los ID de las tareas
type TaskID string

// ITaskDependency interfaz para implementar dependencias entre las tareas de la planificación
type ITaskDependency interface {
	GetTaskID() TaskID
}

// ITask interfaz para implementar el acceso a la información de una tarea a planificar
type ITask interface {
	// Clave por la que identificar una tarea
	GetID() TaskID
	// Tipo de recurso que podrá ser asignado a esta tarea
	GetResourceType() string
	// Número de orden de la tarea dentro de la lista de tareas
	GetOrder() int
	// Duración de la tarea en días
	GetDuration() int
	// Lista de tareas a las que bloquea y no se pueden empezar hasta que estuviera completada
	GetBlocksTo() []ITaskDependency
	// Lista de tareas que bloquean a esta tarea y no podría empezarse hasta que estuvieran completadas
	GetBlocksBy() []ITaskDependency
	// Recurso asignado
	GetResource() IResource
	SetResource(resource IResource)
	// Fecha planificada de comienzo de la tarea
	GetStartDate() time.Time
	SetStartDate(date time.Time)
	// Fecha planificada de fin de la tarea
	GetEndDate() time.Time
	SetEndDate(date time.Time)
	// Porcentaje real completado
	GetComplete() int
	SetComplete(n int)
	// Porcentaje completado según lo planificado
	GetEstimatedComplete() int
	SetEstimatedComplete(n int)
	// Fecha en la que realmente finalizó
	GetRealEndDate() time.Time
	SetRealEndDate(date time.Time)
}

// ProjectPlanID alias para el identificador de un plan de proyecto
type ProjectPlanID string

// IProjectPlan interfaz para implementar el acceso a la información necesaria para poder crear o revisar una planificación
type IProjectPlan interface {
	// Identificador del plan de proyecto
	GetID() ProjectPlanID
	// Fecha planificada de comienzo del proyecto
	GetStartDate() time.Time
	SetStartDate(date time.Time)
	// Fecha planificada de fin del proyecto
	GetEndDate() time.Time
	SetEndDate(date time.Time)
	// Lista de tareas planificadas en el proyecto
	GetTasks() []ITask
	// Lista de recursos
	GetResources() []IResource
	// Días de fiesta
	GetFeastDays() []IHolidays
	// Porcentaje real completado
	GetComplete() int
	SetComplete(n int)
	// Porcentaje completado según lo planificado
	GetEstimatedComplete() int
	SetEstimatedComplete(n int)
	// Avance o retraso real en días
	GetRealAdvancedOrDelayed() float64
	SetRealAdvancedOrDelayed(n float64)
	// Si está archivado o no
	IsArchived() bool
	SetArchived(archived bool)
	// Fecha en la que el proyecto fue archivado
	GetArchivedDate() time.Time
	SetArchivedDate(date time.Time)
	SortTasksByOrder()
}

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
