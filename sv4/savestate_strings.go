// generated by stringer -type=SaveStateType,scenarioGoalType,ParkFlag -output=savestate_strings.go; DO NOT EDIT

package sv4

import "fmt"

const (
	_SaveStateType_name_0 = "TypeRCT"
	_SaveStateType_name_1 = "TypeAACF"
	_SaveStateType_name_2 = "TypeLL"
)

var (
	_SaveStateType_index_0 = [...]uint8{0, 7}
	_SaveStateType_index_1 = [...]uint8{0, 8}
	_SaveStateType_index_2 = [...]uint8{0, 6}
)

func (i SaveStateType) String() string {
	switch {
	case i == 108156:
		return _SaveStateType_name_0
	case i == 110001:
		return _SaveStateType_name_1
	case i == 120001:
		return _SaveStateType_name_2
	default:
		return fmt.Sprintf("SaveStateType(%d)", i)
	}
}

const _scenarioGoalType_name = "GoalGuestsAndRatingParkValueHaveFunCompetitionTenCoasters6ExcitementMaintainMonthlyRideIncomeTenCoasters7ExcitementFiveCoasters"

var _scenarioGoalType_index = [...]uint8{0, 19, 28, 35, 46, 68, 76, 93, 115, 127}

func (i scenarioGoalType) String() string {
	if i >= scenarioGoalType(len(_scenarioGoalType_index)-1) {
		return fmt.Sprintf("scenarioGoalType(%d)", i)
	}
	return _scenarioGoalType_name[_scenarioGoalType_index[i]:_scenarioGoalType_index[i+1]]
}

const _ParkFlag_name = "ParkOpenedProhibitLandModifitationsProhibitRemovingSceneryShowRealNamesProhibitAboveTreeLevelLowIntensityPeepsProhibitAdvertisingCheatsDetectedHighIntensityPeepsNoMoneyModeGuestHighDifficultyForcedFreeEntryRatingHighDifficultyProhibitToggleRealNames"

var _ParkFlag_map = map[ParkFlag]string{
	1:     _ParkFlag_name[0:10],
	4:     _ParkFlag_name[10:35],
	8:     _ParkFlag_name[35:58],
	16:    _ParkFlag_name[58:71],
	32:    _ParkFlag_name[71:93],
	64:    _ParkFlag_name[93:110],
	128:   _ParkFlag_name[110:129],
	256:   _ParkFlag_name[129:143],
	512:   _ParkFlag_name[143:161],
	2048:  _ParkFlag_name[161:172],
	4096:  _ParkFlag_name[172:191],
	8192:  _ParkFlag_name[191:206],
	16384: _ParkFlag_name[206:226],
	32768: _ParkFlag_name[226:249],
}

func (i ParkFlag) String() string {
	if str, ok := _ParkFlag_map[i]; ok {
		return str
	}
	return fmt.Sprintf("ParkFlag(%d)", i)
}
