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

	prueba(r)
}

func prueba(r gplan.Resource) {

	for _, h := range r.GetHolidays().Iterable() {
		fmt.Printf("%v\n", h.GetFrom())
	}

}
