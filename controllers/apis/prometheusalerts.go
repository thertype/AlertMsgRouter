package apis

import (
	"AlertMsgRouter/common"
	collector "AlertMsgRouter/controllers/metrics"
	"bytes"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/common/model"
	"strings"
	"text/template"
	"time"
)

type PrometheusAlertController struct {
	beego.Controller
}

type KV map[string]string

// Alert holds one alert for notification templates.
type Alert struct {
	Status       string    `json:"status"`
	Labels       KV        `json:"labels"`
	Annotations  KV        `json:"annotations"`
	StartsAt     time.Time `json:"startsAt"`
	EndsAt       time.Time `json:"endsAt"`
	GeneratorURL string    `json:"generatorURL"`
	Fingerprint  string    `json:"fingerprint"`
}

type Alerts []Alert

type Data struct {
	Receiver string `json:"receiver"`
	Status   string `json:"status"`
	Alerts   Alerts `json:"alerts"`

	GroupLabels       KV `json:"groupLabels"`
	CommonLabels      KV `json:"commonLabels"`
	CommonAnnotations KV `json:"commonAnnotations"`

	ExternalURL string `json:"externalURL"`
}

type WebhookMessage Data

// Firing returns the subset of alerts that are firing.
func (as Alerts) Firing() []Alert {
	res := []Alert{}
	for _, a := range as {
		if a.Status == string(model.AlertFiring) {
			res = append(res, a)
		}
	}
	return res
}

// Resolved returns the subset of alerts that are resolved.
func (as Alerts) Resolved() []Alert {
	res := []Alert{}
	for _, a := range as {
		if a.Status == string(model.AlertResolved) {
			res = append(res, a)
		}
	}
	return res
}

func (c *PrometheusAlertController) PrometheusAlert() {
	logsign := "[" + common.LogsSign() + "]"
	var p_json interface{}

	logs.Info("PrometheusAlert---RequestBody--: %+v\n ", string(c.Ctx.Input.RequestBody))
	//logs.Debug(logsign, strings.Replace(string(c.Ctx.Input.RequestBody), "\n", "", -1))
	json.Unmarshal(c.Ctx.Input.RequestBody, &p_json)

	P_wxurl := "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=b4d28872-6e5e-4990-a522-c69121966684"
	//P_wxurl := ""

	P_type := "wxwork"
	//P_tpl := c.Input().Get("tpl")
	//P_ddurl := c.Input().Get("ddurl")
	//P_wxurl := c.Input().Get("wxurl")
	//P_fsurl := c.Input().Get("fsurl")
	//P_phone := c.Input().Get("phone")

	funcMap := template.FuncMap{
		"GetCSTtime": common.GetCSTtime,
		"TimeFormat": common.TimeFormat,
	}
	message := ""
	buf := new(bytes.Buffer)
	//tpltext_Tpl := "{{ $var := .externalURL}}{{ range $k,$v:=.alerts }}"
	tpltext_Tpl := "{{ range $k,$v:=.alerts }}{{if eq $v.status \"resolved\"}}????????????????????????IP{{$v.labels.instance}}{{$v.annotations.description}}{{else}}????????????????????????IP{{$v.labels.instance}}{{$v.annotations.description}}{{end}}{{ end }}"

	tpl, err := template.New("").Funcs(funcMap).Parse(tpltext_Tpl)

	if err != nil {
		logs.Error(logsign, err.Error())
		message = err.Error()
	} else {
		tpl.Execute(buf, p_json)
		logs.Info("PrometheusAlert---buf.String()--: %+v\n ", buf.String())
		message = SendMessagePrometheusAlert(buf.String(), P_type, P_wxurl, logsign)
	}

	c.Data["json"] = message
	c.ServeJSON()

	//c.Data["json"] = &common.Res{
	//	Code: 0,
	//	Msg:  "",
	//	Data: p_json,
	//}
	//c.ServeJSON()
}

func SendMessagePrometheusAlert(message, ptype, pwxurl, logsign string) string {
	//Title := beego.AppConfig.String("title")
	ret := ""
	collector.AlertsFromCounter.WithLabelValues("PrometheusAlert", message, "", "", "").Add(1)
	switch ptype {
	//????????????
	case "wxwork":
		Wxurl := strings.Split(pwxurl, ",")
		for _, url := range Wxurl {
			ret += PostToWeiXin(message, url, logsign)
		}

	////????????????
	//case "dd":
	//	Ddurl := strings.Split(pddurl, ",")
	//	for _, url := range Ddurl {
	//		ret += PostToDingDing(Title+"????????????", message, url, logsign)
	//	}
	//
	////????????????
	//case "fs":
	//	Fsurl := strings.Split(pfsurl, ",")
	//	for _, url := range Fsurl {
	//		ret += PostToFS(Title+"????????????", message, url, logsign)
	//	}
	//
	////???????????????
	//case "txdx":
	//	ret = PostTXmessage(message, pphone, logsign)
	////???????????????
	//case "hwdx":
	//	ret = ret + PostHWmessage(message, pphone, logsign)
	////???????????????
	//case "bddx":
	//	ret = ret + PostBDYmessage(message, pphone, logsign)
	////???????????????
	//case "alydx":
	//	ret = ret + PostALYmessage(message, pphone, logsign)
	////???????????????
	//case "txdh":
	//	ret = PostTXphonecall(message, pphone, logsign)
	////???????????????
	//case "alydh":
	//	ret = ret + PostALYphonecall(message, pphone, logsign)
	////???????????????
	//case "rlydh":
	//	ret = ret + PostRLYphonecall(message, pphone, logsign)
	////????????????
	//case "7moordx":
	//	ret = ret + Post7MOORmessage(message, pphone, logsign)
	////??????????????????
	//case "7moordh":
	//	ret = ret + Post7MOORphonecall(message, pphone, logsign)
	////??????
	//case "email":
	//	ret = ret + SendEmail(message, email, logsign)
	//// Telegram
	//case "tg":
	//	ret = ret + SendTG(message, logsign)
	// Workwechat
	//case "workwechat":
	//	ret = ret + SendWorkWechat(ptouser,ptoparty,ptotag,message, logsign)
	////??????Hi(??????)
	//case "rl":
	//	ret += PostToRuLiu(pgroupid, message, beego.AppConfig.String("BDRL_URL"), logsign)
	//????????????
	default:
		ret = "????????????"
	}
	return ret
}
