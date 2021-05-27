package cpu

import (
	"AlertMsgRouter/controllers/metrics/base"
	"github.com/astaxie/beego/logs"
)

type Controller struct {
	base.Controller
}

func (o *Controller) Get() {
	logs.Info("CPU-Get-%v\n ", Register.GetAll())
	o.Data["json"] = Register.GetAll()
	o.ServeJSON()
}
