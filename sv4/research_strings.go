// generated by stringer -type=ResearchRate,ResearchFlag -output=research_strings.go; DO NOT EDIT

package sv4

import "fmt"

const _ResearchRate_name = "NoResearchMinimumResearchNormalResearchMaximumResearch"

var _ResearchRate_index = [...]uint8{0, 10, 25, 39, 54}

func (i ResearchRate) String() string {
	if i >= ResearchRate(len(_ResearchRate_index)-1) {
		return fmt.Sprintf("ResearchRate(%d)", i)
	}
	return _ResearchRate_name[_ResearchRate_index[i]:_ResearchRate_index[i+1]]
}

const (
	_ResearchFlag_name_0 = "RollercoastersResearchThrillRidesResearch"
	_ResearchFlag_name_1 = "GentleRidesResearch"
	_ResearchFlag_name_2 = "ShopsResearch"
	_ResearchFlag_name_3 = "ThemingResearch"
	_ResearchFlag_name_4 = "RideImprovementsResearch"
)

var (
	_ResearchFlag_index_0 = [...]uint8{0, 22, 41}
	_ResearchFlag_index_1 = [...]uint8{0, 19}
	_ResearchFlag_index_2 = [...]uint8{0, 13}
	_ResearchFlag_index_3 = [...]uint8{0, 15}
	_ResearchFlag_index_4 = [...]uint8{0, 24}
)

func (i ResearchFlag) String() string {
	switch {
	case 1 <= i && i <= 2:
		i -= 1
		return _ResearchFlag_name_0[_ResearchFlag_index_0[i]:_ResearchFlag_index_0[i+1]]
	case i == 4:
		return _ResearchFlag_name_1
	case i == 8:
		return _ResearchFlag_name_2
	case i == 16:
		return _ResearchFlag_name_3
	case i == 32:
		return _ResearchFlag_name_4
	default:
		return fmt.Sprintf("ResearchFlag(%d)", i)
	}
}
