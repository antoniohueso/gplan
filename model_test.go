package gplan_test

import (
	"sort"
	"time"

	"github.com/antoniohueso/gplan"
)

//---- Holidays

// Holidays contiene información de un rango de fechas de vacaciones o días de fiesta
type Holidays struct {
	// Fecha de vacaciones/fiesta desde
	From time.Time `json:"from"`
	// Fecha de vacaciones/fiesta Hasta
	To time.Time `json:"to"`
}

// NewHolidays crea un nuevo rango de fechas de vacaciones
func NewHolidaysBase(from time.Time, to time.Time) *Holidays {
	return &Holidays{
		From: from,
		To:   to,
	}
}

func (s *Holidays) GetFrom() time.Time {
	return s.From
}

func (s *Holidays) GetTo() time.Time {
	return s.To
}

//---- Resource

// Resource Contiene información de un recurso
type Resource struct {
	// Nombre del recurso
	ID gplan.ResourceID `json:"name"`
	// Descripción del recurso
	Description string `json:"description"`
	// Clasifica el tipo de recurso, Deberá coincidir con el tipo de tarea a asignarle.
	Type string `json:"type"`
	// Fecha desde la que estará disponible para poder asignarle una tarea
	AvailableFrom time.Time `json:"availableFrom"`
	// Siguiente fecha de disponibilidad. De uso interno para calcular la siguiente fecha de disponibilidad
	nextAvailableDate time.Time
	// Días de vacaciones del recurso
	Holidays []*Holidays
}

// NewResource crea un nuevo recurso
func NewResource(id gplan.ResourceID, description string, resourceType string, availableFrom time.Time, holidays []*Holidays) *Resource {

	if holidays == nil {
		holidays = []*Holidays{}
	}

	return &Resource{
		ID:                id,
		Description:       description,
		Type:              resourceType,
		AvailableFrom:     availableFrom,
		nextAvailableDate: availableFrom,
		Holidays:          holidays,
	}
}

func (s *Resource) GetID() gplan.ResourceID {
	return s.ID
}

func (s *Resource) GetDescription() string {
	return s.Description
}

func (s *Resource) GetType() string {
	return s.Type
}

func (s *Resource) GetAvailableFrom() time.Time {
	return s.AvailableFrom
}

func (s *Resource) SetAvailableFrom(t time.Time) {
	s.AvailableFrom = t
}

func (s *Resource) GetNextAvailableDate() time.Time {
	return s.nextAvailableDate
}

func (s *Resource) SetNextAvailableDate(t time.Time) {
	s.nextAvailableDate = t
}

func (s *Resource) GetHolidays() []gplan.Holidays {
	var slice = []gplan.Holidays{}

	for i := range s.Holidays {
		slice = append(slice, s.Holidays[i])
	}

	return slice
}

//---- TaskDependency

// TaskDependency contiene información de una dependencia de una tarea
type TaskDependency struct {
	TaskID gplan.TaskID `json:"taskId"`
}

// GetTaskID implementa TaskDependency
func (s *TaskDependency) GetTaskID() gplan.TaskID {
	return s.TaskID
}

// NewTaskDependency Crea un nuevo objeto TaskDependency
func NewTaskDependency(taskID gplan.TaskID) *TaskDependency {
	return &TaskDependency{TaskID: taskID}
}

//--- Task

// Task Contiene información de una tarea
type Task struct {
	// Clave por la que identificar una tarea
	ID gplan.TaskID `json:"id"`
	// Descripción de la tarea
	Summary string `json:"summary"`
	// Tipo de recurso que podrá ser asignado a esta tarea
	ResourceType string `json:"label"`
	// Número de orden de la tarea dentro de la lista de tareas
	Order uint `json:"order"`
	// Duración de la tarea en días
	Duration uint `json:"duration"`
	// Fecha planificada de comienzo de la tarea
	StartDate time.Time `json:"startDate"`
	// Fecha planificada de fin de la tarea
	EndDate time.Time `json:"endDate"`
	// Porcentaje real completado
	RealProgress uint `json:"realProgress"`
	// Porcentaje completado según lo planificado
	ExpectedProgress uint `json:"expectedProgress"`
	// Fecha real de finalización
	RealEndDate time.Time `json:"realEndDate"`
	// Recurso asignado
	Resource *Resource
	// IDs de las tareas a las que bloquea esta tarea
	BlocksTo []*TaskDependency
	// IDs de las tareas que bloquean a esta tarea
	BlocksBy []*TaskDependency
}

