package collector

import (
	"AlertMsgRouter/controllers/metrics/collector"
	_ "fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"sync"
	"time"
)

// 指标结构体
type Metrics struct {
	metrics map[string]*prometheus.Desc
	mutex   sync.Mutex
}

var (
	AlertsFromCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "alers_from_count",
		Help: "count alers from any where",
	},
		[]string{"from", "message", "level", "host", "index"},
	)
	//model.AlertsFromCounter.WithLabelValues("from","to","message","level","host","index").Add(1)
	AlertToCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "alers_to_count",
		Help: "count alers to any where",
	},
		[]string{"to", "message", "phone"},
	)
	//model.AlertToCounter.WithLabelValues("to","message").Add(1)
	AlertFailedCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "alers_send_failed_count",
		Help: "count alers send failed",
	},
		[]string{"to", "message", "phone"},
	)
	//model.AlertFailedCounter.WithLabelValues("to","message","phone").Add(1)
)

//func MetricsInit() {
//	prometheus.MustRegister(AlertsFromCounter)
//	prometheus.MustRegister(AlertToCounter)
//	prometheus.MustRegister(AlertFailedCounter)
//}

/**
 * 函数：newGlobalMetric
 * 功能：创建指标描述符
 */
func newGlobalMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}

/**
 * 工厂方法：NewMetrics
 * 功能：初始化指标信息，即Metrics结构体
 */
func NewMetrics(namespace string) *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"my_counter_metric": newGlobalMetric(namespace, "server_A_Name", "The description of my_counter_metric", []string{"name"}),
			"my_ip_metric":      newGlobalMetric(namespace, "server_A_IP", "The description of my_counter_metric", []string{"ip"}),
			"my_gauge_metric":   newGlobalMetric(namespace, "server_Network", "The description of my_gauge_metric", []string{"name"}),
			"my_cpu_metric":     newGlobalMetric(namespace, "server_Cpu", "The description of my_cpu_metric", []string{"name"}),
			"my_disk_metric":    newGlobalMetric(namespace, "server_Disk", "The description of my_cpu_metric", []string{"name"}),
			"my_mem_metric":     newGlobalMetric(namespace, "server_Mem", "The description of my_cpu_metric", []string{"name"}),
		},
	}
}

/**
 * 接口：Describe
 * 功能：传递结构体中的指标描述符到channel
 */
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

/**
 * 接口：Collect
 * 功能：抓取最新的数据，传递给channel
 */
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // 加锁
	defer c.mutex.Unlock()

	mockCounterMetricData, mockGaugeMetricData, mockCpuMetricData, mockDiskMetricData, mockMemMetricData, mockIpMetricData := c.GenerateMockData()
	for host, currentValue := range mockCounterMetricData {
		ch <- prometheus.MustNewConstMetric(c.metrics["my_counter_metric"], prometheus.CounterValue, float64(currentValue), host)
	}
	for host, currentValue := range mockGaugeMetricData {
		ch <- prometheus.MustNewConstMetric(c.metrics["my_gauge_metric"], prometheus.GaugeValue, float64(currentValue), host)
	}
	for host, currentValue := range mockCpuMetricData {
		ch <- prometheus.MustNewConstMetric(c.metrics["my_cpu_metric"], prometheus.GaugeValue, float64(currentValue), host)
	}
	for host, currentValue := range mockDiskMetricData {
		ch <- prometheus.MustNewConstMetric(c.metrics["my_disk_metric"], prometheus.GaugeValue, float64(currentValue), host)
	}
	for host, currentValue := range mockMemMetricData {
		ch <- prometheus.MustNewConstMetric(c.metrics["my_mem_metric"], prometheus.GaugeValue, float64(currentValue), host)
	}
	for host, currentValue := range mockIpMetricData {
		ch <- prometheus.MustNewConstMetric(c.metrics["my_ip_metric"], prometheus.CounterValue, float64(currentValue), host)
	}
}

/**
 * 函数：GenerateMockData
 * 功能：生成模拟数据
 */
func (c *Metrics) GenerateMockData() (mockCounterMetricData map[string]float64, mockGaugeMetricData map[string]float64, mockCpuMetricData map[string]float64, mockDiskMetricData map[string]float64, mockMemMetricData map[string]float64, mockIpMetricData map[string]float64) {
	n, _ := host.Info()
	d, _ := disk.Usage("/")
	v, _ := mem.VirtualMemory()
	boottime, _ := host.BootTime()
	btime := time.Unix(int64(boottime), 0).Format("2006-01-02 15:04:05")
	cc, _ := cpu.Percent(time.Second, false)
	nv, _ := net.IOCounters(true)
	ip := collector.GetOutboundIP()
	diskTotal := float64(d.Total / 1024 / 1024 / 1024)
	diskFree := float64(d.Free / 1024 / 1024 / 1024)
	diskUsage := float64(d.UsedPercent)
	memTotal := float64(v.Total / 1024 / 1024)
	memFree := float64(v.Available / 1024 / 1024)
	memUsed := float64(v.Used / 1024 / 1024)
	memUsage := float64(v.UsedPercent)
	cpuUsage := float64(cc[0])
	bytesRecv := float64(nv[0].BytesRecv)
	bytesSent := float64(nv[0].BytesSent)

	running, zombie, sleep := collector.GetProcessesStatus()

	mockCounterMetricData = map[string]float64{
		n.Hostname: 0,
		btime:      0,
	}

	mockIpMetricData = map[string]float64{

		ip: 0,
	}
	mockCpuMetricData = map[string]float64{
		"Cpu Usage(%)":      cpuUsage,
		"Processes running": running,
		"Processes zombie":  zombie,
		"Processes sleep":   sleep,
	}
	mockDiskMetricData = map[string]float64{
		"Disk Total(GB)": diskTotal,
		"Disk Free(GB)":  diskFree,
		"Disk Usage(%)":  diskUsage,
	}
	mockMemMetricData = map[string]float64{
		"Mem Total(MB)": memTotal,
		"Mem Free(MB)":  memFree,
		"Mem Used(MB)":  memUsed,
		"Mem Usage(%)":  memUsage,
	}

	mockGaugeMetricData = map[string]float64{
		//	"disk Total(GB)":diskTotal,
		//	"disk Free(GB)":diskFree,
		//	"disk Usage(%)":diskUsage,

		//	"mem Total(MB)":memTotal,
		//	"mem Free(MB)":memFree,
		//	"mem Used(MB)":memUsed,
		//	"mem Usage(%)":memUsage,

		//	"cpu Usage(%)":cpuUsage,
		"Bytes Recv(bytes)": bytesRecv,
		"Bytes Sent(bytes)": bytesSent,
	}
	return
}
