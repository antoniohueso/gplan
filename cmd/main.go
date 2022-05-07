package main

import (
	"fmt"
	"log"
	"time"

	"github.com/antoniohueso/gplan"
)

func main() {

	p, err := time.Parse(time.RFC3339, "2022-05-13T22:00:00.000+00:00")
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("%s %v", p, gplan.IsLaborableDay(p, nil))
}
