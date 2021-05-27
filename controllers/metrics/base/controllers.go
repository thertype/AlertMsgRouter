package base

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type Controller struct {
	beego.Controller
}

func (o *Controller) JsonMsg(code int, message string) {
	o.Data["json"] = map[string]interface{}{"code": code, "msg": message}
	o.ServeJSON()
}

func (o *Controller) DataJsonMsg(code int, total int64, data interface{}) {
	o.Data["json"] = map[string]interface{}{"code": code, "count": total, "data": data}
	o.ServeJSON()
}

func (o *Controller) LoadRequest(request interface{}) (status bool) {
	err := json.Unmarshal(o.Ctx.Input.RequestBody, request)
	if err != nil {
		o.JsonMsg(2, "非法请求")
		logs.Error(err)
		return false
	}
	return true
}
