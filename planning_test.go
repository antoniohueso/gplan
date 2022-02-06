package gplan_test

import (
	"fmt"
	"log"
	"time"

	"com.github.antoniohueso/gplan"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Planning", func() {

	var (
		resources []*gplan.Resource
		feastDays []*gplan.Holidays
	)

	BeforeEach(func() {

		resources = []*gplan.Resource{
			gplan.NewResource("ahg", "Antonio Hueso", "backend", parseDate("2021-06-10"), []*gplan.Holidays{
				gplan.NewHolidays(parseDate("2021-06-17"), parseDate("2021-06-18")),
			}),
			gplan.NewResource("cslopez", "Carlos Sobrino", "backend", parseDate("2021-06-07"), nil),
			gplan.NewResource("David.Attrache", "David Attrache", "maquetacion", parseDate("2021-06-07"), nil),
			gplan.NewResource("Noemi", "Noe Medina", "maquetacion", parseDate("2021-06-07"), nil),
		}

		feastDays = []*gplan.Holidays{
			gplan.NewHolidays(parseDate("2021-06-10"), parseDate("2021-06-11")),
		}

	})

	Describe("NewPlanning - Validaciones", func() {

		println(resources)
		println(feastDays)

		Context("Si recibe una lista de tareas vacía", func() {
			It("Debe devolver un error", func() {

				plan := gplan.NewProjectPlan("test", nil, nil, nil)

				err := gplan.Planning(time.Now(), plan)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("la lista de tareas a planificar está vacía")))
			})
		})

		Context("Si recibe una lista de recursos vacía", func() {
			It("Debe devolver un error", func() {

				plan := gplan.NewProjectPlan("test", []*gplan.Task{{}}, nil, nil)
				err := gplan.Planning(time.Now(), plan)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("la lista de recursos a asignar está vacía")))
			})
		})

		Context("Si hay tareas con una duración < 1", func() {

			It("Debe devolver un error", func() {
				plan := gplan.NewProjectPlan("test", []*gplan.Task{
					gplan.NewTask("Task-2", "summary Task-2", "maquetación", 100, 1),
					gplan.NewTask("Task-1", "summary Task-1", "maquetación", 100, 0)},
					[]*gplan.Resource{{}},
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

				plan := gplan.NewProjectPlan("test", []*gplan.Task{
					gplan.NewTask("Task-2", "summary Task-2", "maquetación", 100, 1),
					gplan.NewTask("Task-1", "summary Task-1", "backend", 100, 2)},
					[]*gplan.Resource{
						gplan.NewResource("ahg", "Antonio Hueso", "backend", time.Now(), nil),
					},
					nil)

				err := gplan.Planning(time.Now(), plan)

				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("no hay recursos para los tipos de estas tareas")))
				Expect(len(err.Tasks)).Should(Equal(1))
				Expect(err.Tasks[0]).Should(BeEquivalentTo("Task-2"))
			})
		})

		Context("Si no hay tareas del tipo necesario para los tipos de recursos que lleregan en la lista de recursos", func() {

			It("Debe devolver un error", func() {

				plan := gplan.NewProjectPlan("test",
					[]*gplan.Task{
						gplan.NewTask("Task-2", "summary Task-2", "maquetación", 100, 1),
					},
					[]*gplan.Resource{
						gplan.NewResource("ahg", "Antonio Hueso", "backend", time.Now(), nil),
					},
					nil)

				err := gplan.Planning(time.Now(), plan)

				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("no existen tareas para el recurso ahg de tipo backend")))
			})
		})

		Context("Si hay referencias circulares simples", func() {

			It("Debe devolver un error", func() {
				task1 := gplan.NewTaskWithBlocks("Task-1", "summary Task-1", "backend", 10, 1, []*gplan.TaskDependency{gplan.NewTaskDependency("Task-2")}, nil)
				task2 := gplan.NewTaskWithBlocks("Task-2", "summary Task-2", "backend", 20, 2, []*gplan.TaskDependency{gplan.NewTaskDependency("Task-1")}, nil)

				plan := gplan.NewProjectPlan("test",
					[]*gplan.Task{task1, task2},
					[]*gplan.Resource{gplan.NewResource("ahg", "Antonio Hueso", "backend", time.Now(), nil)},
					nil)

				err := gplan.Planning(time.Now(), plan)

				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("Task-1 contiene una referencia circular en su cadena de dependencias")))
			})
		})

		Context("Si hay referencias circulares complejas", func() {

			It("Debe devolver un error", func() {
				task1 := gplan.NewTaskWithBlocks("Task-1", "summary Task-1", "backend", 10, 1, []*gplan.TaskDependency{gplan.NewTaskDependency("Task-2")}, nil)
				task2 := gplan.NewTaskWithBlocks("Task-2", "summary Task-2", "backend", 20, 2, []*gplan.TaskDependency{gplan.NewTaskDependency("Task-3")}, nil)
				task3 := gplan.NewTaskWithBlocks("Task-3", "summary Task-3", "backend", 30, 2, []*gplan.TaskDependency{gplan.NewTaskDependency("Task-1")}, nil)

				plan := gplan.NewProjectPlan("test",
					[]*gplan.Task{task1, task2, task3},
					[]*gplan.Resource{gplan.NewResource("ahg", "Antonio Hueso", "backend", time.Now(), nil)},
					nil)

				err := gplan.Planning(time.Now(), plan)

				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("Task-1 contiene una referencia circular en su cadena de dependencias")))
			})
		})

		Context("Si hay tareas que bloquean a tareas que no existen en la lista de tareas", func() {

			It("Debe devolver un error", func() {
				task1 := gplan.NewTaskWithBlocks("Task-1", "summary Task-1", "backend", 10, 1, []*gplan.TaskDependency{gplan.NewTaskDependency("Task-No-Existe")}, nil)
				task2 := gplan.NewTaskWithBlocks("Task-2", "summary Task-2", "backend", 20, 2, nil, nil)
				task3 := gplan.NewTaskWithBlocks("Task-3", "summary Task-3", "backend", 30, 2, nil, nil)

				plan := gplan.NewProjectPlan("test",
					[]*gplan.Task{task1, task2, task3},
					[]*gplan.Resource{gplan.NewResource("ahg", "Antonio Hueso", "backend", time.Now(), nil)},
					nil)
				err := gplan.Planning(time.Now(), plan)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("hay tareas bloquedas o que bloquean a otras que no existen en la lista de tareas")))
				Expect(len(err.Tasks)).Should(Equal(1))
				Expect(err.Tasks[0]).Should(BeEquivalentTo("Task-No-Existe"))
			})
		})

		Context("Si hay tareas que bloquean a otras cuyo orden es inferior", func() {

			It("Debe devolver un error", func() {
				task1 := gplan.NewTaskWithBlocks("Task-1", "summary Task-1", "backend", 20, 1, []*gplan.TaskDependency{gplan.NewTaskDependency("Task-2")}, nil)
				task2 := gplan.NewTaskWithBlocks("Task-2", "summary Task-2", "backend", 10, 2, nil, nil)

				plan := gplan.NewProjectPlan("test",
					[]*gplan.Task{task1, task2},
					[]*gplan.Resource{gplan.NewResource("ahg", "Antonio Hueso", "backend", time.Now(), nil)},
					nil)
				err := gplan.Planning(time.Now(), plan)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Message).Should(Equal(fmt.Errorf("tareas bloqueadas por tareas con un orden superior")))
				Expect(len(err.Tasks)).Should(Equal(1))
				Expect(err.Tasks[0]).Should(BeEquivalentTo("Task-2"))
			})
		})

	})

	Describe("NewPlanning - Creación", func() {
		Context("Creamos un plan sin dependencias con dos días de fiesta", func() {

			It("El plan debe ser igual al siguiente plan", func() {

				tasks := []*gplan.Task{
					gplan.NewTask("Tarea1", "Summary", "backend", 10, 10),
					gplan.NewTask("Tarea2", "Summary", "maquetacion", 50, 2),
					gplan.NewTask("Tarea3", "Summary", "backend", 20, 1),
					gplan.NewTask("Tarea4", "Summary", "maquetacion", 60, 20),
					gplan.NewTask("Tarea5", "Summary", "backend", 30, 7),
					gplan.NewTask("Tarea6", "Summary", "backend", 40, 9),
				}

				plan := gplan.NewProjectPlan("test-plan", tasks, resources, feastDays)
				err := gplan.Planning(parseDate("2021-06-07"), plan)

				Expect(err).Should(BeNil())
				comparePlan(plan.Tasks, []string{
					"2021-06-07 2021-06-22 cslopez",
					"2021-06-14 2021-06-14 ahg",
					"2021-06-15 2021-06-25 ahg",
					"2021-06-23 2021-07-05 cslopez",
					"2021-06-07 2021-06-08 David.Attrache",
					"2021-06-07 2021-07-06 Noemi",
				})
				Expect(plan.StartDate.Format("2006-01-02")).Should(Equal("2021-06-07"))
				Expect(plan.EndDate.Format("2006-01-02")).Should(Equal("2021-07-06"))
			})
		})

		Context("Creamos un plan con dependencias con dos días de fiesta", func() {

			It("El plan debe ser igual al siguiente plan", func() {

				tasks := []*gplan.Task{
					gplan.NewTask("Tarea1", "Summary", "backend", 30, 10),
					gplan.NewTask("Tarea2", "Summary", "maquetacion", 50, 2),
					gplan.NewTask("Tarea3", "Summary", "backend", 40, 1),
					gplan.NewTask("Tarea4", "Summary", "maquetacion", 20, 20),
					gplan.NewTask("Tarea5", "Summary", "backend", 60, 7),
					gplan.NewTask("Tarea6", "Summary", "backend", 10, 9),
				}

				gplan.BlocksTo(tasks[3], tasks[0])
				gplan.BlocksTo(tasks[5], tasks[2])
				gplan.BlocksTo(tasks[5], tasks[1])

				plan := gplan.NewProjectPlan("test-plan", tasks, resources, feastDays)
				err := gplan.Planning(parseDate("2021-06-07"), plan)

				Expect(err).Should(BeNil())
				comparePlan(plan.Tasks, []string{
					"2021-06-07 2021-06-21 cslopez",
					"2021-06-07 2021-07-06 David.Attrache",
					"2021-07-07 2021-07-20 ahg",
					"2021-06-22 2021-06-22 cslopez",
					"2021-06-22 2021-06-23 Noemi",
					"2021-06-23 2021-07-01 cslopez",
				})
				Expect(plan.StartDate.Format("2006-01-02")).Should(Equal("2021-06-07"))
				Expect(plan.EndDate.Format("2006-01-02")).Should(Equal("2021-07-20"))
			})
		})

	})

})

func printPlan(plan *gplan.ProjectPlan) {
	fmt.Println()
	for _, task := range plan.Tasks {
		fmt.Printf("%s|%02d|%s|%s|%s\n", task.ID, task.Duration, task.StartDate.Format("2006-01-02"), task.EndDate.Format("2006-01-02"),
			task.Resource.GetID())
	}
}

func comparePlan(tasks []*gplan.Task, compare []string) {
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
