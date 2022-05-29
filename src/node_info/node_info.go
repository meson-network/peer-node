package node_info

import (
	"fmt"
	"time"

	"github.com/meson-network/peer-node/src/cdn_cache_folder"
	"github.com/meson-network/peer_common/heart_beat"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type NodeInfo struct {
	Stor_total_bytes int64 `json:"stor_total_bytes"`
	Stor_used_bytes  int64 `json:"stor_used_bytes"`

	heart_beat.HardwareInfo
}

var info = &NodeInfo{}

func GetNodeInfo() *NodeInfo {
	refreshInfo()
	return info
}

func refreshInfo() {
	//hardware
	//cpu
	if c, err := cpu.Info(); err == nil {
		info.Cpu = c[0].ModelName
		info.Cpu_count = int64(len(c))
	}

	//cpu usage
	if percent, err := cpu.Percent(time.Second, false); err == nil || len(percent) > 0 {
		info.Cpu_percentage = percent[0]
	}

	//system info
	if h, err := host.Info(); err == nil {
		info.Op_sys = fmt.Sprintf("%v:%v(%v):%v", h.OS, h.Platform, h.PlatformFamily, h.PlatformVersion)
	}

	//memory
	if v, err := mem.VirtualMemory(); err == nil {
		info.Mem_total_bytes = int64(v.Total)
		info.Mem_used_bytes = int64(v.Used)
	}

	//disk
	if d, err := disk.Usage("/"); err == nil {
		info.Disk_total_bytes = int64(d.Total)
		info.Disk_used_bytes = int64(d.Used)
	}

	//cdn cache space
	info.Stor_total_bytes = cdn_cache_folder.GetInstance().Cache_provide_size
	info.Stor_used_bytes = cdn_cache_folder.GetInstance().Cache_used_size
}
