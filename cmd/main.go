package main

import (
	"fmt"
	"time"

	"github.com/antoniohueso/gplan"
)

func main() {

	fechaUTC := time.Date(2022, 3, 10, 00, 00, 0, 0, time.Local).UTC()

	fmt.Printf("%d\n", gplan.CalculateLaborableDays(time.Now(), fechaUTC, nil))

}
