package gplan

import (
	"math"
	"time"

	"github.com/antoniohueso/gplan/dateutil"
)

// Review actualiza el estado de avance de una planificación a fecha de reviewDate
func Review(plan ProjectPlan, reviewDate time.Time) *Error {

	var (
		tasks     []Task     = plan.GetTasks()
		feastDays []Holidays = plan.GetFeastDays()
	)

	// Si el plan no está planificado aun retorna error
	for _, task := range tasks {
		if task.GetResource() == nil {
			return newTextError("Hay tareas sin planificar aun")
		}
	}

	// Calcula el avance estimado
	plan.SetEstimatedComplete(calculateEstimatedAdvance(
		tasks,
		feastDays,
		reviewDate))

	// Calcula el avance real
	plan.SetComplete(calculateRealAdvance(tasks))

	// Si el plan está al 100% calcula la fecha de finalización real mayor para ponerla como fecha final real del proyecto
	// Y calcular el retraso o adelanto
	var lastResolvedDate time.Time
	if plan.GetComplete() == 100 || plan.IsArchived() {
		for _, task := range tasks {
			if task.GetRealEndDate().After(lastResolvedDate) {
				lastResolvedDate = task.GetRealEndDate()
			}
		}
	}

	estimatedAdvancedOrDelayedDays := calculateEstimatedAdvancedOrDelayedDays(
		reviewDate,
		plan.GetEstimatedComplete(),
		plan.GetComplete(),
		plan.GetStartDate(),
		plan.GetEndDate(),
		lastResolvedDate,
		feastDays)

	plan.SetRealAdvancedOrDelayed(math.Round(estimatedAdvancedOrDelayedDays*100) / 100)

	return nil
}

// calculateEstimatedAdvance Calcula el % de avance en el que deberíamos estar en el día actual en función de la planificación.
// Sigue la programación de las tareas de manera que si currDate > que la fecha de fin de la tarea esta se considera como que debería
// estar completa al 100% y si currDate es < startDate entonces se considera que debería estar al 0%. Si currDate está entre
// startDate y endDate de una tarea calcula el % que debería llevar hasta el día actual.
func calculateEstimatedAdvance(tasks []Task, feastDays []Holidays, currDate time.Time) int {

	var estimatedAdvanced = 0
	var totalDuration = 0

	for _, task := range tasks {

		totalDuration += task.GetDuration()

		// Está en el intervalo de fecha de comienzo y final de la tarea, calcula el avance estimado
		if dateutil.IsBetween(currDate, task.GetStartDate(), task.GetEndDate()) {
			var holidaysAndFeastDays []Holidays
			// concatena días de fiesta y vacaciones del recurso
			holidaysAndFeastDays = append(holidaysAndFeastDays, task.GetResource().GetHolidays()...)
			holidaysAndFeastDays = append(holidaysAndFeastDays, feastDays...)

			currDays := calculateLaborableDays(task.GetStartDate(), currDate.AddDate(0, 0, -1), holidaysAndFeastDays)
			estimatedAdvanced += currDays
			task.SetEstimatedComplete((currDays * 100) / task.GetDuration())
		} else if currDate.After(task.GetEndDate()) {
			// Se ha pasado de la fecha fin, debería estar al 100%
			task.SetEstimatedComplete(100)
			estimatedAdvanced += task.GetDuration()
		} else {
			// Es menor que la fecha de comienzo, debería estar al 0%
			task.SetEstimatedComplete(0)
		}
	}

	// Se suman las duraciones que deberían estar completas o a medio completar y se calcula el % con respecto al
	// total de la duración
	return (estimatedAdvanced * 100) / totalDuration
}

// calculateRealAdvance Calcula el % de avance real sumando la duración completada realmente y sacando el % con respecto de la suma
// de todas las duraciones.
func calculateRealAdvance(tasks []Task) int {

	var totalCompleteXDuration int = 0
	var totalDuration = 0

	for _, task := range tasks {
		totalCompleteXDuration += task.GetComplete() * task.GetDuration()
		totalDuration += task.GetDuration()
	}

	return totalCompleteXDuration / totalDuration
}

func calculateEstimatedAdvancedOrDelayedDays(
	reviewDate time.Time,
	shouldCompleted int,
	realCompleted int,
	startDate time.Time,
	endDate time.Time,
	lastResolvedDate time.Time,
	feastDays []Holidays) float64 {

	var result float64 = 0.0

	// Si la planificación se ha completado
	if realCompleted == 100 {

		// Si la fecha en de la última resolución es > que la fecha de fin de planificación
		if lastResolvedDate.After(endDate) {
			// Calcula los días que van desde la fecha final + 1 y la fecha de resolución y serán días de retraso
			result = float64(calculateLaborableDays(endDate.AddDate(0, 0, 1), lastResolvedDate, feastDays))
		} else if lastResolvedDate.Equal(endDate) {
			// Si la última fecha de resolución coincide con la fecha de fin de proyecto no hay retraso ni adelanto
			result = 0.0
		} else {
			// Calcula los días que van desde la fecha de resolución de la última tarea gasta la fecha final
			// y serán días de adelanto (por eso se multiplica por -1)
			result = float64(calculateLaborableDays(lastResolvedDate, endDate, feastDays)) * -1
		}

	} else {
		// Si no está completado ni debería estarlo
		// Se calculan los días de adelanto o de retraso normalmente
		days := float64(calculateLaborableDays(startDate, endDate, feastDays))
		diff := float64(shouldCompleted - realCompleted)
		result = (diff * days) / 100.0

		if shouldCompleted == 100 {
			// Si no se ha completado, pero debería estar ya al 100%
			// Calcula como días de retraso los que van desde la fecha de fin + 1 y la fecha de hoy - 1, ya que no se debe
			// tener en cuenta la fecha actual en el seguimiento
			result2 := float64(calculateLaborableDays(endDate.AddDate(0, 0, 1), reviewDate.AddDate(0, 0, -1), feastDays))
			// se los suma a los anteriores
			result += result2
		}
	}

	return result
}
