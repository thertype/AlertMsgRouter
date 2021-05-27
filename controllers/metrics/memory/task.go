package memory

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
	"github.com/rcrowley/go-metrics"
	"github.com/shirou/gopsutil/mem"
)

var (
	StatusError error
	Register    metrics.Registry
)

func init() {
	Register = metrics.NewRegistry()
	tk := toolbox.NewTask("GetMemDetailedTask", "*/20 * * * * *", GetMemDetailedTask)
	toolbox.AddTask("GetMemDetailedTask", tk)
}

func GetMemDetailedTask() error {
	Register.UnregisterAll()
	status, err := mem.VirtualMemory()
	if err != nil {
		StatusError = err
		return err
	}
	metrics.GetOrRegisterGaugeFloat64("percent", Register).Update(status.UsedPercent)
	metrics.GetOrRegisterGauge("total", Register).Update(int64(status.Total))
	metrics.GetOrRegisterGauge("used", Register).Update(int64(status.Used))

	//logs.Info("GetMemDetailedTask-Task-Register.GetAll()-%v\n ", Register.GetAll())
	//logs.Info("GetMemDetailedTask-Task-status-%v\n ", status)

	return nil

}

func TotalSize() uint64 {
	info, err := mem.VirtualMemory()
	if err != nil {
		logs.Error(err)
		return 0
	} else {
		return info.Total
	}
}
