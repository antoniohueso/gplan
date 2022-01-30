package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"com.github.antoniohueso/gplan"
	"com.github.antoniohueso/gplan/defmodel"
)

type arr []interface{}

func main() {

	h := defmodel.NewHolidays(time.Now(), time.Now())

	r := defmodel.NewResource("ahg", "Antonio Hueso", "backend", time.Now(), []gplan.Holidays{h})

	b, err := json.Marshal(r)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(b))

}
