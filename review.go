package gplan

import (
	"math"
	"time"

	"github.com/antoniohueso/gplan/dateutil"
)

// TODO: Introducir los cambios que calculalos días si está al 100% o está archivado

// Review actualiza el estado de avance de una planificación a fecha de reviewDate
func Review(plan ProjectPlan, reviewDate time.Time) *Error {

	var (
		tasks     = plan.GetTasks()
		feastDays = plan.GetFeastDays()
	)

	// Si el plan no está planificado aun retorna error
	for _, task := range tasks {
		if task.GetResourceID() == nil {
			return newTextError("Hay tareas sin planificar aun")
		}
	}

	// Calcula el avance estimado
	CalculateExpectedProgress(plan, reviewDate)

	// Calcula el avance real
	CalculateRealProgress(plan)

	// Calcula el progreso (positivo o negativo) en días
	CalculateProgressDays(plan, reviewDate)

	// Calcula el total de tareas completadas
	CalculateTotalTasksCompleted(plan)

	plan.SetWorkdaysToEndDate(CalculateLaborableDays(time.Now(), plan.GetEndDate(), feastDays))

	return nil
}

// CalculateExpectedProgress Calcula el % de avance en el que deberíamos estar en el día actual en función de la planificación.
// Sigue la programación de las tareas de manera que si currDate > que la fecha de fin de la tarea esta se considera como que debería
// estar completa al 100% y si currDate es < startDate entonces se considera que debería estar al 0%. Si currDate está entre
// startDate y endDate de una tarea calcula el % que debería llevar hasta el día actual.
func CalculateExpectedProgress(plan ProjectPlan, reviewDate time.Time) {
	var (
		expectedProgressDuration uint
		totalDuration            uint
		feastDays                = plan.GetFeastDays()
		resourcesIdx             = make(map[ResourceID]Resource)
	)

	// Crea el índice de recursos
	for _, r := range plan.GetResources() {
		resourcesIdx[r.GetID()] = r
	}

	for _, task := range plan.GetTasks() {

		if dateutil.IsLte(reviewDate, plan.GetStartDate()) || dateutil.IsLte(reviewDate, task.GetStartDate()) {
			// Si la fecha de revisión es <= que la fecha de comienzo del plan o que la fecha de comienzo de la tarea
			// debería estar al 0%
			task.SetExpectedProgress(0)
			task.SetExpectedCompleteDuration(0)
		} else if dateutil.IsGt(reviewDate, task.GetEndDate()) {
			// Si la fecha de revisión es > que la fecha de fin
			// Se ha pasado de la fecha fin, debería estar al 100%
			task.SetExpectedProgress(100)
			task.SetExpectedCompleteDuration(task.GetDuration())
		} else {
			// Si la fecha de revisión está entre la fecha de inicio y la de fin de la tarea, calcula el progreso esperado en base a la duración
			// que debería llevar
			var holidaysAndFeastDays []Holidays
			// concatena días de fiesta y vacaciones del recurso
			var resource = resourcesIdx[*task.GetResourceID()]
			holidaysAndFeastDays = append(holidaysAndFeastDays, resource.GetHolidays()...)
			holidaysAndFeastDays = append(holidaysAndFeastDays, feastDays...)

			currDays := CalculateLaborableDays(task.GetStartDate(), reviewDate.AddDate(0, 0, -1), holidaysAndFeastDays)
			task.SetExpectedProgress((currDays * 100) / task.GetDuration())
			task.SetExpectedCompleteDuration(currDays)
		}

		expectedProgressDuration += task.GetExpectedCompleteDuration()
		totalDuration += task.GetDuration()
	}

	// Se suman las duraciones que deberían estar completas o a medio completar y se calcula el % con respecto al
	// total de la duración
	plan.SetExpectedProgress((expectedProgressDuration * 100) / totalDuration)
}

// CalculateRealProgress Calcula el % de avance real
func CalculateRealProgress(plan ProjectPlan) {

	var (
		totalCompleteXDuration uint
		totalDuration          uint
	)

	for _, task := range plan.GetTasks() {
		// Calcula la duración real completada de la tarea en función del % realcompletado
		task.SetRealCompleteDuration((task.GetDuration() * task.GetRealProgress()) / 100)

		totalCompleteXDuration += task.GetRealCompleteDuration()
		totalDuration += task.GetDuration()
	}

	plan.SetRealProgress(totalCompleteXDuration * 100 / totalDuration)
}

