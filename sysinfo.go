package snake

type snakesysinfo struct {
	Input []string
}

// SysInfo ...
type SysInfo interface {
	Add()
	Disk()
	Points()
	Hostname()
	OS()
	PlatformVersion()
	VolumeName() //硬盘名
}
