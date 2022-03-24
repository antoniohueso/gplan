package gplan

import (
	"math"
	"time"

	"github.com/antoniohueso/gplan/dateutil"
)

// Review actualiza el estado de avance de una planificación a fecha de reviewDate
func Review(plan ProjectPlan, reviewDate time.Time) *Error {

	var (
		tasks        = plan.GetTasks()
		feastDays    = plan.GetFeastDays()
		resourcesIdx = make(map[ResourceID]Resource)
	)

	// Si el plan no está planificado aun retorna error
	for _, task := range tasks {
		if task.GetResourceID() == nil {
			return newTextError("Hay tareas sin planificar aun")
		}
	}

	// Crea el índice de recursos
	for _, r := range plan.GetResources() {
		resourcesIdx[r.GetID()] = r
	}

	// Calcula el avance estimado
	plan.SetExpectedProgress(CalculateExpectedProgress(
		tasks,
		feastDays,
		resourcesIdx,
		reviewDate))

	// Calcula el avance real
	plan.SetRealProgress(CalculateRealProgress(tasks))

	// Si el plan está al 100% calcula la fecha de finalización real mayor para ponerla como fecha final real del proyecto
	// Y calcular el retraso o adelanto
	var lastResolvedDate time.Time
	if plan.GetRealProgress() == 100 || plan.IsArchived() {
		for _, task := range tasks {
			if dateutil.IsGt(task.GetRealEndDate(), lastResolvedDate) {
				lastResolvedDate = task.GetRealEndDate()
			}
		}
	}

	var realProgressDays = calculateProgressDays(
		reviewDate,
		plan.GetExpectedProgress(),
		plan.GetRealProgress(),
		plan.GetStartDate(),
		plan.GetEndDate(),
		lastResolvedDate,
		feastDays)

	// Calcula el avance o retraso en la planificación en días
	plan.SetRealProgressDays(math.Round(realProgressDays*100) / 100)

	// Calcula el avance o el retraso en función de Avance o retraso en días y teniendo en cuenta los días de fiesta
	// Redondea a la alta los días de retraso y a la baja los de adelanto de manera que si es -1.2 será -1 y si es 1.2 será 2.
	// Es decir 1.3 días de adelanto para gplan será un día de adelanto y 1.3 días de retraso serán 2 días
	plan.SetEstimatedEndDate(
		CalculateLaborableDate(plan.GetEndDate(), int(math.Ceil(plan.GetRealProgressDays())), feastDays))

	// Calcula el total de tareas completadas
	plan.SetCompleteTasks(CalculateTotalTasksCompleted(tasks))

	plan.SetWorkdaysToEndDate(CalculateLaborableDays(time.Now(), plan.GetEndDate(), feastDays))

	return nil
}

// CalculateExpectedProgress Calcula el % de avance en el que deberíamos estar en el día actual en función de la planificación.
// Sigue la programación de las tareas de manera que si currDate > que la fecha de fin de la tarea esta se considera como que debería
// estar completa al 100% y si currDate es < startDate entonces se considera que debería estar al 0%. Si currDate está entre
// startDate y endDate de una tarea calcula el % que debería llevar hasta el día actual.
func CalculateExpectedProgress(tasks []Task, feastDays []Holidays, resourceIdx map[ResourceID]Resource, currDate time.Time) uint {

	var expectedProgressDays uint
	var totalDuration uint

	for _, task := range tasks {

		totalDuration += task.GetDuration()

		// Está en el intervalo de fecha de comienzo y final de la tarea, calcula el avance estimado
		if dateutil.IsBetween(currDate, task.GetStartDate(), task.GetEndDate()) {
			var holidaysAndFeastDays []Holidays
			// concatena días de fiesta y vacaciones del recurso
			var resource = resourceIdx[*task.GetResourceID()]
			holidaysAndFeastDays = append(holidaysAndFeastDays, resource.GetHolidays()...)
			holidaysAndFeastDays = append(holidaysAndFeastDays, feastDays...)

			currDays := CalculateLaborableDays(task.GetStartDate(), currDate.AddDate(0, 0, -1), holidaysAndFeastDays)
			expectedProgressDays += currDays
			task.SetExpectedProgress((currDays * 100) / task.GetDuration())
		} else if dateutil.IsGt(currDate, task.GetEndDate()) {
			// Se ha pasado de la fecha fin, debería estar al 100%
			task.SetExpectedProgress(100)
			expectedProgressDays += task.GetDuration()
		} else {
			// Es menor que la fecha de comienzo, debería estar al 0%
			task.SetExpectedProgress(0)
		}
	}

	// Se suman las duraciones que deberían estar completas o a medio completar y se calcula el % con respecto al
	// total de la duración
	// TODO: Esto lo calcula siempre a la baja,probar algún dia a redondear normalmente con math.Round(valor)
	return (expectedProgressDays * 100) / totalDuration
}

