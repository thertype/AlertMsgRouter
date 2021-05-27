package collector

import (
	"github.com/shirou/gopsutil/process"
)

//type Pstatus struct {
//	Running int64
//	Zombie int64
//	Sleep int64
//}

func GetProcessesStatus() (float64, float64, float64) {
	status1, _ := process.Processes()

	//var pstatus Pstatus
	var running float64
	var zombie float64
	var sleep float64

	for _, pro := range status1 {
		b, _ := pro.Status()
		//logs.Info("GetStatusTask-b-%v\n ", b)

		switch b {
		case "R":
			running += 1
		case "Z":
			zombie += 1
		case "S":
			sleep += 1
		}
	}

	//logs.Info("GetStatusTask-Register.GetAll()-%v\n ", Register.GetAll())

	return running, zombie, sleep
}
