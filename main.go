package main

import (
	"AlertMsgRouter/controllers/apis"
	collector "AlertMsgRouter/controllers/metrics"
	"flag"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Set during go build
	// version   string
	// gitCommit string

	// 命令行参数
	listenAddr       = flag.String("web.listen-port", "9000", "An port to listen on for web interface and telemetry.")
	metricsPath      = flag.String("web.telemetry-path", "/metrics", "A path under which to expose metrics.")
	metricsNamespace = flag.String("metric.namespace", "AlertMsgRouter", "Prometheus metrics namespace, as the prefix of metrics name")
)

func init() {
	//beego router
	ns := beego.NewNamespace("/api/v1",
		beego.NSRouter("/prometheusalert", &apis.PrometheusAlertController{}, "get,post:PrometheusAlert"),
	)
	beego.AddNamespace(ns)
	beego.SetStaticPath("/swagger", "swagger")

	//
	//beego.Router("/prometheusalert", &apis.PrometheusAlertController{}, "get,post:PrometheusAlert")
	//beego.Router("/cpu", &cpu.Controller{}, "get:Get")
	//beego.Router("/mem", &memory.Controller{}, "get:Get")
	//beego.Router("/mem", &collector.Controller{}, "get:Get")

	//metrics for self -- beego handler
	flag.Parse()
	metrics := collector.NewMetrics(*metricsNamespace)
	registry := prometheus.NewRegistry()
	registry.MustRegister(metrics)

	//registry.MustRegister(collector.AlertFailedCounter)
	//registry.MustRegister(metrics)
	//registry.MustRegister(metrics)
	beego.Handler(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	//metrics for self -- http handler
	//http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	//logs.Info("Main-metrics-%v\n ", registry)
	//
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write([]byte(`<html>
	//		<head><title>A Prometheus Exporter</title></head>
	//		<body>
	//		<h1>A Prometheus Exporter</h1>
	//		<p><a href='/metrics'>Metrics</a></p>
	//		</body>
	//		</html>`))
	//})
	//log.Info("Starting Server at http://localhost:%s%s", *listenAddr, *metricsPath)
	//log.Fatal(http.ListenAndServe(":"+*listenAddr, nil))

}

func main() {
	orm.Debug = true
	logtype := beego.AppConfig.String("logtype")
	if logtype == "console" {
		logs.SetLogger(logtype)
	} else if logtype == "file" {
		logpath := beego.AppConfig.String("logpath")
		logs.SetLogger(logtype, `{"filename":"`+logpath+`"}`)
	}
	logs.Info("[main] 当前版本（Version）4.3.4")
	//beego.Handler("/metrics",promhttp.Handler())
	//beego.Router("/cpu", &cpu.Controller{}, "get:Get")

	beego.Run()
}

//var (
//	// Set during go build
//	// version   string
//	// gitCommit string
//
//	// 命令行参数
//	listenAddr  = flag.String("web.listen-port", "9000", "An port to listen on for web interface and telemetry.")
//	metricsPath = flag.String("web.telemetry-path", "/metrics", "A path under which to expose metrics.")
//	metricsNamespace = flag.String("metric.namespace", "AlertMsgRouter", "Prometheus metrics namespace, as the prefix of metrics name")
//)

//func main() {
//	flag.Parse()
//	metrics := collector.NewMetrics(*metricsNamespace)
//	registry := prometheus.NewRegistry()
//	registry.MustRegister(metrics)
//
//	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
//	logs.Info("Main-metrics-%v\n ", metrics)
//
//	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte(`<html>
//			<head><title>A Prometheus Exporter</title></head>
//			<body>
//			<h1>A Prometheus Exporter</h1>
//			<p><a href='/metrics'>Metrics</a></p>
//			</body>
//			</html>`))
//	})
//	log.Info("Starting Server at http://localhost:%s%s", *listenAddr, *metricsPath)
//	log.Fatal(http.ListenAndServe(":"+*listenAddr, nil))
//}