// CalculateRealProgress Calcula el % de avance real con la siguiente fórmula:
// [Suma el resultado de multiplicar % real de cada tarea * su duración y el resultado lo divide por el total de la duración]
func CalculateRealProgress(tasks []Task) uint {

	var totalCompleteXDuration uint
	var totalDuration uint

	for _, task := range tasks {
		totalCompleteXDuration += task.GetRealProgress() * task.GetDuration()
		totalDuration += task.GetDuration()
	}

	return totalCompleteXDuration / totalDuration
}

// calculateProgressDays Calcula la los días de retraso o adelanto que llevamos, si es positivo el valor será retraso
// si es negativo será adelanto
func calculateProgressDays(
	reviewDate time.Time,
	shouldCompleted uint,
	realCompleted uint,
	startDate time.Time,
	endDate time.Time,
	lastResolvedDate time.Time,
	feastDays []Holidays) float64 {

	var result float64

	// Si la planificación se ha completado
	if realCompleted == 100 {

		// Si la fecha en de la última resolución es > que la fecha de fin de planificación
		if dateutil.IsGt(lastResolvedDate, endDate) {
			// Calcula los días que van desde la fecha final + 1 y la fecha de resolución y serán días de retraso
			result = float64(CalculateLaborableDays(endDate.AddDate(0, 0, 1), lastResolvedDate, feastDays))
		} else if dateutil.IsEqual(lastResolvedDate, endDate) {
			// Si la última fecha de resolución coincide con la fecha de fin de proyecto no hay retraso ni adelanto
			result = 0.0
		} else {
			// Calcula los días que van desde la fecha de resolución de la última tarea gasta la fecha final
			// y serán días de adelanto (por eso se multiplica por -1)
			result = float64(CalculateLaborableDays(lastResolvedDate, endDate, feastDays)) * -1
		}

	} else {
		// Si no está completado ni debería estarlo
		// Se calculan los días de adelanto o de retraso normalmente
		days := float64(CalculateLaborableDays(startDate, endDate, feastDays))

		// Calcula los días que deberían estar completados según el % avance estimado
		daysShouldCompleted := (float64(shouldCompleted) * days) / 100.0
		// Calcula los días que deberían estar completados según el % avance real
		daysRealCompleted := (float64(realCompleted) * days) / 100.0
		// resta los días estimados que deberían estar completados - los días completados reales
		// Si el importe es positivo, son días de atraso y si el importe es 0 no hay retraso ni adelanto, vamos al día.
		result = daysShouldCompleted - daysRealCompleted

		// Forma anterior de calcularlo me gusta más la anterior, es más intuitiva
		//diff := float64(shouldCompleted - realCompleted)
		//result = (diff * days) / 100.0

		if shouldCompleted == 100 {
			// Si no se ha completado, pero debería estar ya al 100%
			// Calcula como días de retraso los que van desde la fecha de fin + 1 y la fecha de hoy - 1, ya que no se debe
			// tener en cuenta la fecha actual en el seguimiento
			result2 := float64(CalculateLaborableDays(endDate.AddDate(0, 0, 1), reviewDate.AddDate(0, 0, -1), feastDays))
			// se los suma a los anteriores
			result += result2
		}
	}

	return result
}

// CalculateTotalTasksCompleted Calcula el número de tareas completas
func CalculateTotalTasksCompleted(tasks []Task) uint {

	var total uint

	for _, task := range tasks {
		if task.GetRealProgress() == 100 {
			total++
		}
	}

	return total
}
