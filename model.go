package gplan

import "time"

// Holidays interface a implementar para los rangos de vacaciones o días de fiesta
type Holidays interface {
	// fecha desde
	GetFrom() time.Time
	// fecha Hasta
	GetTo() time.Time
}

// ResourceID alias del ID único de un recurso
type ResourceID string

// Resource interface a implementar para definir un recurso disponible para una planificación
type Resource interface {
	// Nombre del recurso
	GetID() ResourceID
	// Descripción del recurso
	GetDescription() string
	// Clasifica el tipo de recurso, Deberá coincidir con el tipo de tarea a asignarle.
	GetType() string
	// Fecha desde la que estará disponible para poder asignarle una tarea
	GetAvailableFrom() time.Time
	SetAvailableFrom(time.Time)
	// Siguiente fecha de disponibilidad. De uso interno para calcular la siguiente fecha de disponibilidad
	GetNextAvailableDate() time.Time
	SetNextAvailableDate(time.Time)
	// Vacaciones del recurso
	GetHolidays() []Holidays
}

// TaskDependency interface a implementar para definir una dependencia entre tareas
type TaskDependency interface {
	// ID de la tarea de la que depende
	GetTaskID() TaskID
}

// TaskID Alias para los ID de las tareas
type TaskID string

// Task interface que implementa una tarea de la planificación
type Task interface {
	// Clave por la que identificar una tarea
	GetID() TaskID
	// Descripción de la tarea
	GetSummary() string
	// Tipo de recurso que podrá ser asignado a esta tarea
	GetResourceType() string
	// Número de orden de la tarea dentro de la lista de tareas
	GetOrder() uint
	// Duración de la tarea en días
	GetDuration() uint
	// Fecha planificada de comienzo de la tarea
	GetStartDate() time.Time
	SetStartDate(time.Time)
	// Fecha planificada de fin de la tarea
	GetEndDate() time.Time
	SetEndDate(time.Time)
	// Porcentaje real completado
	GetRealProgress() uint
	SetRealProgress(uint)
	// Porcentaje completado según lo planificado
	GetExpectedProgress() uint
	SetExpectedProgress(uint)
	// Fecha real de finalización
	GetRealEndDate() time.Time
	SetRealEndDate(time.Time)
	// Recurso asignado
	GetResource() Resource
	SetResource(Resource)
	// IDs de las tareas a las que bloquea esta tarea
	GetBlocksTo() []TaskDependency
	// IDs de las tareas que bloquean a esta tarea
	GetBlocksBy() []TaskDependency
}

// ProjectPlanID alias para el identificador de un plan de proyecto
type ProjectPlanID string

// ProjectPlan interface a implementar para  Contiene información de la planificación del proyecto
type ProjectPlan interface {
	// Identificador del plan de proyecto
	GetID() ProjectPlanID
	// Fecha planificada de comienzo del proyecto
	GetStartDate() time.Time
	SetStartDate(time.Time)
	// Fecha planificada de fin del proyecto
	GetEndDate() time.Time
	SetEndDate(time.Time)
	// Porcentaje real completado
	GetRealProgress() uint
	SetRealProgress(uint)
	// Porcentaje completado según lo planificado
	GetExpectedProgress() uint
	SetExpectedProgress(uint)
	// Avance o retraso real en días
	GetRealProgressDays() float64
	SetRealProgressDays(float64)
	// Fecha estimada de fin calculada en cada revisión en función de los días de avance o retraso y los días de fiesta
	GetEstimatedEndDate() time.Time
	SetEstimatedEndDate(time.Time)
	// Jornadas de trabajo que tiene la planificación
	GetWorkdays() uint
	SetWorkdays(uint)
	// Tareas completadas
	GetCompleteTasks() uint
	SetCompleteTasks(uint)
	// Número Total de tareas
	GetTotalTasks() uint
	SetTotalTasks(uint)
	// Indica si está archivado o no
	IsArchived() bool
	// Fecha en la que el proyecto fue archivado
	GetArchivedDate() time.Time
	// Tareas de la planificación
	GetTasks() []Task
	// Recursos disponibles para la planificación
	GetResources() []Resource
	// Festivos
	GetFeastDays() []Holidays
	// Método que ordena las tareas por el campo Orden
	SortTasksByOrder()
}

// scheduledTaskInfo datos de la planificación de una tarea
type scheduledTaskInfo struct {
	// Recurso asignado
	Resource Resource
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
