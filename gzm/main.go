package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func helloHandler(c *gin.Context) {
	totalCostTime.Observe([]string{c.Request.Method, c.Request.RequestURI}, 0.1)
	totalRequestQps.Inc([]string{c.Request.Method, c.Request.RequestURI})
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello q1mi!",
	})
}

var totalCostTime *ginmetrics.Metric
var totalRequestQps *ginmetrics.Metric

func GinMetricsAddMetric(m *ginmetrics.Monitor) {
	totalCostTime = &ginmetrics.Metric{
		Type:        ginmetrics.Histogram,
		Name:        "total_cost_time",
		Description: "请求处理耗时",
		Labels:      []string{"method", "url"},
		Buckets:     []float64{0.001, 0.005, 0.01, 0.02, 0.04, 0.08, 0.16, 0.3, 1.2, 5, 10},
	}
	m.AddMetric(totalCostTime)
	totalRequestQps = &ginmetrics.Metric{
		Type:        ginmetrics.Counter,
		Name:        "total_request_qps",
		Description: "请求次数",
		Labels:      []string{"method", "url"},
	}
	m.AddMetric(totalRequestQps)
}

func GinMetricsRegist(r *gin.Engine) {
	m := ginmetrics.GetMonitor()
	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(30)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	// 添加指标
	GinMetricsAddMetric(m)
	// set middleware for gin
	m.Use(r)
}

func main() {
	r := gin.Default()
	GinMetricsRegist(r)
	r.GET("/hello", helloHandler)
	if err := r.Run(); err != nil {
		fmt.Printf("startup service failed, err:%v\n", err)
	}
}
