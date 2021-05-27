package cpu

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/cpu"
	"math"
)

var (
	GetCpuStatusError error
	HistoryInfo       []cpu.TimesStat
	Register          metrics.Registry
)

func init() {
	Register = metrics.NewRegistry()
	//	logs.Info("CPU-Task-INIT-%v\n ", Register.GetAll())
	tk := toolbox.NewTask("GetCpuDetailedTask", "*/20 * * * * *", GetCpuDetailedTask)
	toolbox.AddTask("GetCpuDetailedTask", tk)
}

func GetCpuDetailedTask() error {
	//	logs.Info("CPU-Task-GetCpuDetailedTask-%v\n ", Register.GetAll())

	var Detailed []cpu.TimesStat
	if len(HistoryInfo) == 0 {
		HistoryInfo, GetCpuStatusError = cpu.Times(false)
	}
	Detailed, GetCpuStatusError = cpu.Times(false)
	if GetCpuStatusError != nil {
		return GetCpuStatusError
	}
	nowBusy := Detailed[0].User + Detailed[0].System + Detailed[0].Nice + Detailed[0].Iowait + Detailed[0].Irq +
		Detailed[0].Softirq + Detailed[0].Steal

	historyBusy := HistoryInfo[0].User + HistoryInfo[0].System + HistoryInfo[0].Nice + HistoryInfo[0].Iowait + HistoryInfo[0].Irq +
		HistoryInfo[0].Softirq + HistoryInfo[0].Steal

	if nowBusy <= historyBusy {
		metrics.GetOrRegisterGaugeFloat64("percent", Register).Update(0)
	} else {
		if Detailed[0].Total() <= HistoryInfo[0].Total() {
			metrics.GetOrRegisterGaugeFloat64("percent", Register).Update(100)
		} else {
			used := math.Min(100, math.Max(0, (nowBusy-historyBusy)/(Detailed[0].Total()-HistoryInfo[0].Total())*100))
			metrics.GetOrRegisterGaugeFloat64("percent", Register).Update(used)
		}
	}
	HistoryInfo = Detailed
	//	logs.Info("CPU-Task-GetCpuDetailedTask-HistoryInfo-%v\n ", HistoryInfo)
	return nil
}

func Count() int {
	count, err := cpu.Counts(true)
	if err != nil {
		logs.Error(err.Error())
		return 0
	}
	return count
}
