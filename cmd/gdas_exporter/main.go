package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/DesistDaydream/gdas-exporter/pkg/collector"
	"github.com/DesistDaydream/gdas-exporter/pkg/gdasclient"
	"github.com/DesistDaydream/gdas-exporter/pkg/logging"
	"github.com/DesistDaydream/gdas-exporter/pkg/scraper"
	"github.com/coreos/go-systemd/daemon"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// scrapers 列出了应该注册的所有 Scraper(抓取器)，以及默认情况下是否应该启用它们
// 用一个 map 来定义这些抓取器是否开启，key 为 collector.Scraper 接口类型，value 为 bool 类型。
// 凡是实现了 collector.Scraper 接口的结构体，都可以做作为该接口类型的值
var scrapers = map[scraper.Scraper]bool{
	collector.ScrapeMagazines{}:  true,
	collector.ScrapeNodeList{}:   true,
	collector.ScrapeUser{}:       true,
	collector.ScrapeTotalspace{}: true,
	collector.ScrapePools{}:      true,
	collector.ScrapeDas{}:        true,
}

// DumpStacks is
func DumpStacks() {
	buf := make([]byte, 16384)
	buf = buf[:runtime.Stack(buf, true)]
	logrus.Printf("=== BEGIN goroutine stack dump ===\n%s\n=== END goroutine stack dump ===", buf)
}

func main() {
	// ####################################
	// ######## 设置命令行标志，开始 ########
	// ####################################
	listenAddress := pflag.StringP("listen-address", "l", ":8003", "Address to listen on for web interface and telemetry.")
	metricsPath := pflag.String("telemetry-path", "/metrics", "Path under which to expose metrics.")

	// 设置日志相关命令行标志
	logFlags := logging.LoggingFlags{}
	logFlags.AddFlags()

	// 设置关于抓取 Metric 目标客户端的一些信息的标志
	opts := &gdasclient.GdasOpts{}
	opts.AddFlag()

	// scraperFlags 也是一个 map，并且 key 为 collector.Scraper 接口类型，这一小段代码主要有下面几个作用
	// 1.生成抓取器的命令行标志，用于通过命令行控制开启哪些抓取器，说白了就是控制采集哪些指标
	// 2.下面的 for 循环会通过命令行 flag 获取到的值，放到 scraperFlags 这个 map 中
	// 3.然后在后面注册 Exporter 之前，先通过这个 map 中的键值对判断是否要把 value 为 true 的 抓取器 注册进去
	scraperFlags := map[scraper.Scraper]*bool{}
	for scraper, enabledByDefault := range scrapers {
		defaultOn := false
		if enabledByDefault {
			defaultOn = true
		}
		// 设置命令行 flag
		f := pflag.Bool("collect."+scraper.Name(), defaultOn, scraper.Help())
		// 将命令行 flag 中获取到的值，赋到 map 中，作为 map 的 value
		scraperFlags[scraper] = f
	}
	// 解析命令行标志,即：将命令行标志的值传递到代码的变量中。若不解析，则所有通过命令行标志设置的变量是没有值的。
	pflag.Parse()
	// ####################################
	// ######## 设置命令行标志，结束 ########
	// ####################################

	// 初始化日志
	if err := logging.LogInit(logFlags.LogLevel, logFlags.LogOutput, logFlags.LogFormat); err != nil {
		logrus.Fatal("初始化日志失败", err)
	}

	// 下面的都是 Exporter 运行的最主要逻辑了
	//
	// 获取所有通过命令行标志，设置开启的 scrapers(抓取器)。
	// 不包含默认开启的，默认开启的在代码中已经指定了。
	enabledScrapers := []scraper.Scraper{}
	for scraper, enabled := range scraperFlags {
		if *enabled {
			logrus.Info("Scraper enabled ", scraper.Name())
			enabledScrapers = append(enabledScrapers, scraper)
		}
	}
	// 实例化 Exporter，其中包括所有自定义的 Metrics
	GdasClient := gdasclient.NewGdasClient(opts)
	exporter := scraper.NewExporter(GdasClient, enabledScrapers)
	// 实例化一个注册器,并使用这个注册器注册 exporter
	reg := prometheus.NewRegistry()
	reg.MustRegister(exporter)

	// 设置路由信息
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>` + collector.Name() + `</title></head>
             <body>
             <h1>` + collector.Name() + `</h1>
             <p><a href='` + *metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	http.Handle(*metricsPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{ErrorLog: logrus.StandardLogger()}))

	http.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok")
	})

	// 启动前检查并启动 Exporter
	logrus.Info("Listening on address ", *listenAddress)
	daemon.SdNotify(false, daemon.SdNotifyReady)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		logrus.Fatal(err)
	}
}
