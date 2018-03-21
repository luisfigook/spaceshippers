package main

var statInformation map[StatType]StatInfo

type SystemType int

const (
	SYS_PROPULSION SystemType = iota
	SYS_POWER
	SYS_COMPUTER
	SYS_NAVIGATION
	SYS_SCANNER
	SYS_LIFESUPPORT
	SYS_WEAPON
	SYS_SHIELD
	SYS_MAXSYSTEMS
)

type StatType int

const (
	//SYS_PROPULSION stats
	STAT_SUBLIGHT_THRUST StatType = iota
	STAT_SUBLIGHT_FUELUSE
	STAT_SUBLIGHT_POWER
	STAT_FTL_THRUST
	STAT_FTL_FUELUSE
	STAT_FTL_POWER

	//SYS_NAVIGATION stats
	STAT_NAV_POWER
	STAT_NAV_COMPUTER

	STAT_MAXSTATS
)

type StatInfo struct {
	Name        string
	Description string
	System      SystemType
	Stat        StatType
}

type RoomStat struct {
	Stat     StatType
	Modifier int
}

func (rs RoomStat) GetName() string {
	return statInformation[rs.Stat].Name
}

func (rs RoomStat) GetDesc() string {
	return statInformation[rs.Stat].Description
}

func (rs RoomStat) GetSystem() SystemType {
	return statInformation[rs.Stat].System
}

func init() {
	statInformation = make(map[StatType]StatInfo)

	statInformation[STAT_SUBLIGHT_THRUST] = StatInfo{
		Name:        "Sublight Thrust Factor",
		Description: "The thrust provided by the engine for sublight-travel",
		System:      SYS_PROPULSION,
		Stat:        STAT_SUBLIGHT_THRUST,
	}
	statInformation[STAT_SUBLIGHT_FUELUSE] = StatInfo{
		Name:        "Sublight Fuel Use",
		Description: "The fuel used by the engine for sublight-travel",
		System:      SYS_PROPULSION,
		Stat:        STAT_SUBLIGHT_FUELUSE,
	}
	statInformation[STAT_SUBLIGHT_POWER] = StatInfo{
		Name:        "Sublight Power Draw",
		Description: "The power used by the engine for sublight-travel",
		System:      SYS_PROPULSION,
		Stat:        STAT_SUBLIGHT_POWER,
	}
	statInformation[STAT_FTL_THRUST] = StatInfo{
		Name:        "FTL Thrust Factor",
		Description: "The thrust provided by the engine for FTL-travel",
		System:      SYS_PROPULSION,
		Stat:        STAT_FTL_THRUST,
	}
	statInformation[STAT_FTL_FUELUSE] = StatInfo{
		Name:        "FTL Fuel Use",
		Description: "The fuel used by the engine for FTL-travel",
		System:      SYS_PROPULSION,
		Stat:        STAT_FTL_FUELUSE,
	}
	statInformation[STAT_FTL_POWER] = StatInfo{
		Name:        "FTL Power Draw",
		Description: "The power used by the engine for FTL-travel",
		System:      SYS_PROPULSION,
		Stat:        STAT_FTL_POWER,
	}
}
