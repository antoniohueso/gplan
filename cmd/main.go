package main

import (
	"fmt"
	"time"

	"com.github.antoniohueso/gplan"
)

type holidaysExt struct {
	Name string
	*gplan.Holidays
}

func main() {

	h := &holidaysExt{
		Holidays: gplan.NewHolidays(time.Now(), time.Now()),
		Name:     "Pollas",
	}
	h2 := &holidaysExt{
		Holidays: gplan.NewHolidays(time.Now(), time.Now()),
		Name:     "Pollas2",
	}

	r := gplan.NewResource("Antoño", "Hueso", "backend", time.Now().AddDate(0, 1, 0), h, h2)

	fmt.Printf("%v\n", r.GetHolidays()[0].GetFrom())
}
