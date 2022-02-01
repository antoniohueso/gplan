package main

import (
	"fmt"
	"time"

	"com.github.antoniohueso/gplan"
	"com.github.antoniohueso/gplan/defmodel"
)

func main() {

	h := defmodel.NewHolidays(time.Now(), time.Now())
	h1 := defmodel.NewHolidays(time.Now().AddDate(0, 1, 0), time.Now().AddDate(0, 1, 0))

	r := defmodel.NewResource("ahg", "Antonio Hueso", "backend", time.Now(), defmodel.ArrayOfHolidays{h, h1})

	var tasks []*defmodel.Task

	for i := 0; i < 10; i++ {
		tasks = append(tasks, defmodel.NewTask(gplan.TaskID(fmt.Sprintf("Tarea-%d", i)), "summary sample", "backend", i, i*2))
	}

	plan := defmodel.NewProjectPlan("myPlan", tasks, defmodel.ArrayOfResources{r}, defmodel.ArrayOfHolidays{h})

	prueba(plan)
}

func prueba(p gplan.ProjectPlan) {

	p.GetTasks().SortByOrder()

	for _, t := range p.GetTasks().Iterable() {
		fmt.Printf("%v\n", t)
	}

}