// NewTask crea una nueva tarea
func NewTask(id gplan.TaskID, summary string, resourceType string, order uint, duration uint) *Task {
	return &Task{
		ID:           id,
		Summary:      summary,
		ResourceType: resourceType,
		Order:        order,
		Duration:     duration,
	}
}

// NewTaskWithBlocks crea una nueva tarea con bloqueos
func NewTaskWithBlocks(id gplan.TaskID, summary string, resourceType string, order uint, duration uint, blocksTo []*TaskDependency, blocksBy []*TaskDependency) *Task {

	if blocksTo == nil {
		blocksTo = []*TaskDependency{}
	}

	if blocksBy == nil {
		blocksBy = []*TaskDependency{}
	}

	return &Task{
		ID:           id,
		Summary:      summary,
		ResourceType: resourceType,
		Order:        order,
		Duration:     duration,
		BlocksTo:     blocksTo,
		BlocksBy:     blocksBy,
	}
}

// BlocksTo Crea una referencia de una tarea que bloquea a otra tarea
func BlocksTo(taskBlocks *Task, taskBloked *Task) {
	if !hasTaskDependency(taskBloked.ID, taskBlocks.BlocksTo) {
		taskBlocks.BlocksTo = append(taskBlocks.BlocksTo, NewTaskDependency(taskBloked.ID))
	}
	if !hasTaskDependency(taskBlocks.ID, taskBloked.BlocksBy) {
		taskBloked.BlocksBy = append(taskBloked.BlocksBy, NewTaskDependency(taskBlocks.ID))
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

func (s *Task) GetID() gplan.TaskID {
	return s.ID
}

func (s *Task) GetSummary() string {
	return s.Summary
}

func (s *Task) GetResourceType() string {
	return s.ResourceType
}

func (s *Task) GetOrder() uint {
	return s.Order
}

func (s *Task) GetDuration() uint {
	return s.Duration
}

func (s *Task) GetStartDate() time.Time {
	return s.StartDate
}

func (s *Task) SetStartDate(t time.Time) {
	s.StartDate = t
}

func (s *Task) GetEndDate() time.Time {
	return s.EndDate
}

func (s *Task) SetEndDate(t time.Time) {
	s.EndDate = t
}

func (s *Task) GetRealEndDate() time.Time {
	return s.RealEndDate
}

func (s *Task) SetRealEndDate(t time.Time) {
	s.RealEndDate = t
}

func (s *Task) GetRealProgress() uint {
	return s.RealProgress
}

func (s *Task) SetRealProgress(n uint) {
	s.RealProgress = n
}

func (s *Task) GetExpectedProgress() uint {
	return s.ExpectedProgress
}

func (s *Task) SetExpectedProgress(n uint) {
	s.ExpectedProgress = n
}

func (s *Task) GetResource() gplan.Resource {
	return s.Resource
}

func (s *Task) SetResource(r gplan.Resource) {
	s.Resource = r.(*Resource)
}

func (s *Task) GetBlocksTo() []gplan.TaskDependency {
	var slice = []gplan.TaskDependency{}

	for i := range s.BlocksTo {
		slice = append(slice, s.BlocksTo[i])
	}

	return slice
}

func (s *Task) GetBlocksBy() []gplan.TaskDependency {
	var slice = []gplan.TaskDependency{}

	for i := range s.BlocksBy {
		slice = append(slice, s.BlocksBy[i])
	}

	return slice
}

//---- ProjectPlan

// ProjectPlan Contiene información de la planificación del proyecto
// Solo pongo etiquetas bson a aquellas por las que voy a buscar o de lo contrario no funciona
type ProjectPlan struct {
	// Identificador del plan de proyecto
	ID gplan.ProjectPlanID `json:"id" bson:"_id"`
	// Fecha planificada de comienzo del proyecto
	StartDate time.Time `json:"startDate"`
	// Fecha planificada de fin del proyecto
	EndDate time.Time `json:"endDate"`
	// Porcentaje real completado
	RealProgress uint `json:"realProgress"`
	// Porcentaje completado según lo planificado
	ExpectedProgress uint `json:"expectedProgress"`
	// Avance o retraso real en días
	RealProgressDays float64 `json:"realProgressDays"`
	// Fecha estimada de fin calculada en cada revisión en función de los días de avance o retraso y los días de fiesta
	EstimatedEndDate time.Time `json:"estimatedEndDate"`
	Workdays         uint      `json:"workDays"`
	// Tareas completadas
	CompleteTasks uint `json:"completeTasks" bson:",omitempty"`
	// Número Total de tareas
	TotalTasks uint `json:"totalTasks"`
	// Indica si está archivado o no
	Archived bool `json:"archived" bson:"archived"`
	// Fecha en la que el proyecto fue archivado
	ArchivedDate time.Time `json:"archivedDate"`
	// Lista de tareas
	Tasks []*Task
	// Lista de recursos
	Resurces []*Resource
	// Festivos
	FeastDays []*Holidays
}

// NewProjectPlan crea un nuevo plan de proyecto para poder ser planificado o revisado
func NewProjectPlan(id gplan.ProjectPlanID) *ProjectPlan {
	return &ProjectPlan{
		ID: id,
	}
}

func (s *ProjectPlan) GetID() gplan.ProjectPlanID {
	return s.ID
}

func (s *ProjectPlan) GetStartDate() time.Time {
	return s.StartDate
}

func (s *ProjectPlan) SetStartDate(t time.Time) {
	s.StartDate = t
}

func (s *ProjectPlan) GetEndDate() time.Time {
	return s.EndDate
}

func (s *ProjectPlan) SetEndDate(t time.Time) {
	s.EndDate = t
}

func (s *ProjectPlan) GetRealProgress() uint {
	return s.RealProgress
}

func (s *ProjectPlan) SetRealProgress(n uint) {
	s.RealProgress = n
}

func (s *ProjectPlan) GetExpectedProgress() uint {
	return s.ExpectedProgress
}

func (s *ProjectPlan) SetExpectedProgress(n uint) {
	s.ExpectedProgress = n
}

func (s *ProjectPlan) GetEstimatedEndDate() time.Time {
	return s.EstimatedEndDate
}

func (s *ProjectPlan) SetEstimatedEndDate(t time.Time) {
	s.EstimatedEndDate = t
}

func (s *ProjectPlan) GetWorkdays() uint {
	return s.Workdays
}

func (s *ProjectPlan) SetWorkdays(n uint) {
	s.Workdays = n
}

func (s *ProjectPlan) GetTotalTasks() uint {
	return s.TotalTasks
}

func (s *ProjectPlan) SetTotalTasks(n uint) {
	s.TotalTasks = n
}

func (s *ProjectPlan) GetCompleteTasks() uint {
	return s.CompleteTasks
}

func (s *ProjectPlan) SetCompleteTasks(n uint) {
	s.CompleteTasks = n
}

func (s *ProjectPlan) IsArchived() bool {
	return s.Archived
}

func (s *ProjectPlan) GetArchivedDate() time.Time {
	return s.ArchivedDate
}

func (s *ProjectPlan) GetTasks() []gplan.Task {
	var slice = []gplan.Task{}

	for i := range s.Tasks {
		slice = append(slice, s.Tasks[i])
	}

	return slice
}

func (s *ProjectPlan) GetResources() []gplan.Resource {
	var slice = []gplan.Resource{}

	for i := range s.Resurces {
		slice = append(slice, s.Resurces[i])
	}

	return slice
}

func (s *ProjectPlan) GetFeastDays() []gplan.Holidays {
	var slice = []gplan.Holidays{}

	for i := range s.FeastDays {
		slice = append(slice, s.FeastDays[i])
	}

	return slice
}

// SortTasksByOrder Implementa SortTasksByOrder
func (s *ProjectPlan) SortTasksByOrder() {
	// Ordena las tareas por número de orden para poder planificarlas
	sort.Slice(s.Tasks, func(i, j int) bool {
		return s.Tasks[i].Order < s.Tasks[j].Order
	})
}
