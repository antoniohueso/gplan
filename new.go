package gplan

import (
	"sort"
	"time"

	"github.com/antoniohueso/gplan/dateutil"
)

// Planning Crea una nueva planificación a partir de unas tareas, recursos, vacaciones y fecha de comienzo.
func Planning(startDate time.Time, plan ProjectPlan) *Error {

	// Hay que ordenarlas antes de que asihne las tareas ya que las extrae a otro array distinto
	plan.SortTasksByOrder()

	var (
		tasks      = plan.GetTasks()
		resources  = plan.GetResources()
		feastDays  = plan.GetFeastDays()
		tasksIndex map[TaskID]Task
		err        *Error
	)

	tasksIndex, err = validateTasks(tasks, resources)
	if err != nil {
		return err
	}

	// Si la fecha de disponibilidad ldel recurso es menor que la fecha en la que debe comenzar el proyecto sele pone la fecha en la que debe comenzar el proyecto
	// para que no haya ninguna tarea que comeice antes
	for _, resource := range resources {
		if dateutil.IsLt(resource.GetNextAvailableDate(), startDate) {
			resource.SetNextAvailableDate(startDate)
		}
	}

	if feastDays == nil {
		feastDays = []Holidays{}
	}

	// Planifica las tareas
	for _, task := range tasks {

		err = assignTask(task, resources, feastDays, tasksIndex)
		if err != nil {
			return err
		}
	}

	plan.SetStartDate(tasks[0].GetStartDate())
	plan.SetEndDate(tasks[0].GetEndDate())

	// Busca la fecha de comienzo y de fin del proyecto que será la menor fecha de inicio de una tarea y la mayor
	// fecha de fin del proyecto
	for _, task := range tasks {
		if dateutil.IsLt(task.GetStartDate(), plan.GetStartDate()) {
			plan.SetStartDate(task.GetStartDate())
		}
		if dateutil.IsGt(task.GetEndDate(), plan.GetEndDate()) {
			plan.SetEndDate(task.GetEndDate())
		}
	}

	plan.SetWorkdays(CalculateLaborableDays(plan.GetStartDate(), plan.GetEndDate(), feastDays))
	plan.SetTotalTasks(uint(len(tasks)))

	return nil
}

