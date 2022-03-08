package main

import (
	"fmt"
	"time"

	"github.com/antoniohueso/gplan"
)

func main() {

	fechaUTC := time.Date(2022, 3, 30, 00, 00, 0, 0, time.Local).UTC()

	fmt.Printf("%s\n", fechaUTC)

	fmt.Printf("%s\n", gplan.CalculateLaborableDate(fechaUTC, -2, nil))

}
