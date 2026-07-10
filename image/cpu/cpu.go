package cpu

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CPUMonitor struct {
	lastUsage int64
	lastTime  time.Time
	numCPUs   float64
}

func readUsageUsec() (int64, error) {
	data, err := os.ReadFile("/sys/fs/cgroup/cpu.stat")
	if err != nil {
		return 0, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "usage_usec") {
			fields := strings.Fields(line)
			return strconv.ParseInt(fields[1], 10, 64)
		}
	}
	return 0, fmt.Errorf("usage_usec não encontrado")
}

// retorna quantos CPUs o container tem direito (ex: 0.5, 1.0, 2.0)
func readCPULimit() (float64, error) {
	data, err := os.ReadFile("/sys/fs/cgroup/cpu.max")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(strings.TrimSpace(string(data)))
	if fields[0] == "max" {
		return -1, nil // sem limite -> usar número de CPUs da máquina
	}
	quota, _ := strconv.ParseFloat(fields[0], 64)
	period, _ := strconv.ParseFloat(fields[1], 64)
	return quota / period, nil
}

func NewCPUMonitor() *CPUMonitor {
	numCPUs, _ := readCPULimit()
	usage, _ := readUsageUsec()
	return &CPUMonitor{lastUsage: usage, lastTime: time.Now(), numCPUs: numCPUs}
}

func (m *CPUMonitor) Percent() float64 {
	usage, _ := readUsageUsec()
	now := time.Now()

	deltaUsage := float64(usage - m.lastUsage)
	deltaTime := now.Sub(m.lastTime).Microseconds()

	m.lastUsage = usage
	m.lastTime = now

	if deltaTime == 0 {
		return 0
	}
	return (deltaUsage / (float64(deltaTime) * m.numCPUs)) * 100
}