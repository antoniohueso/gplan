package gplan

import (
	"sort"
	"time"

	"github.com/antoniohueso/gplan/dateutil"
)

// Planning Crea una nueva planificación a partir de unas tareas, recursos, vacaciones y fecha de comienzo.
func Planning(startDate time.Time, projectPlan IProjectPlan) *Error {

	var plan *ProjectPlan = projectPlan.GetProjectPlan()

	// Ordena las tareas por número de orden para poder planificarlas
	sort.Slice(plan.Tasks, func(i, j int) bool {
		return plan.Tasks[i].GetTask().Order < plan.Tasks[j].GetTask().Order
	})

	var (
		tasks      []ITask     = plan.Tasks
		resources  []IResource = plan.Resources
		feastDays  []IHolidays = plan.FeastDays
		tasksIndex map[TaskID]ITask
		err        *Error
	)

	tasksIndex, err = validateTasks(tasks, resources)
	if err != nil {
		return err
	}

	// Si la fecha de disponibilidad ldel recurso es menor que la fecha en la que debe comenzar el proyecto sele pone la fecha en la que debe comenzar el proyecto
	// para que no haya ninguna tarea que comeice antes
	for _, resource := range resources {
		if dateutil.IsLt(resource.GetResource().nextAvailableDate, startDate) {
			resource.GetResource().nextAvailableDate = startDate
		}
	}

	if feastDays == nil {
		feastDays = []IHolidays{}
	}

	// Planifica las tareas
	for _, task := range tasks {

		err = assignTask(task, resources, feastDays, tasksIndex)
		if err != nil {
			return err
		}
	}

	plan.StartDate = tasks[0].GetTask().StartDate
	plan.EndDate = tasks[0].GetTask().EndDate

	// Busca la fecha de comienzo y de fin del proyecto que será la menor fecha de inicio de una tarea y la mayor
	// fecha de fin del proyecto
	for _, task := range tasks {
		if dateutil.IsLt(task.GetTask().StartDate, plan.StartDate) {
			plan.StartDate = task.GetTask().StartDate
		}
		if dateutil.IsGt(task.GetTask().EndDate, plan.EndDate) {
			plan.EndDate = task.GetTask().EndDate
		}
	}

	plan.Workdays = CalculateLaborableDays(plan.StartDate, plan.EndDate, plan.FeastDays)

	return nil
}

