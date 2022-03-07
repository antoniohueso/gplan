package main

import (
	"fmt"
	"time"
)

func main() {

	fechaUTC := time.Date(2022, 3, 6, 23, 30, 0, 0, time.UTC)
	fechaLocal := fechaUTC.Local().Add(2 * time.Hour)

	fmt.Printf("%s\n", fechaUTC)
	fmt.Println(fechaUTC.Year())
	fmt.Println(fechaUTC.Month())
	fmt.Println(fechaUTC.Day())
	fmt.Println(fechaUTC.Weekday())

	fmt.Printf("%s\n", fechaLocal)

	fmt.Println(fechaLocal.Year())
	fmt.Println(fechaLocal.Month())
	fmt.Println(fechaLocal.Day())
	fmt.Println(fechaLocal.Weekday())

	fmt.Println(fechaUTC.Local())
	fmt.Println(fechaLocal.Local())

	f(fechaUTC)
	f(fechaLocal)

	fmt.Println("POLLAS", fechaUTC.Local())
	fmt.Println("POLLAS", fechaLocal.Local())

	d1 := time.Date(fechaUTC.Local().Year(), fechaUTC.Local().Month(), fechaUTC.Local().Day(), 0, 0, 0, 0, time.UTC)
	d2 := time.Date(fechaLocal.Local().Year(), fechaUTC.Local().Month(), fechaUTC.Local().Day(), 0, 0, 0, 0, time.UTC)

	fmt.Println(d1.Equal(d2))

	/*
		fecha := time.Date(2022, 3, 29, 23, 0, 0, 0, time.UTC)

		fmt.Println(fecha.Local())

		res := gplan.CalculateLaborableDate(fecha, -2, nil)

		fmt.Println(res.Weekday())

		fmt.Println(res.Local())
		fmt.Println(res)
	*/
}

func f(fe time.Time) {
	fe = time.Date(fe.Local().Year(), fe.Local().Month(), fe.Local().Day(), 0, 0, 0, 0, time.Local)
	fmt.Println(fe)
}