// CalculateProgressDays Calcula la los días de retraso o adelanto que llevamos, si es positivo el valor será retraso
// si es negativo será adelanto
func CalculateProgressDays(plan ProjectPlan, reviewDate time.Time) {

	var (
		feastDays                = plan.GetFeastDays()
		expectedCompleteDuration float64
		realCompleteDuration     float64
		realProgressDays         float64
		totalDuration            float64
		workDays                 = float64(plan.GetWorkdays())
	)

	// Si la planificación se ha completado calcula los días desde la última última fecha completada
	if plan.GetRealProgress() == 100 || plan.IsArchived() {

		// Busca la última fecha completada
		var realEndDate time.Time
		for _, task := range plan.GetTasks() {
			if dateutil.IsGt(task.GetRealEndDate(), realEndDate) {
				realEndDate = task.GetRealEndDate()
			}
		}

		// Si la fecha en de la última resolución es > que la fecha de fin de planificación
		if dateutil.IsGt(realEndDate, plan.GetEndDate()) {
			// Calcula los días que van desde la fecha final + 1 y la fecha de resolución y serán días de retraso
			realProgressDays = float64(CalculateLaborableDays(plan.GetEndDate().AddDate(0, 0, 1), realEndDate, feastDays))
		} else if dateutil.IsEqual(realEndDate, plan.GetEndDate()) {
			// Si la última fecha de resolución coincide con la fecha de fin de proyecto no hay retraso ni adelanto
			realProgressDays = 0.0
		} else {
			// Calcula los días que van desde la fecha de resolución de la última tarea gasta la fecha final
			// y serán días de adelanto (por eso se multiplica por -1)
			realProgressDays = float64(CalculateLaborableDays(realEndDate, plan.GetEndDate(), feastDays)) * -1
		}

		// Redondea a 1 decimal
		plan.SetRealProgressDays(math.Round(realProgressDays*10.0) / 10.0)
		// La fecha estimada de fin es la fecha real de fin del proyecto
		plan.SetEstimatedEndDate(realEndDate)
	} else {
		// No está completado ni archivado, suma las duraciones esperadas completas y las duraciones reales completas
		for _, task := range plan.GetTasks() {
			expectedCompleteDuration += float64(task.GetExpectedCompleteDuration())
			realCompleteDuration += float64(task.GetRealCompleteDuration())
			totalDuration += float64(task.GetDuration())
		}

		// Fórmula (duración esperada - duración real) * las jornadas laborales del plan / duración total del plan
		realProgressDays = ((expectedCompleteDuration - realCompleteDuration) * workDays) / totalDuration

		// Si la fecha de revisión -1  es > que la fecha de fin del plan, calcula los días que hay desde la fecha de finalización hasta la fecha de revisión -1
		// y se los suma a los días
		if dateutil.IsGt(reviewDate.AddDate(0, 0, -1), plan.GetEndDate()) {
			realProgressDays += float64(CalculateLaborableDays(plan.GetEndDate().AddDate(0, 0, 1), reviewDate.AddDate(0, 0, -1), feastDays))
		}

		// Redondea a 1 decimal
		plan.SetRealProgressDays(math.Round(realProgressDays*10.0) / 10.0)

		// Calcula el avance o el retraso en función de Avance o retraso en días y teniendo en cuenta los días de fiesta
		// Redondea a la alta los días de retraso y a la baja los de adelanto de manera que si es -1.2 será -1 y si es 1.2 será 2.
		// Es decir 1.3 días de adelanto para gplan será un día de adelanto y 1.3 días de retraso serán 2 días
		plan.SetEstimatedEndDate(
			CalculateLaborableDate(plan.GetEndDate(), int(math.Ceil(plan.GetRealProgressDays())), feastDays))
	}

}

// CalculateTotalTasksCompleted Calcula el número de tareas completas
func CalculateTotalTasksCompleted(plan ProjectPlan) {

	var total uint

	for _, task := range plan.GetTasks() {
		if task.GetRealProgress() == 100 {
			total++
		}
	}

	plan.SetCompleteTasks(total)
}