// validateTasks comprueba que las tareas tengan el formato correcto para poder realizar la planificación
func validateTasks(aTasks []ITask, resources []IResource) (map[TaskID]ITask, *Error) {
	var (
		taskIDSErrors []TaskID
		tasks         []ITask = aTasks
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
		if task.GetTask().Duration < 1 {
			taskIDSErrors = append(taskIDSErrors, task.GetTask().ID)
		}
	}

	if len(taskIDSErrors) > 0 {
		return nil, newError("las siguientes tareas tienen la duración inferior a un día", taskIDSErrors)
	}

	// No puede haber tareas con una orden menor que 1
	for _, task := range tasks {
		if task.GetTask().Order < 1 {
			taskIDSErrors = append(taskIDSErrors, task.GetTask().ID)
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
		typeOfTasks[task.GetTask().ResourceType] = true
	}

	// Crea el índice de tipos de resource
	for _, resource := range resources {
		typeOfResources[resource.GetResource().Type] = true
	}

	// Tiene que haber tareas para los tipos de recurso que llegan
	for _, resource := range resources {
		if _, exist := typeOfTasks[resource.GetResource().Type]; !exist {
			return nil, newTextError("no existen tareas para el recurso %s de tipo %s", resource.GetResource().ID, resource.GetResource().Type)
		}
	}

	// Tiene que haber recursos para las tareas del tipo que llevan asociado
	for _, task := range tasks {
		if _, exist := typeOfResources[task.GetTask().ResourceType]; !exist {
			taskIDSErrors = append(taskIDSErrors, task.GetTask().ID)
		}
	}
	if len(taskIDSErrors) > 0 {
		return nil, newError("no hay recursos para los tipos de estas tareas", taskIDSErrors)
	}

	// Crea el índice de tareas para poder saber a qué tarea corresponde un tareaID
	var tasksIndex = make(map[TaskID]ITask, len(tasks))
	for _, task := range tasks {
		tasksIndex[task.GetTask().ID] = task
	}

	var tasksNotExists = map[TaskID]bool{}

	// No puede haber tareas en blocksBy ni blocksTo que no estén en la lista de tareas
	for _, task := range tasks {

		// Revisa las tareas a las que bloquea
		for _, dep := range task.GetTask().BlocksTo {
			if _, exist := tasksIndex[dep.GetTaskDependency().TaskID]; !exist {
				tasksNotExists[dep.GetTaskDependency().TaskID] = true
			}
		}

		// Revisa las tareas que le bloquean
		for _, dep := range task.GetTask().BlocksBy {
			if _, exist := tasksIndex[dep.GetTaskDependency().TaskID]; !exist {
				tasksNotExists[dep.GetTaskDependency().TaskID] = true
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
		for _, dep := range task.GetTask().BlocksTo {
			taskBlocked := tasksIndex[dep.GetTaskDependency().TaskID]
			if taskBlocked != nil {
				if task.GetTask().Order >= taskBlocked.GetTask().Order {
					taskIDSErrors = append(taskIDSErrors, dep.GetTaskDependency().TaskID)
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
func checkForCircularDependencies(task ITask, tasksIndex map[TaskID]ITask, linkedblocksTo []TaskID) *Error {

	for _, tID := range linkedblocksTo {
		if tID == task.GetTask().ID {
			return newTextError("%s contiene una referencia circular en su cadena de dependencias", tID)
		}
	}

	// Recorre las tareas a las que está bloqueando llamando recursivamente a la propia función
	for _, dep := range task.GetTask().BlocksTo {
		// Añade task.ID a la secuencia de bloqueos para detectar una posible referencia circular a ella misma en
		// dicho encadenamiento
		linkedblocksTo = append(linkedblocksTo, task.GetTask().ID)
		err := checkForCircularDependencies(tasksIndex[dep.GetTaskDependency().TaskID], tasksIndex, linkedblocksTo)
		if err != nil {
			return err
		}
	}
	return nil
}

// assignTask Asigna la tarea al recurso que en una simulación de planificación pueda acabarla antes
func assignTask(task ITask, resources []IResource, feastDays []IHolidays, tasksIndex map[TaskID]ITask) *Error {

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

	task.GetTask().StartDate = startDate

	// Le asigna los datos de la planificación
	scheduledTaskInfo = bestScheduledTask(task, resources, feastDays)

	task.GetTask().Resource = scheduledTaskInfo.Resource
	task.GetTask().StartDate = scheduledTaskInfo.StartDate
	task.GetTask().EndDate = scheduledTaskInfo.EndDate

	return nil
}

// getRealStartDate Retorna la fecha real de comienzo de una tarea teniendo en cuenta que si tiene bloqueos
// calcula la fecha mayor de las que bloquean a la tarea, ya que no se podrá empezar antes.
func getRealStartDate(task ITask, tasksIndex map[TaskID]ITask) (time.Time, *Error) {

	var blocksBy []ITaskDependency = task.GetTask().BlocksBy

	// Si no tiene tareas que la bloqueen retorna la misma fecha de comienzo
	if len(blocksBy) == 0 {
		return task.GetTask().StartDate, nil
	}

	var tasksBlocksBy []ITask

	// Las tareas que la bloquean en una lista de tareas ordenadas primero por las que bloquean a otros deberían
	// estar ya planificadas, si no lo están damos error
	for _, dep := range blocksBy {

		taskBlocksBy := tasksIndex[dep.GetTaskDependency().TaskID]

		if taskBlocksBy.GetTask().Resource == nil {
			// Esto debería ser muy improbable que se dé si el código funciona como debe...
			return time.Time{}, newTextError("La tarea %s no está planificada y bloquea a la tarea %s que está en planificación",
				dep.GetTaskDependency().TaskID, task.GetTask().ID)
		}

		tasksBlocksBy = append(tasksBlocksBy, taskBlocksBy)
	}

	// Ordena blocksBy por fecha de fin de tarea descendiente para encontrar la fecha mayor
	sort.Slice(tasksBlocksBy, func(i, j int) bool {
		return dateutil.IsGt(tasksBlocksBy[i].GetTask().EndDate, tasksBlocksBy[j].GetTask().EndDate)
	})

	// Retorna la fecha mayor + 1 día ya que debe comenzar al día siguiente de la fecha de fin de la última tarea que la bloquean
	return tasksBlocksBy[0].GetTask().EndDate.AddDate(0, 0, 1), nil

}

// bestScheduledTask Calcula la planificación de la tarea para cada recurso y retorna la que termine antes.
func bestScheduledTask(task ITask, resources []IResource, feastDays []IHolidays) *scheduledTaskInfo {

	var (
		availableResources []IResource
		scheduledTasks     []*scheduledTaskInfo
		bestScheduled      *scheduledTaskInfo
	)

	// Filtra los recursos del tipo de tarea
	for _, resource := range resources {
		if resource.GetResource().Type == task.GetTask().ResourceType {
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
	bestScheduled.Resource.GetResource().nextAvailableDate = bestScheduled.EndDate.AddDate(0, 0, 1)

	return bestScheduled
}

// scheduledTask planifica una tarea para un recurso
func scheduledTask(task ITask, resource IResource, feastDays []IHolidays) *scheduledTaskInfo {

	var (
		holidaysAndFeastDays []IHolidays
		realStartDate        time.Time
		duration             uint
		endDate              time.Time
	)

	// concatena días de fiesta y vacaciones del recurso
	holidaysAndFeastDays = append(holidaysAndFeastDays, resource.GetResource().Holidays...)
	holidaysAndFeastDays = append(holidaysAndFeastDays, feastDays...)

	// Si la fecha en la que debe comenzar la tarea es superior a la fecha en la que el recurso estaría disponible
	// Ponemos esa fecha como fecha en la que el recurso estaría disponible para el cálculo y si no se pone la fecha en
	// la que estaría disponible el recurso
	if dateutil.IsGt(task.GetTask().StartDate, resource.GetResource().nextAvailableDate) {
		realStartDate = task.GetTask().StartDate
	} else {
		realStartDate = resource.GetResource().nextAvailableDate
	}

	duration = task.GetTask().Duration
	endDate = realStartDate

	for {

		if IsLaborableDay(endDate, holidaysAndFeastDays) {
			// Si es el primer día que empieza a contar actualiza la fecha de comienzo, puede que aunque la fecha de
			// comienzo inicial sea hoy, hoy y mañana sean fiesta por lo que comenzaría dos días después
			if duration == task.GetTask().Duration {
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
