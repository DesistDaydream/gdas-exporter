package collector

import (
	"fmt"
	"strconv"

	"github.com/DesistDaydream/gdas-exporter/pkg/gdasclient"
	"github.com/DesistDaydream/gdas-exporter/pkg/scraper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	// check interface
	_ scraper.Scraper = ScrapeDas{}
	// 全局盘库总数
	DasTotalCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "gdas_das_total_count"),
		"全局盘库总数",
		[]string{}, nil,
	)
	// 盘库状态
	DasStatus = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "gdas_das_status"),
		"盘库状态.0-连接正常,-203-盘匣弹出中,-210-仓架解锁中,-202-系统繁忙,-102-断开连接,-100&&-103-识别中",
		[]string{"dam_name", "ip", "da_name", "da_no", "da_vendor"}, nil,
	)
	// 盘库注册、断开状态
	DasOffline = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "gdas_das_offline"),
		"盘库注册、断开状态.0-已断开,1-已注册",
		[]string{"dam_name", "ip", "da_name", "da_no", "da_vendor"}, nil,
	)
	// 盘库所具有的槽位总数
	DasSlotCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "gdas_das_slot_count"),
		"盘库所具有的槽位总数，槽位就是可以插盘匣的位置",
		[]string{"dam_name", "ip", "da_name", "da_no", "da_vendor"}, nil,
	)
)

// ScrapeDas 是将要实现 Scraper 接口的一个 Metric 结构体
type ScrapeDas struct{}

// Name 指定自己定义的 抓取器 的名字，与 Metric 的名字不是一个概念，但是一般保持一致
// 该方法用于为 ScrapeDas 结构体实现 Scraper 接口
func (ScrapeDas) Name() string {
	return "gdas_das_info"
}

// Help 指定自己定义的 抓取器 的帮助信息，这里的 Help 的内容将会作为命令行标志的帮助信息。与 Metric 的 Help 不是一个概念。
// 该方法用于为 ScrapeDas 结构体实现 Scraper 接口
func (ScrapeDas) Help() string {
	return "Gdas Das Info"
}

// Scrape 从客户端采集数据，并将其作为 Metric 通过 channel(通道) 发送。主要就是采集 Gdas 集群信息的具体行为。
// 该方法用于为 ScrapeDas 结构体实现 Scraper 接口
func (ScrapeDas) Scrape(client *gdasclient.GdasClient, ch chan<- prometheus.Metric) (err error) {
	data, err := client.Services.Das.GetDas()
	if err != nil {
		logrus.Errorf("获取全局盘库信息失败:%v", err)
		return err
	} else if data == nil {
		return fmt.Errorf("获取到的全局盘库信息为空")
	}

	// 全局盘库总数
	ch <- prometheus.MustNewConstMetric(DasTotalCount, prometheus.GaugeValue, float64(len(data.DaInfo)))

	for i := 0; i < len(data.DaInfo); i++ {
		//盘库状态
		ch <- prometheus.MustNewConstMetric(DasStatus, prometheus.GaugeValue, float64(data.DaInfo[i].DaStatus),
			data.DaInfo[i].DamName,
			data.DaInfo[i].IP,
			data.DaInfo[i].DaName,
			strconv.Itoa(data.DaInfo[i].DaNo),
			data.DaInfo[i].DaVendor,
		)
		// 盘库注册、断开状态
		ch <- prometheus.MustNewConstMetric(DasOffline, prometheus.GaugeValue, float64(data.DaInfo[i].Offline),
			data.DaInfo[i].DamName,
			data.DaInfo[i].IP,
			data.DaInfo[i].DaName,
			strconv.Itoa(data.DaInfo[i].DaNo),
			data.DaInfo[i].DaVendor,
		)
		// 盘库所具有的槽位总数
		ch <- prometheus.MustNewConstMetric(DasSlotCount, prometheus.GaugeValue, float64(data.DaInfo[i].SlotCount),
			data.DaInfo[i].DamName,
			data.DaInfo[i].IP,
			data.DaInfo[i].DaName,
			strconv.Itoa(data.DaInfo[i].DaNo),
			data.DaInfo[i].DaVendor,
		)
	}

	return nil
}
