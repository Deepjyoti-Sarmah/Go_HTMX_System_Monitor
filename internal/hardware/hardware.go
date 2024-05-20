package hardware

import (
	"runtime"
	"strconv"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

const megabyteDiv uint64 = 1024 * 1024
const gigabyteDiv uint64 = megabyteDiv * 1024

func GetSystemSection() (string, error) {
	runTimeOS := runtime.GOOS

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	hostStat, err := host.Info()
	if err != nil {
		return "", err
	}

	html := "<div class='system-data'><table class='table-auto'><tbody>"
	html = html + "<tr><td class='font-semibold text-lg'>Operating System:</td> <td><i class='fa fa-brands fa-linux'></i> " + runTimeOS + "</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Platform:</td><td> <i class='fa fa-brands fa-fedora'></i> " + hostStat.Platform + "</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Hostname:</td><td>" + hostStat.Hostname + "</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Number of processes running:</td><td>" + strconv.FormatUint(hostStat.Procs, 10) + "</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Total memory:</td><td>" + strconv.FormatUint(vmStat.Total/megabyteDiv, 10) + " MB</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Free memory:</td><td>" + strconv.FormatUint(vmStat.Free/megabyteDiv, 10) + " MB</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Percentage used memory:</td><td>" + strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%</td></tr></tbody></table>"

	html = html + "</div>"

	return html, nil
}

func GetDiskSection() (string, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", err
	}

	html := "<div class='disk-data'><table class='table-auto'><tbody>"
	html = html + "<tr><td class='font-semibold text-lg'>Total disk space:</td><td>" + strconv.FormatUint(diskStat.Total/gigabyteDiv, 10) + " GB</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Used disk space:</td><td>" + strconv.FormatUint(diskStat.Used/gigabyteDiv, 10) + " GB</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Free disk space:</td><td>" + strconv.FormatUint(diskStat.Free/gigabyteDiv, 10) + " GB</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Percentage disk space usage:</td><td>" + strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%</td></tr>"

	html = html + "</div>"

	return html, nil
}

// func GetCpuSection() (string, error) {
// 	cpuStat, err := cpu.Info()
// 	if err != nil {
// 		return "", err
// 	}

// 	percentage, err := cpu.Percent(0, true)
// 	if err != nil {
// 		return "", nil
// 	}

// 	html := "<div class='cpu-data'><table class='table-auto'><tbody>"

// 	if len(cpuStat) != 0 {
// 		html = html + "<tr><td class='font-semibold text-lg'>Model Name:</td><td>" + cpuStat[0].ModelName + "</td></tr>"
// 		html = html + "<tr><td class='font-semibold text-lg'>Family:</td><td>" + cpuStat[0].Family + "</td></tr>"
// 		html = html + "<tr><td class='font-semibold text-lg'>Speed:</td><td>" + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz</td></tr>"
// 	}

// 	firstCpus := percentage[:len(percentage)/2]
// 	secondCpus := percentage[len(percentage)/2:]

// 	html = html + "<tr><td class='font-semibold text-lg'>Cores: </td><td><div class='row mb-4'><div class='col-md-6'><table class='table-auto'><tbody>"
// 	for idx, cpupercent := range firstCpus {
// 		html = html + "<tr><td class='font-semibold text-lg'>CPU [" + strconv.Itoa(idx) + "]: " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%</td></tr>"
// 	}

// 	html = html + "</tbody></table></div><div class='col-md-6'><table class='table-auto'><tbody>"
// 	for idx, cpupercent := range secondCpus {
// 		html = html + "<tr><td class='font-semibold text-lg'>CPU [" + strconv.Itoa(idx+8) + "]: " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%</td></tr>"
// 	}

// 	html = html + "</tbody></table></div></div></td></tr></tbody></table></div>"

// 	return html, nil
// }

func GetCpuSection() (string, error) {
    cpuStat, err := cpu.Info()
    if err != nil {
        return "", err
    }

    percentage, err := cpu.Percent(0, true)
    if err != nil {
        return "", nil
    }

    html := "<div class='grid grid-cols-2 gap-2'>"

    if len(cpuStat) != 0 {
        html += "<span id='cpu-model'>" + cpuStat[0].ModelName + "</span>"
        html += "<span id='cpu-speed'>" + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz</span>"
    }

    for idx, cpupercent := range percentage {
        html += "<div class='bg-gray-200 rounded-md py-1 px-2'><span class='font-semibold'>CPU [" + strconv.Itoa(idx) + "]:</span> " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%</div>"
    }

    html += "</div>"

    return html, nil
}
