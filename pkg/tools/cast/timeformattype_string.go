package cast

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[timeFormatNoTimezone-0]
	_ = x[timeFormatNamedTimezone-1]
	_ = x[timeFormatNumericTimezone-2]
	_ = x[timeFormatNumericAndNamedTimezone-3]
	_ = x[timeFormatTimeOnly-4]
}

const _timeFormatType_name = "timeFormatNoTimezonetimeFormatNamedTimezonetimeFormatNumericTimezonetimeFormatNumericAndNamedTimezonetimeFormatTimeOnly"

var _timeFormatType_index = [...]uint8{0, 20, 43, 68, 101, 119}

func (i timeFormatType) String() string {
	if i < 0 || i >= timeFormatType(len(_timeFormatType_index)-1) {
		return "timeFormatType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _timeFormatType_name[_timeFormatType_index[i]:_timeFormatType_index[i+1]]
}
