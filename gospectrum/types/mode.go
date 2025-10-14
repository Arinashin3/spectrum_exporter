package types

type Mode string

const (
	ModeUnmanaged Mode = "unmanaged"
	ModeManaged   Mode = "managed"
	ModeImage     Mode = "image"
	ModeArray     Mode = "array"
)

var ModeEnum = map[Mode]float64{
	ModeUnmanaged: 0,
	ModeManaged:   1,
	ModeImage:     2,
	ModeArray:     3,
}

func (_s Mode) Enum() float64 {
	return ModeEnum[_s]
}
