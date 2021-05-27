package common

import (
	"bytes"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"html/template"
	tmplhtml "html/template"
	"io/ioutil"
	"regexp"
	"strings"
	tmpltext "text/template"
	"time"
)

var logger log.Logger

// 定义接收的数据格式
type Alerts struct {
	Status          string      `json:"status,omitempty"`
	Labels          Labels      `json:"labels,omitempty"`
	Annotations     Annotations `json:"annotations,omitempty"`
	StartsAt        time.Time   `json:"startsAt,omitempty"`
	EndsAt          time.Time   `json:"endsAt,omitempty"`
	GeneratorURL    string      `json:"generatorURL,omitempty"`
	AlertManagerURL string      `json:"alertManagerURL,omitempty"`
	Nodeip          string      `json:"nodeip,omitempty"`
	Podip           string      `json:"podip,omitempty"`
	Timestamp       string      `json:"timestamp,omitempty"`
	Develop         string      `json:"develop,omitempty"`
	Manager         string      `json:"manager,omitempty"`
	Ops             string      `json:"ops,omitempty"`
	Ops_backup      string      `json:"ops_backup,omitempty"`
	Group_name      string      `json:"group_name,omitempty"`
	//Develop_list    string      `json:"develop_list,omitempty"`
	Duration      string `json:"duration,omitempty"`
	DurationHuman string `json:"duration_human,omitempty"`
	RetryCount    string `json:"retry_count,omitempty"`
	Ticket_id     string `json:"ticket_id,omitempty"`
}

type Labels struct {
	AlertName      template.HTML `json:"alertname,omitempty"`
	Cluster        string        `json:"cluster,omitempty"`
	Replica        string        `json:"replica,omitempty"`
	Namespace      string        `json:"namespace,omitempty"`
	Class          string        `json:"class,omitempty"`
	Project_name   string        `json:"project_name,omitempty"`
	App            string        `json:"app,omitempty"`
	Env            string        `json:"env,omitempty"`
	ConsumerGroup  string        `json:"consumergroup,omitempty"`
	Topic          string        `json:"topic,omitempty"`
	Container      string        `json:"container,omitempty"`
	Deployment     string        `json:"deployment,omitempty"`
	Unique_appname string        `json:"unique_appname,omitempty"`
	//Unique_appname_list string        `json:"unique_appname_list,omitempty"`
	Statefulset    string `json:"statefulset,omitempty"`
	Daemonset      string `json:"daemonset,omitempty"`
	Node           string `json:"node,omitempty"`
	Name           string `json:"name,omitempty"`
	Ip_wjs         string `json:"ip_wjs,omitempty"`
	Instance       string `json:"instance,omitempty"`
	Username       string `json:"username,omitempty"`
	Mountpoint     string `json:"mountpoint,omitempty"`
	Fstype         string `json:"fstype,omitempty"`
	Mongodbcluster string `json:"mongodbcluster,omitempty"`
	Set            string `json:"set,omitempty"`
	State          string `json:"state,omitempty"`
	Group          string `json:"group,omitempty"`
	Ceph_daemon    string `json:"ceph_daemon,omitempty"`
	Pool_id        string `json:"pool_id,omitempty"`
	Device         string `json:"device,omitempty"`
	Pod            string `json:"pod,omitempty"`
	Desp           string `json:"desp,omitempty"`
	Job            string `json:"job,omitempty"`
	Reason         string `json:"reason,omitempty"`
	Phase          string `json:"phase,omitempty"`
	Severity       string `json:"severity,omitempty"`
	Type           string `json:"type,omitempty"`
}

type Annotations struct {
	Result string `json:"result,omitempty"`
}

type AlertList []Alerts

type MessageContent struct {
	Firing   string `json:"firing,omitempty"`
	Resolved string `json:"resolved,omitempty"`
}

//func TimeFormat(t time.Time) string {
//	return t.Format("2006-01-02 15:04:05")
//}

func Compute(index int) int {
	return index % 2
}

func TimeFormatNew(t time.Time) string {
	interval, _ := time.ParseDuration("+8h")
	return t.Add(interval).Format("2006-01-02 15:04:05")
}

func TimeConvertCST(t string) string {
	interval, _ := time.ParseDuration("+8h")
	tstr, _ := time.Parse(time.RFC3339, t)
	return tstr.Add(interval).Format("2006-01-02 15:04:05")
}

func DefaultFuncs() map[string]interface{} {
	type FuncMap map[string]interface{}
	var DefaultFuncs = FuncMap{
		"compute":        Compute,
		"timeFormat":     TimeFormat,
		"timeConvertCST": TimeConvertCST,
		"timeFormatNew":  TimeFormatNew,
		"toUpper":        strings.ToUpper,
		"toLower":        strings.ToLower,
		"title":          strings.Title,
		// join is equal to strings.Join but inverts the argument order
		// for easier pipelining in templates.
		"join": func(sep string, s []string) string {
			return strings.Join(s, sep)
		},
		"match": regexp.MatchString,
		"safeHtml": func(text string) tmplhtml.HTML {
			return tmplhtml.HTML(text)
		},
		"reReplaceAll": func(pattern, repl, text string) string {
			re := regexp.MustCompile(pattern)
			return re.ReplaceAllString(text, repl)
		},
		"stringSlice": func(s ...string) []string {
			return s
		},
	}
	return DefaultFuncs
}

type Sender struct {
	Typ        string    `json:",omitempty"`
	Status     string    `json:",omitempty"`
	NotifyType string    `json:",omitempty"`
	Alerts     AlertList `json:",omitempty"`
}

type Senders []Sender

func RenderText(sender Sender, fileName string) (content string, err error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("err: %s\n", err)
			level.Error(logger).Log("type", "panic", "message", err)
		}
	}()
	//fmt.Printf("sender: %s\n", sender)

	if len(sender.Alerts) == 0 {
		return
	}

	// 读取配置文件conf中的Alerts.Templates中定义的模板文件内容
	temp, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	}
	// Create a template, add the function map, and parse the text.
	tmpl, err := tmpltext.New("*.tmpl").Funcs(DefaultFuncs()).Parse(string(temp))
	if err != nil {
		fmt.Printf("err: %s\n", err)
		level.Error(logger).Log("message", string(temp))
		level.Error(logger).Log("message", err.Error())
	}

	// Run the template to verify the output.
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, sender)
	if err != nil {
		//log.Fatalf("execution: %s", err.Error())
		level.Error(logger).Log("type", "render", "fileName", fileName, "message", err.Error())
	}
	content = buf.String()
	return content, nil
}
