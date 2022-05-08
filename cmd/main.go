package main

import (
	"fmt"
	"time"
)

func main() {

	f1 := time.Date(2022, time.May, 14, 0, 0, 0, 0, time.Local)

	fmt.Println(f1, f1.Weekday())
	fmt.Println(f1.UTC(), f1.UTC().Weekday())

}
