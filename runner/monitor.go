package runner

import (
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

type MonitorClient struct {
	Wait int `json:"Wait"`
}

// NewSystemClient the Monitor for machine usage,
// returns Cpu Total Usage Percent, Cpu Count and Mem Usage Percent.
func NewSystemClient(wait int) *MonitorClient {
	return &MonitorClient{
		Wait: wait,
	}
}

type UsageStates struct {
	CpuPercent float64 `json:"CpuPercent"`
	CpuCount   int     `json:"CpuCount"`
	MemPercent float64 `json:"MemPercent"`
}

func (s *MonitorClient) Read() (*UsageStates, error) {
	cpuPercent, err := cpu.Percent(time.Duration(s.Wait), false)
	if err != nil {
		return nil, err
	}
	cpuCount, err := cpu.Counts(false)
	if err != nil {
		return nil, err
	}
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}
	memPercent := vmStat.UsedPercent
	return &UsageStates{
		CpuPercent: cpuPercent[0],
		CpuCount:   cpuCount,
		MemPercent: memPercent,
	}, err
}
