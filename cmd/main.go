package main

import "fmt"

type Task struct {
	Name string
	Age  int
}

type Project struct {
	TotalTask int
	Tasks     []Task
}

type ImmutableSlice[K any] interface {
	Add(v K) ImmutableSlice[K]
	Get(index int) K
	Set(v K, index int) ImmutableSlice[K]
	Slice() []K
}

type immutableSlice[K any] struct {
	slice []K
}

func (s immutableSlice[K]) Add(v K) ImmutableSlice[K] {
	var newSlice = copySlice(s.slice)
	s.slice = newSlice
	s.slice = append(s.slice, v)
	return s
}

func (s immutableSlice[K]) Set(v K, index int) ImmutableSlice[K] {
	var newSlice = copySlice(s.slice)
	s.slice = newSlice
	s.slice[index] = v
	return s
}

func (s immutableSlice[K]) Get(index int) K {
	return s.slice[index]
}

func (s immutableSlice[K]) Slice() []K {
	return copySlice(s.slice)
}

func copySlice[K any](slice []K) []K {
	var newSlice []K
	newSlice = append(newSlice, slice...)
	return newSlice
}

func NewImmutableSlice[K any](slice ...K) ImmutableSlice[K] {
	return immutableSlice[K]{
		slice: slice,
	}
}

func main() {

	var p = NewImmutableSlice[Task]()

	for i := 0; i < 10; i++ {
		p = p.Add(Task{Name: fmt.Sprintf("Pollas %d", i), Age: 10 + i})
	}

	var s = p.Get(0)
	s.Name = "Pene"
	p.Set(s, 0)

	for _, s := range p.Slice() {
		fmt.Println(s.Name, s.Age)
	}

}

func prueba(p Project) Project {
	p.TotalTask = 10
	p.Tasks = []Task{
		{Name: "Antoño", Age: 10},
		{Name: "Antoño1", Age: 11},
		{Name: "Antoño2", Age: 12},
	}
	return p
}

func prueba2(p Project) Project {
	p.Tasks[0].Name = "Pollas"
	return p
}
