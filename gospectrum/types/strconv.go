package types

import (
	"strconv"
	"time"
)

type Bool string

func (_b Bool) Bool() bool {
	b, _ := strconv.ParseBool(string(_b))
	return b
}

type Number string

func (_n Number) Int() int {
	i, _ := strconv.Atoi(string(_n))
	return i
}

func (_n Number) Float() float64 {
	f, _ := strconv.ParseFloat(string(_n), 64)
	return f
}

type Timestamp string

const TimeLayout = "060102150405"

func (_t Timestamp) Time() time.Time {
	t, _ := time.ParseInLocation(TimeLayout, string(_t), time.Local)
	return t
}
