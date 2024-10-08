// Code generated by "stringer -type=FaultLevel"; DO NOT EDIT.

package fault

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[FATAL-0]
	_ = x[ERROR-1]
	_ = x[WARNING-2]
	_ = x[NOTICE-3]
	_ = x[DEBUG-4]
}

const _FaultLevel_name = "FATALERRORWARNINGNOTICEDEBUG"

var _FaultLevel_index = [...]uint8{0, 5, 10, 17, 23, 28}

func (i FaultLevel) String() string {
	if i < 0 || i >= FaultLevel(len(_FaultLevel_index)-1) {
		return "FaultLevel(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FaultLevel_name[_FaultLevel_index[i]:_FaultLevel_index[i+1]]
}
