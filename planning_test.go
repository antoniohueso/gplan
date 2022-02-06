package gplan_test

import (
	"fmt"
	"log"
	"time"

	"com.github.antoniohueso/gplan"
	"com.github.antoniohueso/gplan/defmodel"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Planning", func() {

	var (
		resources []*defmodel.Resource
		feastDays []*defmodel.Holidays
	)

	BeforeEach(func() {

		resources = []*defmodel.Resource{
			defmodel.NewResource("ahg", "Antonio Hueso", "backend", parseDate("2021-06-10"), []*defmodel.Holidays{
				defmodel.NewHolidays(parseDate("2021-06-17"), parseDate("2021-06-18")),
			}),
			defmodel.NewResource("cslopez", "Carlos Sobrino", "backend", parseDate("2021-06-07"), nil),
			defmodel.NewResource("David.Attrache", "David Attrache", "maquetacion", parseDate("2021-06-07"), nil),
			defmodel.NewResource("Noemi", "Noe Medina", "maquetacion", parseDate("2021-06-07"), nil),
		}

		feastDays = []*defmodel.Holidays{
			defmodel.NewHolidays(parseDate("2021-06-10"), parseDate("2021-06-11")),
		}

	})

	Describe("NewPlanning - Validaciones", func() {

		println(resources)
		println(feastDays)

		Context("Si recibe una lista de tareas vacía", func() {
			It("Debe devolver un error", func() {

				plan := defmodel.NewProjectPlan("test", defmodel.ArrayOfTasks{}, nil, nil)

				err := gplan.Planning(time.Now(), plan)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("la lista de tareas a planificar está vacía")))
			})
		})

		Context("Si recibe una lista de recursos vacía", func() {
			It("Debe devolver un error", func() {

				plan := defmodel.NewProjectPlan("test", defmodel.ArrayOfTasks{{}}, nil, nil)
				err := gplan.Planning(time.Now(), plan)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("la lista de recursos a asignar está vacía")))
			})
		})

		Context("Si hay tareas con una duración < 1", func() {

			It("Debe devolver un error", func() {
				plan := defmodel.NewProjectPlan("test", defmodel.ArrayOfTasks{
					defmodel.NewTask("Task-2", "summary Task-2", "maquetación", 100, 1),
					defmodel.NewTask("Task-1", "summary Task-1", "maquetación", 100, 0)},
					[]*defmodel.Resource{{}},
					nil)

				err := gplan.Planning(time.Now(), plan)

				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("las siguientes tareas tienen la duración inferior a un día")))
				Expect(len(err.Tasks)).Should(Equal(1))
				Expect(err.Tasks[0]).Should(BeEquivalentTo("Task-1"))
			})
		})

		Context("Si hay tareas para tipos de recursos que no están en la lista de recursos", func() {

			It("Debe devolver un error", func() {

				plan := defmodel.NewProjectPlan("test", defmodel.ArrayOfTasks{
					defmodel.NewTask("Task-2", "summary Task-2", "maquetación", 100, 1),
					defmodel.NewTask("Task-1", "summary Task-1", "backend", 100, 2)},
					[]*defmodel.Resource{
						defmodel.NewResource("ahg", "Antonio Hueso", "backend", time.Now(), nil),
					},
					nil)

				err := gplan.Planning(time.Now(), plan)

				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("no hay recursos para los tipos de estas tareas")))
				Expect(len(err.Tasks)).Should(Equal(1))
				Expect(err.Tasks[0]).Should(BeEquivalentTo("Task-2"))
			})
		})

	})

})

func printPlan(plan *defmodel.ProjectPlan) {
	fmt.Println()
	for _, task := range plan.Tasks {
		fmt.Printf("%s|%02d|%s|%s|%s\n", task.ID, task.Duration, task.StartDate.Format("2006-01-02"), task.EndDate.Format("2006-01-02"),
			task.Resource.GetID())
	}
}

func comparePlan(tasks []*defmodel.Task, compare []string) {
	var sTasks []string

	for i := range tasks {
		s := fmt.Sprintf("%s %s %s", tasks[i].StartDate.Format("2006-01-02"), tasks[i].EndDate.Format("2006-01-02"),
			tasks[i].Resource.GetID())
		sTasks = append(sTasks, s)
		//log.Println(s)
	}
	for i := range tasks {
		Expect(sTasks[i]).
			Should(Equal(compare[i]))
	}
}

func parseDate(date string) time.Time {
	f, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Fatal(err)
	}
	return f
}
