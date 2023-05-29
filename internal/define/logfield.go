package define

type log string

const (
	// default 7 log field
	Message log = "x_msg"
	Level   log = "x_level"
	Date    log = "x_date"
	File    log = "x_file"
	Host    log = "x_host"
	PSM     log = "x_psm"
	Version log = "x_version" // int

	// 7 need filled log
	Tag      log = "x_tag"      // string
	Duration log = "x_duration" // float64
	// gin or rpc
	HandlerPath          log = "x_handler_path"   // string
	HandlerExecutePeriod log = "x_handler_period" // int64

	TraceId   PT = "x_trace_id"   // string
	StartTime PT = "x_start_time" // string(int64)
	CallFrom  PT = "x_call_from"  // string
	//IsTest    PT = "x_is_test"
)

func (pt log) String() string {
	return string(pt)
}

// the value it related must be string
type PT string

func (pt PT) String() string {
	return string(pt)
}

var (
	CtxPassThrough = []PT{TraceId, StartTime, CallFrom}
)
