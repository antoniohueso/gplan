package main

type Task struct {
	Name string
	Age  int
}

type Project struct {
	TotalTask int
	Tasks     []*Task
}

func main() {

	var project = &Project{
		TotalTask: 10,
		Tasks: []*Task{
			&Task{Name: "Antoño", Age: 10},
			&Task{Name: "Antoño1", Age: 11},
			&Task{Name: "Antoño2", Age: 12},
		},
	}

	var p = *project
	var tasks []*Task
	for _, t := range p.Tasks {
		var tt = *t
		tasks = append(tasks, &tt)
	}
	p.Tasks = tasks

	p.TotalTask = 2
	p.Tasks[0].Name = "Pollas"

	println(p.TotalTask, " ", project.TotalTask)
	println(p.Tasks[0].Name, " ", project.Tasks[0].Name)
}
