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

	html := "<div class='system-data gap-4 p-4'><table class='table-auto'><tbody>"
	html = html + "<tr><td class='font-semibold text-lg'>Operating System:</td> <td class='text-right'>" + runTimeOS + "</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Platform:</td><td class='text-right'> " + hostStat.Platform + "</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Hostname:</td><td class='text-right'>" + hostStat.Hostname + "</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Processes running:</td><td class='text-right'>" + strconv.FormatUint(hostStat.Procs, 10) + "</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Total memory:</td><td class='text-right'>" + strconv.FormatUint(vmStat.Total/megabyteDiv, 10) + " MB</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Free memory:</td><td class='text-right'>" + strconv.FormatUint(vmStat.Free/megabyteDiv, 10) + " MB</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Memory usage:</td><td class='text-right'>" + strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64) + "%</td></tr></tbody></table>"

	html = html + "</div>"

	return html, nil
}

func GetDiskSection() (string, error) {
	diskStat, err := disk.Usage("/")
	if err != nil {
		return "", err
	}

	html := "<div class='disk-data gap-4 p-4'><table class='table-auto'><tbody>"
	html = html + "<tr><td class='font-semibold text-lg'>Total disk space:</td><td class='text-right'>" + strconv.FormatUint(diskStat.Total/gigabyteDiv, 10) + " GB</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Used disk space:</td><td class='text-right'>" + strconv.FormatUint(diskStat.Used/gigabyteDiv, 10) + " GB</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Free disk space:</td><td class='text-right'>" + strconv.FormatUint(diskStat.Free/gigabyteDiv, 10) + " GB</td></tr>"
	html = html + "<tr><td class='font-semibold text-lg'>Space usage:</td><td class='text-right'>" + strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64) + "%</td></tr>"

	html = html + "</div>"

	return html, nil
}

func GetCpuSection() (string, error) {
    cpuStat, err := cpu.Info()
    if err != nil {
        return "", err
    }

    percentage, err := cpu.Percent(0, true)
    if err != nil {
        return "", nil
    }

    html := "<div class='cpu-data p-4 gap-4'>"

    if len(cpuStat) != 0 {
		html += "<p class='mb-2'><span class='font-semibold text-lg'>Model Name: </span> <span id='cpu-model' class='text-right'>" + cpuStat[0].ModelName + "</span></p>" 
        html += "<p class='mb-2'><span class='font-semibold text-lg'>Speed: </span> <span id='cpu-model' class='text-right'>" + strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64) + " MHz</span></p>"
    }
    html += "</div>"

	html += "<div class='cores-data grid grid-cols-2 md:grid-cols-4 gap-4 p-4'>"
    for idx, cpupercent := range percentage {
        html += "<div class='bg-gray-200 rounded-md py-1 px-2'><span class='font-semibold'>CPU [" + strconv.Itoa(idx) + "]:</span> " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%</div>"
    }

    html += "</div>"

    return html, nil
}

