package main

import (
	"time"

	"com.github.antoniohueso/gplan/defmodel"
)

func main() {
	h := defmodel.NewHolidays(time.Now(), time.Now())
	prueba(h)
}

func prueba(p interface{}) {

}