// validateTasks comprueba que las tareas tengan el formato correcto para poder realizar la planificación
func validateTasks(aTasks []Task, resources []Resource) (map[TaskID]Task, *Error) {
	var (
		taskIDSErrors []TaskID
		tasks         = aTasks
	)

	// La lista de tareas no puede estar vacía
	if len(tasks) == 0 {
		return nil, newTextError("la lista de tareas a planificar está vacía")
	}

	// La lista de recursos no puede estar vacía
	if len(resources) == 0 {
		return nil, newTextError("la lista de recursos a asignar está vacía")
	}

	// No puede haber tareas con una duración menor que 1
	for _, task := range tasks {
		if task.GetDuration() < 1 {
			taskIDSErrors = append(taskIDSErrors, task.GetID())
		}
	}

	if len(taskIDSErrors) > 0 {
		return nil, newError("las siguientes tareas tienen la duración inferior a un día", taskIDSErrors)
	}

	// No puede haber tareas con una orden menor que 1
	for _, task := range tasks {
		if task.GetOrder() < 1 {
			taskIDSErrors = append(taskIDSErrors, task.GetID())
		}
	}

	if len(taskIDSErrors) > 0 {
		return nil, newError("las siguientes tareas tienen un orden inferior a 1", taskIDSErrors)
	}

	// Ha de haber tareas de tipo 'ResourceType' para los tipos de recurso de tipo 'Type' y viceversa
	var (
		typeOfTasks     = make(map[string]bool)
		typeOfResources = make(map[string]bool)
	)

	// Crea el índice de tipos de tarea
	for _, task := range tasks {
		typeOfTasks[task.GetResourceType()] = true
	}

	// Crea el índice de tipos de resource
	for _, resource := range resources {
		typeOfResources[resource.GetType()] = true
	}

	// Tiene que haber tareas para los tipos de recurso que llegan
	for _, resource := range resources {
		if _, exist := typeOfTasks[resource.GetType()]; !exist {
			return nil, newTextError("no existen tareas para el recurso %s de tipo %s", resource.GetID(), resource.GetType())
		}
	}

	// Tiene que haber recursos para las tareas del tipo que llevan asociado
	for _, task := range tasks {
		if _, exist := typeOfResources[task.GetResourceType()]; !exist {
			taskIDSErrors = append(taskIDSErrors, task.GetID())
		}
	}
	if len(taskIDSErrors) > 0 {
		return nil, newError("no hay recursos para los tipos de estas tareas", taskIDSErrors)
	}

	// Crea el índice de tareas para poder saber a qué tarea corresponde un tareaID
	var tasksIndex = make(map[TaskID]Task, len(tasks))
	for _, task := range tasks {
		tasksIndex[task.GetID()] = task
	}

	var tasksNotExists = map[TaskID]bool{}

	// No puede haber tareas en blocksBy ni blocksTo que no estén en la lista de tareas
	for _, task := range tasks {

		// Revisa las tareas a las que bloquea
		for _, dep := range task.GetBlocksTo() {
			if _, exist := tasksIndex[dep.GetTaskID()]; !exist {
				tasksNotExists[dep.GetTaskID()] = true
			}
		}

		// Revisa las tareas que le bloquean
		for _, dep := range task.GetBlocksBy() {
			if _, exist := tasksIndex[dep.GetTaskID()]; !exist {
				tasksNotExists[dep.GetTaskID()] = true
			}
		}
	}

	for taskID := range tasksNotExists {
		taskIDSErrors = append(taskIDSErrors, taskID)
	}

	if len(taskIDSErrors) > 0 {
		return nil, newError("hay tareas bloquedas o que bloquean a otras que no existen en la lista de tareas", taskIDSErrors)
	}

	// No puede haber referencias circulares, es decir, tareas que se bloqueen a sí mismas
	for _, task := range tasks {
		err := checkForCircularDependencies(task, tasksIndex, []TaskID{})
		if err != nil {
			return nil, err
		}
	}

	// No puede haber tareas con un orden superior que bloqueen a otras con un número de orden inferior o igual
	for _, task := range tasks {
		for _, dep := range task.GetBlocksTo() {
			taskBlocked := tasksIndex[dep.GetTaskID()]
			if taskBlocked != nil {
				if task.GetOrder() >= taskBlocked.GetOrder() {
					taskIDSErrors = append(taskIDSErrors, dep.GetTaskID())
				}
			}
		}
	}

	if len(taskIDSErrors) > 0 {
		return nil, newError("tareas bloqueadas por tareas con un orden superior", taskIDSErrors)
	}

	return tasksIndex, nil
}

// checkForCircularDependencies Chequea que no haya referencias circulares en los bloqueos de las tareas, es decir que
// una tarea se bloquee a sí misma.
// Por ejemplo: A → B → C → A ...
func checkForCircularDependencies(task Task, tasksIndex map[TaskID]Task, linkedblocksTo []TaskID) *Error {

	for _, tID := range linkedblocksTo {
		if tID == task.GetID() {
			return newTextError("%s contiene una referencia circular en su cadena de dependencias", tID)
		}
	}

	// Recorre las tareas a las que está bloqueando llamando recursivamente a la propia función
	for _, dep := range task.GetBlocksTo() {
		// Añade task.ID a la secuencia de bloqueos para detectar una posible referencia circular a ella misma en
		// dicho encadenamiento
		linkedblocksTo = append(linkedblocksTo, task.GetID())
		err := checkForCircularDependencies(tasksIndex[dep.GetTaskID()], tasksIndex, linkedblocksTo)
		if err != nil {
			return err
		}
	}
	return nil
}

// assignTask Asigna la tarea al recurso que en una simulación de planificación pueda acabarla antes
func assignTask(task Task, resources []Resource, feastDays []Holidays, tasksIndex map[TaskID]Task) *Error {

	var (
		err               *Error
		startDate         time.Time
		scheduledTaskInfo *scheduledTaskInfo
	)

	// Le asigna la fecha de comienzo
	startDate, err = getRealStartDate(task, tasksIndex)
	if err != nil {
		return err
	}

	task.SetStartDate(startDate)

	// Le asigna los datos de la planificación
	scheduledTaskInfo = bestScheduledTask(task, resources, feastDays)

	var resourceID = scheduledTaskInfo.Resource.GetID()
	task.SetResourceID(&resourceID)
	task.SetStartDate(scheduledTaskInfo.StartDate)
	task.SetEndDate(scheduledTaskInfo.EndDate)

	return nil
}

