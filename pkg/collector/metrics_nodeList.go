package collector

import (
	"fmt"

	"github.com/DesistDaydream/gdas-exporter/pkg/gdasclient"
	"github.com/DesistDaydream/gdas-exporter/pkg/scraper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	// check interface
	_ scraper.Scraper = ScrapeNodeList{}

	// 设置 Metric 的基本信息
	nodelist = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "node_status"),
		"节点状态:0-活跃,1-异常",
		[]string{"dam_name", "ip"}, nil,
	)
)

// ScrapeNodeList 是将要实现 Scraper 接口的一个 Metric 结构体
type ScrapeNodeList struct{}

// Name 指定自己定义的 抓取器 的名字，与 Metric 的名字不是一个概念，但是一般保持一致
// 该方法用于为 ScrapeNodeList 结构体实现 Scraper 接口
func (ScrapeNodeList) Name() string {
	return "gdas_node_info"
}

// Help 指定自己定义的 抓取器 的帮助信息，这里的 Help 的内容将会作为命令行标志的帮助信息。与 Metric 的 Help 不是一个概念。
// 该方法用于为 ScrapeNodeList 结构体实现 Scraper 接口
func (ScrapeNodeList) Help() string {
	return "Gdas Node Info"
}

// Scrape 从客户端采集数据，并将其作为 Metric 通过 channel(通道) 发送。主要就是采集 Gdas 集群信息的具体行为。
// 该方法用于为 ScrapeNodeList 结构体实现 Scraper 接口
func (ScrapeNodeList) Scrape(client *gdasclient.GdasClient, ch chan<- prometheus.Metric) (err error) {
	data, err := client.Services.Node.GetNode()
	if err != nil {
		logrus.Errorf("获取节点列表失败:%v", err)
		return err
	} else if data == nil {
		return fmt.Errorf("获取到的节点列表为空")
	}

	logrus.Debugf("当前共有 %v 个节点", len(data.NodeList))

	for i := 0; i < len(data.NodeList); i++ {
		// 节点状态
		ch <- prometheus.MustNewConstMetric(nodelist, prometheus.GaugeValue, data.NodeList[i].Status,
			data.NodeList[i].IP,
			data.NodeList[i].DamName,
		)
	}
	return nil
}
