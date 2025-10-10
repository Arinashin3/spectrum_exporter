package unit

import (
	"regexp"
	"strconv"
	"strings"
)

type IEC string
type Bytes int

const (
	UnitBytes Bytes = 1 << (10 * iota)
	UnitKilobytes
	UnitMegabytes
	UnitGigabytes
	UnitTerabytes
	UnitPetabytes
)

func (_b IEC) Bytes() Bytes {
	var num string
	var unit string
	for _, s := range strings.ToLower(string(_b)) {
		match, _ := regexp.MatchString("[0-9.,]", string(s))
		if match {
			num += string(s)
		} else {
			unit += string(s)
		}
	}
	f, _ := strconv.ParseFloat(num, 64)
	switch unit {
	case "b", "bytes":
		return Bytes(f)
	case "kb", "kilobyte", "kilobytes":
		return Bytes(f * float64(UnitKilobytes))
	case "mb", "megabyte", "megabytes":
		return Bytes(f * float64(UnitMegabytes))
	case "gb", "gigabyte", "gigabytes":
		return Bytes(f * float64(UnitGigabytes))
	case "tb", "terabyte", "terabytes":
		return Bytes(f * float64(UnitTerabytes))
	case "pb", "petabyte", "petabytes":
		return Bytes(f * float64(UnitPetabytes))
	}
	return 0
}

func (b Bytes) ToKiB() float64 {
	return float64(b) / float64(UnitKilobytes)
}
func (b Bytes) ToMiB() float64 {
	return float64(b) / float64(UnitMegabytes)
}
func (b Bytes) ToGiB() float64 {
	return float64(b) / float64(UnitGigabytes)
}
func (b Bytes) ToTiB() float64 {
	return float64(b) / float64(UnitTerabytes)
}
func (b Bytes) ToPiB() float64 {
	return float64(b) / float64(UnitPetabytes)
}
