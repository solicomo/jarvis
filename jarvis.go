package jarvis

const (
	DIR_PUBLIC = "/public/"
)

const (
	URL_INDEX    = "/"
	URL_DETECTOR = "/detector/"
	URL_REPORT   = "/report"
	URL_LOGIN    = "/login"
	URL_PING     = "/ping"
)

type MetricConfig struct {
	Name     string
	Type     string
	Detector string
	Params   []interface{}
	MD5      string
}

type NodeInfo struct {
	Name   string
	Type   string
	Addr   string
	OS     string
	CPU    string
	Core   string
	Mem    string
	Disk   string
	Uptime string
}

type Login struct {
	NodeInfo
}

type LoginRsp struct {
	ID      string
	Metrics map[string]MetricConfig
}

type Ping struct {
	ID     string
	Uptime string
}

type PingRsp struct {
	Metrics map[string]MetricConfig
}

type CommonRsp struct {
	Status string
}

const (
	COMMON_RSP_OK   = `{"status":"ok"}`
	COMMON_RSP_FAIL = `{"status":"fail"}`
)

type MetricReport struct {
	ID      string
	Metrics map[string]string
}