// getRealStartDate Retorna la fecha real de comienzo de una tarea teniendo en cuenta que si tiene bloqueos
// calcula la fecha mayor de las que bloquean a la tarea, ya que no se podrá empezar antes.
func getRealStartDate(task Task, tasksIndex map[TaskID]Task) (time.Time, *Error) {

	var blocksBy = task.GetBlocksBy()

	// Si no tiene tareas que la bloqueen retorna la misma fecha de comienzo
	if len(blocksBy) == 0 {
		return task.GetStartDate(), nil
	}

	var tasksBlocksBy []Task

	// Las tareas que la bloquean en una lista de tareas ordenadas primero por las que bloquean a otros deberían
	// estar ya planificadas, si no lo están damos error
	for _, dep := range blocksBy {

		taskBlocksBy := tasksIndex[dep.GetTaskID()]

		if taskBlocksBy.GetResourceID() == nil {
			// Esto debería ser muy improbable que se dé si el código funciona como debe...
			return time.Time{}, newTextError("La tarea %s no está planificada y bloquea a la tarea %s que está en planificación",
				dep.GetTaskID(), task.GetID())
		}

		tasksBlocksBy = append(tasksBlocksBy, taskBlocksBy)
	}

	// Ordena blocksBy por fecha de fin de tarea descendiente para encontrar la fecha mayor
	sort.Slice(tasksBlocksBy, func(i, j int) bool {
		return dateutil.IsGt(tasksBlocksBy[i].GetEndDate(), tasksBlocksBy[j].GetEndDate())
	})

	// Retorna la fecha mayor + 1 día ya que debe comenzar al día siguiente de la fecha de fin de la última tarea que la bloquean
	return tasksBlocksBy[0].GetEndDate().AddDate(0, 0, 1), nil

}

// bestScheduledTask Calcula la planificación de la tarea para cada recurso y retorna la que termine antes.
func bestScheduledTask(task Task, resources []Resource, feastDays []Holidays) *scheduledTaskInfo {

	var (
		availableResources []Resource
		scheduledTasks     []*scheduledTaskInfo
		bestScheduled      *scheduledTaskInfo
	)

	// Filtra los recursos del tipo de tarea
	for _, resource := range resources {
		if resource.GetType() == task.GetResourceType() {
			availableResources = append(availableResources, resource)
		}
	}

	// Calcula la planificación de la tarea para cada recurso
	for _, resource := range availableResources {
		scheduledTasks = append(scheduledTasks, scheduledTask(task, resource, feastDays))
	}
	// Busca la mejor fecha
	for _, sh := range scheduledTasks {
		if bestScheduled == nil {
			bestScheduled = sh
		} else {
			if dateutil.IsLt(sh.EndDate, bestScheduled.EndDate) {
				bestScheduled = sh
			}
		}
	}

	// Le pone la fecha siguiente fecha de disponibilidad al recurso asignado
	bestScheduled.Resource.SetNextAvailableDate(bestScheduled.EndDate.AddDate(0, 0, 1))

	return bestScheduled
}

// scheduledTask planifica una tarea para un recurso
func scheduledTask(task Task, resource Resource, feastDays []Holidays) *scheduledTaskInfo {

	var (
		holidaysAndFeastDays []Holidays
		realStartDate        time.Time
		duration             uint
		endDate              time.Time
	)

	// concatena días de fiesta y vacaciones del recurso
	holidaysAndFeastDays = append(holidaysAndFeastDays, resource.GetHolidays()...)
	holidaysAndFeastDays = append(holidaysAndFeastDays, feastDays...)

	// Si la fecha en la que debe comenzar la tarea es superior a la fecha en la que el recurso estaría disponible
	// Ponemos esa fecha como fecha en la que el recurso estaría disponible para el cálculo y si no se pone la fecha en
	// la que estaría disponible el recurso
	if dateutil.IsGt(task.GetStartDate(), resource.GetNextAvailableDate()) {
		realStartDate = task.GetStartDate()
	} else {
		realStartDate = resource.GetNextAvailableDate()
	}

	duration = task.GetDuration()
	endDate = realStartDate

	for {

		if IsLaborableDay(endDate, holidaysAndFeastDays) {
			// Si es el primer día que empieza a contar actualiza la fecha de comienzo, puede que aunque la fecha de
			// comienzo inicial sea hoy, hoy y mañana sean fiesta por lo que comenzaría dos días después
			if duration == task.GetDuration() {
				realStartDate = endDate
			}
			duration--
			if duration == 0 {
				break
			}
		}
		endDate = endDate.AddDate(0, 0, 1)
	}

	return &scheduledTaskInfo{
		Resource:  resource,
		StartDate: realStartDate,
		EndDate:   endDate,
	}

}
