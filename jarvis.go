package jarvis

const (
	PUBLIC_DIR = "/public/"
)

const (
	INDEX_URL  = "/"
	PUBLIC_URL = "/public/"
	REPORT_URL = "/report"
	LOGIN_URL  = "/login"
	PING_URL   = "/ping"
)

type MetricConfig struct {
	Name     string
	Type     string
	Detector string
	Params   []interface{}
	MD5      string
}

type Stat struct {
	ID      string
	Metrics map[string]string
}
