package hardware

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

func GetSystemSection() (string, error)  {
	runTimeOs := runtime.GOOS

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	hostStat, err := host.Info() 
	if err != nil {
		return "", err
	}

	output := 
	fmt.Sprintf("Hostname : %s\nTotal Memory: %d\nMemory: %d\nOS: %s", hostStat.Hostname, vmStat.Total, vmStat.Used, runTimeOs)

	return output, nil
}

func GetCpuSection() (string, error) {
	cpuStat, err := cpu.Info()
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("CPU: %s\nCors: %d", cpuStat[0].ModelName, len(cpuStat))
	return output, nil
}

func GetDiskSection()  (string, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", err
	}

	output := fmt.Sprintf("Total Disk Space: %d\nFree Disk Space: %d", diskStat.Total, diskStat.Free)
	return output, nil
}