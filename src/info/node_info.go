package info

type NodeInfo struct {
	Port int `json:"port"`

	Cpu_cores      uint64  `json:"cpu_cores"`
	Cpu_percentage float64 `json:"cpu_percentage"`
	Mem_total      uint64  `json:"mem_total"`
	Mem_used       uint64  `json:"mem_used"`
	Op_sys         string  `json:"op_sys"`
	Disk_total     uint64  `json:"disk_total"`
	Disk_available uint64  `json:"disk_available"`
}

func GetMachineStatus() {
	//if s.Status == nil {
	//	s.Status = &meson_msg.TerminalStatesMsg{}
	//}
	//
	//if h, err := host.Info(); err == nil {
	//	s.Status.OS = fmt.Sprintf("%v:%v(%v):%v", h.OS, h.Platform, h.PlatformFamily, h.PlatformVersion)
	//}
	//
	//if s.Status.CPU == "" {
	//	if c, err := cpu.Info(); err == nil {
	//		s.Status.CPU = c[0].ModelName
	//	}
	//}
	//
	////need update data
	////memory
	//if v, err := mem.VirtualMemory(); err == nil {
	//	s.Status.MemTotal = int64(v.Total)
	//	s.Status.MemAvailable = int64(v.Available)
	//}
	//
	////cpu usage
	//if percent, err := cpu.Percent(time.Second, false); err != nil || len(percent) <= 0 {
	//	basic.Logger.Debugln("failed to get cup usage", "err", err)
	//} else {
	//	s.Status.CpuUsage = percent[0]
	//}
	//
	//s.Status.Version = versionMgr.GetInstance().CurrentVersion
	//
	////disk
	//total, used, _ := diskMgr.GetInstance().GetSpaceInfo()
	//s.Status.CdnDiskTotal = total
	//s.Status.CdnDiskUsed = used

}
