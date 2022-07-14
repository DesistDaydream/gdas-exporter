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
	_ scraper.Scraper = ScrapeMagazines{}

	// 盘匣状态
	MagazinesStatus = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_status"),
		"盘匣状态.0-正常,1-异常",
		[]string{"dam_name", "ip", "da_name", "da_no", "rfid", "slot_no", "pool_name"}, nil,
	)
	// 盘匣是否已满
	MagazinesFull = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_full"),
		"盘匣空间是否已满.0-未满,1-已满",
		[]string{"dam_name", "ip", "da_name", "da_no", "rfid", "slot_no", "pool_name"}, nil,
	)
	// 盘匣是否已被分配
	MagazinesRfidSts = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_rfid_sts"),
		"盘匣是否已被分配.0-未分配,1-已分配",
		[]string{"dam_name", "ip", "da_name", "da_no", "rfid", "slot_no", "pool_name"}, nil,
	)
)

// ScrapeMagazines 是将要实现 Scraper 接口的一个 Metric 结构体
type ScrapeMagazines struct{}

// Name 指定自己定义的 抓取器 的名字，与 Metric 的名字不是一个概念，但是一般保持一致
// 该方法用于为 ScrapeMagazines 结构体实现 Scraper 接口
func (ScrapeMagazines) Name() string {
	return "gdas_magazines_info"
}

// Help 指定自己定义的 抓取器 的帮助信息，这里的 Help 的内容将会作为命令行标志的帮助信息。与 Metric 的 Help 不是一个概念。
// 该方法用于为 ScrapeMagazines 结构体实现 Scraper 接口
func (ScrapeMagazines) Help() string {
	return "Gdas Magazines Info"
}

// Scrape 从客户端采集数据，并将其作为 Metric 通过 channel(通道) 发送。主要就是采集 Gdas 集群信息的具体行为。
// 该方法用于为 ScrapeMagazines 结构体实现 Scraper 接口
func (ScrapeMagazines) Scrape(client *gdasclient.GdasClient, ch chan<- prometheus.Metric) (err error) {
	data, err := client.Services.Magazines.GetMagazines()
	if err != nil {
		logrus.Errorf("获取盘匣列表失败:%v", err)
		return err
	} else if data == nil {
		return fmt.Errorf("获取到的盘匣列表为空")
	}

	for i := 0; i < len(data.Rfid); i++ {
		// 盘匣状态
		ch <- prometheus.MustNewConstMetric(MagazinesStatus, prometheus.GaugeValue, float64(data.Rfid[i].Status),
			data.Rfid[i].DamName,
			data.Rfid[i].ServerIP,
			data.Rfid[i].DaName,
			strconv.Itoa(data.Rfid[i].DaNo),
			data.Rfid[i].Rfid,
			strconv.Itoa(data.Rfid[i].SlotNo),
			data.Rfid[i].PoolName,
		)
		// 盘匣空间是否已满
		ch <- prometheus.MustNewConstMetric(MagazinesFull, prometheus.GaugeValue, float64(data.Rfid[i].Full),
			data.Rfid[i].DamName,
			data.Rfid[i].ServerIP,
			data.Rfid[i].DaName,
			strconv.Itoa(data.Rfid[i].DaNo),
			data.Rfid[i].Rfid,
			strconv.Itoa(data.Rfid[i].SlotNo),
			data.Rfid[i].PoolName,
		)
		// 盘匣是否已分配
		ch <- prometheus.MustNewConstMetric(MagazinesRfidSts, prometheus.GaugeValue, float64(data.Rfid[i].RfidSts),
			data.Rfid[i].DamName,
			data.Rfid[i].ServerIP,
			data.Rfid[i].DaName,
			strconv.Itoa(data.Rfid[i].DaNo),
			data.Rfid[i].Rfid,
			strconv.Itoa(data.Rfid[i].SlotNo),
			data.Rfid[i].PoolName,
		)
	}
	return nil
}
