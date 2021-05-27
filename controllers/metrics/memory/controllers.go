package memory

import (
	"AlertMsgRouter/controllers/metrics/base"
	"github.com/astaxie/beego"
)

//func init() {
//	beego.Router("/metrics/memory", &Controller{})
//}

func init() {
	beego.Router("/memory", &Controller{})

}

type Controller struct {
	base.Controller
}

func (o *Controller) Get() {
	//logs.Info("memory-Get-%v\n ", Register.GetAll())
	o.Data["json"] = Register.GetAll()
	o.ServeJSON()
}
