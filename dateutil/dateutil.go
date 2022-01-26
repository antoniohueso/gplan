package dateutil

import "time"

// IsLte devuelve True si la fecha que recibe como parámero dayA es menor o igual que dayB
func IsLte(dayA time.Time, dayB time.Time) bool {
	return dayA.Before(dayB) || dayA.Equal(dayB)
}

// IsGte devuelve True si la fecha que recibe como parámero dayA es mayor o igual que dayB
func IsGte(dayA time.Time, dayB time.Time) bool {
	return dayA.After(dayB) || dayA.Equal(dayB)
}

// IsBetween devuelve True si la fecha que recibe como parámero day es mayor o igual que from y menor o igual que to
func IsBetween(day time.Time, from time.Time, to time.Time) bool {
	return IsGte(day, from) && IsLte(day, to)
}
