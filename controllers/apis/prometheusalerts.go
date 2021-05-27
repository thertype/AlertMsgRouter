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
	tpltext_Tpl := "{{ range $k,$v:=.alerts }}{{if eq $v.status \"resolved\"}}恢复信息故障主机IP{{$v.labels.instance}}{{$v.annotations.description}}{{else}}告警信息故障主机IP{{$v.labels.instance}}{{$v.annotations.description}}{{end}}{{ end }}"

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
	//微信渠道
	case "wxwork":
		Wxurl := strings.Split(pwxurl, ",")
		for _, url := range Wxurl {
			ret += PostToWeiXin(message, url, logsign)
		}

	////钉钉渠道
	//case "dd":
	//	Ddurl := strings.Split(pddurl, ",")
	//	for _, url := range Ddurl {
	//		ret += PostToDingDing(Title+"告警消息", message, url, logsign)
	//	}
	//
	////飞书渠道
	//case "fs":
	//	Fsurl := strings.Split(pfsurl, ",")
	//	for _, url := range Fsurl {
	//		ret += PostToFS(Title+"告警消息", message, url, logsign)
	//	}
	//
	////腾讯云短信
	//case "txdx":
	//	ret = PostTXmessage(message, pphone, logsign)
	////华为云短信
	//case "hwdx":
	//	ret = ret + PostHWmessage(message, pphone, logsign)
	////百度云短信
	//case "bddx":
	//	ret = ret + PostBDYmessage(message, pphone, logsign)
	////阿里云短信
	//case "alydx":
	//	ret = ret + PostALYmessage(message, pphone, logsign)
	////腾讯云电话
	//case "txdh":
	//	ret = PostTXphonecall(message, pphone, logsign)
	////阿里云电话
	//case "alydh":
	//	ret = ret + PostALYphonecall(message, pphone, logsign)
	////容联云电话
	//case "rlydh":
	//	ret = ret + PostRLYphonecall(message, pphone, logsign)
	////七陌短信
	//case "7moordx":
	//	ret = ret + Post7MOORmessage(message, pphone, logsign)
	////七陌语音电话
	//case "7moordh":
	//	ret = ret + Post7MOORphonecall(message, pphone, logsign)
	////邮件
	//case "email":
	//	ret = ret + SendEmail(message, email, logsign)
	//// Telegram
	//case "tg":
	//	ret = ret + SendTG(message, logsign)
	// Workwechat
	//case "workwechat":
	//	ret = ret + SendWorkWechat(ptouser,ptoparty,ptotag,message, logsign)
	////百度Hi(如流)
	//case "rl":
	//	ret += PostToRuLiu(pgroupid, message, beego.AppConfig.String("BDRL_URL"), logsign)
	//异常参数
	default:
		ret = "参数错误"
	}
	return ret
}
