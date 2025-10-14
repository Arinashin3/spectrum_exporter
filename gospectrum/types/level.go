package types

type RaidLevel string

const (
	RaidLevelRAID0  RaidLevel = "raid0"
	RaidLevelRAID1  RaidLevel = "raid1"
	RaidLevelRAID5  RaidLevel = "raid5"
	RaidLevelRAID6  RaidLevel = "raid6"
	RaidLevelRAID10 RaidLevel = "raid10"
)

var RaidLevelEnum = map[RaidLevel]float64{
	RaidLevelRAID0:  0,
	RaidLevelRAID1:  1,
	RaidLevelRAID5:  5,
	RaidLevelRAID6:  6,
	RaidLevelRAID10: 10,
}

func (_level RaidLevel) Enum() float64 {
	return RaidLevelEnum[_level]
}
