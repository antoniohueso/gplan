package main

import (
	"fmt"
	"sort"
	"time"

	"com.github.antoniohueso/gplan"
	. "com.github.antoniohueso/gplan/sample"
)

type TaskExt struct {
	*Task
	Status string
}

type ProjectPlanExt struct {
	*ProjectPlan
	Tasks []*TaskExt
}

// GetTasks Devuelve un array de gplan.Task
func (p ProjectPlanExt) GetTasks() []gplan.ITask {
	newArr := make([]gplan.ITask, len(p.Tasks))
	for i := range p.Tasks {
		newArr[i] = p.Tasks[i]
	}
	return newArr
}

// SortTasksByOrder Implementa SortByOrder
func (p ProjectPlanExt) SortTasksByOrder() {
	// Ordena las tareas por número de orden para poder planificarlas
	sort.Slice(p.Tasks, func(i, j int) bool {
		return p.Tasks[i].Order < p.Tasks[j].Order
	})
}

func main() {

	h := NewHolidays(time.Now(), time.Now())
	h1 := NewHolidays(time.Now().AddDate(0, 1, 0), time.Now().AddDate(0, 1, 0))

	r := NewResource("ahg", "Antonio Hueso", "backend", time.Now(), []*Holidays{h, h1})

	var tasks []*TaskExt

	for i := 0; i < 10; i++ {
		task := &TaskExt{
			Task:   NewTask(gplan.TaskID(fmt.Sprintf("Tarea-%d", i)), "summary sample", "backend", i, i*2),
			Status: "Backlog",
		}
		tasks = append(tasks, task)
	}

	plan := &ProjectPlanExt{
		ProjectPlan: &ProjectPlan{},
	}
	plan.Tasks = tasks
	plan.Resources = []*Resource{r}
	plan.FeastDays = []*Holidays{h}

	prueba(plan)
}

func prueba(p gplan.IProjectPlan) {

	p.SortTasksByOrder()

	for _, t := range p.GetTasks() {
		fmt.Printf("%s\n", t.GetID())
	}

}
