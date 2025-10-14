package types

type Type string

const (
	TypeNode      Type = "node"
	TypeExpansion Type = "expansion"
)

var TypeEnum = map[Type]float64{
	TypeNode:      1,
	TypeExpansion: 2,
}

func (_t Type) Enum() float64 {
	return TypeEnum[_t]
}
