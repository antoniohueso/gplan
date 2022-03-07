package dateutil

import "time"

// IsEqual devuelve True si la fecha que recibe como parámetro dayA es igual que dayB
func IsEqual(dayA time.Time, dayB time.Time) bool {
	dayA = toLocalDate(dayA)
	dayB = toLocalDate(dayB)

	return dayA.Equal(dayB)
}

// IsLt devuelve True si la fecha que recibe como parámetro dayA es menor que dayB
func IsLt(dayA time.Time, dayB time.Time) bool {
	dayA = toLocalDate(dayA)
	dayB = toLocalDate(dayB)

	return dayA.Before(dayB)
}

// IsLte devuelve True si la fecha que recibe como parámetro dayA es menor o igual que dayB
func IsLte(dayA time.Time, dayB time.Time) bool {
	dayA = toLocalDate(dayA)
	dayB = toLocalDate(dayB)

	return dayA.Before(dayB) || dayA.Equal(dayB)
}

// IsGt devuelve True si la fecha que recibe como parámetro dayA es mayor que dayB
func IsGt(dayA time.Time, dayB time.Time) bool {
	dayA = toLocalDate(dayA)
	dayB = toLocalDate(dayB)

	return dayA.After(dayB)
}

// IsGte devuelve True si la fecha que recibe como parámetro dayA es mayor o igual que dayB
func IsGte(dayA time.Time, dayB time.Time) bool {
	dayA = toLocalDate(dayA)
	dayB = toLocalDate(dayB)

	return dayA.After(dayB) || dayA.Equal(dayB)
}

// IsBetween devuelve True si la fecha que recibe como parámetro day es mayor o igual que from y menor o igual que to
func IsBetween(day time.Time, from time.Time, to time.Time) bool {
	day = toLocalDate(day)
	from = toLocalDate(from)
	to = toLocalDate(to)

	return IsGte(day, from) && IsLte(day, to)
}

// toLocalDate devuelve una fecha convertida a local y si hora para poder ser comparada con otra fecha
func toLocalDate(d time.Time) time.Time {
	dateLocal := d.Local()
	return time.Date(dateLocal.Year(), dateLocal.Month(), dateLocal.Day(), 0, 0, 0, 0, time.Local)
}
