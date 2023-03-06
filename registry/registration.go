package registry

type Registration struct {
	// 服务名
	ServiceName ServiceName
	// 服务URL
	ServiceURL       string
	RequiredServices []ServiceName
	// 接收服务状态信息
	ServiceUpdateURL string
	// 心跳检测
	HeartbeatURL string
}

type ServiceName string

const (
	LogService     = ServiceName("LogService")
	GradingService = ServiceName("GradingService")
	PortalService  = ServiceName("Portal")
)

type patchEntry struct {
	Name ServiceName
	URL  string
}

type patch struct {
	Added   []patchEntry
	Removed []patchEntry
}
