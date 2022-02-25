package main

import (
	"fmt"
	"time"

	"github.com/antoniohueso/gplan"
	"github.com/antoniohueso/gplan/sample"
)

func main() {

	fecha, err := time.Parse("2006-01-02", "2021-07-20")
	if err != nil {
		panic(err)
	}

	hfrom, _ := time.Parse("2006-01-02", "2021-07-21")
	hto, _ := time.Parse("2006-01-02", "2021-07-22")

	fmt.Println(gplan.CalculateLaborableDate(fecha, -2, []gplan.Holidays{sample.NewHolidays(hfrom, hto)}))

}
