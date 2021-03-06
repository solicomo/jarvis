// +build linux

package detector

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"syscall"
	"unsafe"

	"github.com/likexian/host-stat-go"
	"github.com/toolkits/nux"
)

func int8ToString(bs []int8) string {
	return strings.TrimRight(string(*(*[]byte)(unsafe.Pointer(&bs))), "\x00")
}

func Round(data float64, precision int) (result float64) {
	pow := math.Pow(10, float64(precision))
	digit := pow * data
	_, div := math.Modf(digit)

	if div >= 0.5 {
		result = math.Ceil(digit)
	} else {
		result = math.Floor(digit)
	}
	result = result / pow

	return
}

func (d Detector) BaseMetrics() (result string, err error) {

	var bm struct {
		Load      string
		CPURate   string
		MemRate   string
		SwapRate  string
		DiskRate  string
		DiskRead  string
		DiskWrite string
		NetRead   string
		NetWrite  string
	}

	bm.Load, _ = d.Load()
	bm.CPURate, _ = d.CPURate()
	bm.MemRate, _ = d.MemRate()
	bm.SwapRate, _ = d.SwapRate()
	bm.DiskRate, _ = d.DiskRate()
	bm.DiskRead, _ = d.DiskRead()
	bm.DiskWrite, _ = d.DiskWrite()
	bm.NetRead, _ = d.NetRead()
	bm.NetWrite, _ = d.NetWrite()

	data, err := json.Marshal(bm)

	if err != nil {
		log.Println("[ERRO]", err)
		result = err.Error()
		return
	}

	result = string(data[:])

	return
}

func (Detector) OSVer(params ...string) (result string, err error) {
	var u syscall.Utsname
	if err = syscall.Uname(&u); err != nil {
		return
	}

	sysName := int8ToString(u.Sysname[:])
	nodName := int8ToString(u.Nodename[:])
	release := int8ToString(u.Release[:])
	version := int8ToString(u.Version[:])
	machine := int8ToString(u.Machine[:])
	domName := int8ToString(u.Domainname[:])

	if len(params) == 0 {
		result = sysName + " " + nodName + " " + release + " " + version + " " + machine + " " + domName
	}

	for _, p := range params {
		switch p {
		case "s":
			result += sysName + " "
		case "n":
			result += nodName + " "
		case "r":
			result += release + " "
		case "v":
			result += version + " "
		case "m":
			result += machine + " "
		case "o":
			result += domName + " "
		}
	}
	return
}

func (Detector) Uptime() (result string, err error) {

	days, hours, mins, err := nux.SystemUptime()

	if err == nil {
		result = fmt.Sprintf("%v days, %v:%v", days, hours, mins)
	}

	//	if upt, err := host_stat.GetUptimeStat(); err == nil {
	//		result = fmt.Sprintf("%v", uint64(upt.Uptime))
	//	}

	return
}

func (Detector) Load() (result string, err error) {

	if load, err := host_stat.GetLoadStat(); err == nil {
		result = fmt.Sprintf("%v, %v, %v", load.LoadNow, load.LoadPre, load.LoadFar)
	}

	return
}

func (Detector) CPUName() (result string, err error) {

	if ci, err := host_stat.GetCPUInfo(); err == nil {
		result = ci.ModelName
	}

	return
}

func (Detector) CPUCore() (result string, err error) {

	if ci, err := host_stat.GetCPUInfo(); err == nil {
		result = fmt.Sprintf("%v", ci.CoreCount)
	}

	return
}

func (Detector) CPURate() (result string, err error) {

	if cs, err := host_stat.GetCPUStat(); err == nil {
		result = fmt.Sprintf("%v", cs.UserRate)
	}

	return
}

func (Detector) MemSize() (result string, err error) {

	if ms, err := host_stat.GetMemStat(); err == nil {
		result = fmt.Sprintf("%v", ms.MemTotal)
	}

	return
}

func (Detector) MemRate() (result string, err error) {

	if ms, err := host_stat.GetMemStat(); err == nil {
		result = fmt.Sprintf("%v", ms.MemRate)
	}

	return
}

func (Detector) SwapRate() (result string, err error) {

	if ms, err := host_stat.GetMemStat(); err == nil {
		result = fmt.Sprintf("%v", ms.SwapRate)
	}

	return
}

func (Detector) DiskSize() (result string, err error) {

	if ds, err := host_stat.GetDiskStat(); err == nil {

		disk_total := uint64(0)

		for _, v := range ds {
			disk_total += v.Total
		}

		result = fmt.Sprintf("%v", disk_total)
	}

	return
}

func (Detector) DiskRate() (result string, err error) {

	if ds, err := host_stat.GetDiskStat(); err == nil {

		disk_total := uint64(0)
		disk_used := uint64(0)

		for _, v := range ds {
			disk_total += v.Total
			disk_used += v.Used
		}

		result = fmt.Sprintf("%v", Round(float64(disk_used)/float64(disk_total), 2))
	}

	return
}

func (Detector) DiskRead() (result string, err error) {

	if is, err := host_stat.GetIOStat(); err == nil {

		disk_read := uint64(0)

		for _, v := range is {
			disk_read += v.ReadBytes / 1024
		}

		result = fmt.Sprintf("%v", disk_read)
	}

	return
}

func (Detector) DiskWrite() (result string, err error) {
	if is, err := host_stat.GetIOStat(); err == nil {

		disk_write := uint64(0)

		for _, v := range is {
			disk_write += v.WriteBytes / 1024
		}

		result = fmt.Sprintf("%v", disk_write)
	}

	return
}

func (Detector) NetRead() (result string, err error) {

	if ns, err := host_stat.GetNetStat(); err == nil {

		net_read := uint64(0)

		for _, v := range ns {
			if v.Device != "lo" {
				net_read += v.RXBytes / 1024
			}
		}

		result = fmt.Sprintf("%v", net_read)
	}

	return
}

func (Detector) NetWrite() (result string, err error) {

	if ns, err := host_stat.GetNetStat(); err == nil {

		net_write := uint64(0)

		for _, v := range ns {
			if v.Device != "lo" {
				net_write += v.TXBytes / 1024
			}
		}

		result = fmt.Sprintf("%v", net_write)
	}

	return
}
