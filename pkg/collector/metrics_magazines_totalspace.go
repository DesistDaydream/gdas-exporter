package collector

import (
	"fmt"

	"github.com/DesistDaydream/gdas-exporter/pkg/gdasclient"
	"github.com/DesistDaydream/gdas-exporter/pkg/scraper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	_ scraper.Scraper = ScrapeTotalspace{}
	// 全局盘匣raid0总空间,即裸容量
	MagazinesTotalSpaceRaid0 = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_total_space_raid0"),
		"全部盘匣以raid0级别计算总空间，即裸容量，单位:bytes",
		[]string{}, nil,
	)
	// 全局盘匣实际总空间
	MagazinesTotalSpace = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_total_space"),
		"全部盘匣实际总空间,即做了raid和没做raid的所有盘匣的总容量,单位:Byte",
		[]string{}, nil,
	)
	// 全局盘匣实际剩余总空间
	MagazinesTotalAvailableSpace = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_total_available_space"),
		"全部盘匣实际可用空间,单位:Byte",
		[]string{}, nil,
	)
	// 全局盘匣实际已用总空间
	MagazinesTotalUsedSpace = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_total_used_space"),
		"全部盘匣实际使用空间,单位:Byte",
		[]string{}, nil,
	)
	// 全部盘匣槽位数
	MagzainesTotalSlotCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_total_slot_count"),
		"全部盘匣槽位数",
		[]string{}, nil,
	)
	// 全局盘匣总数
	MagazinesTotalCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_total_count"),
		"全部盘匣总数",
		[]string{}, nil,
	)
	// 全局盘匣已使用数
	MagazinesUsedCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_used_count"),
		"全部盘匣中,已使用的总数",
		[]string{}, nil,
	)
	// 全部盘匣中，未使用的总数
	MagazinesFreeCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_free_count"),
		"全部盘匣中,未使用的总数",
		[]string{}, nil,
	)
	// 全局盘匣异常数
	MagazinesExceptionCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "magazines_exception_count"),
		"全部盘匣中,异常状态的总数",
		[]string{}, nil,
	)
)

// ScrapeTotalspace is
type ScrapeTotalspace struct{}

// Name is
func (ScrapeTotalspace) Name() string {
	return "gdas_magazines_info"
}

// Help is
func (ScrapeTotalspace) Help() string {
	return "Gdas Magazines Info"
}

// Scrape is
func (ScrapeTotalspace) Scrape(client *gdasclient.GdasClient, ch chan<- prometheus.Metric) (err error) {
	data, err := client.Services.Magazines.GetTotalspace()
	if err != nil {
		logrus.Errorf("获取全局存储空间信息失败:%v", err)
		return err
	} else if data == nil {
		return fmt.Errorf("获取到的全局存储空间信息为空")
	}

	// 全局盘匣raid0总空间,即裸容量
	ch <- prometheus.MustNewConstMetric(MagazinesTotalSpaceRaid0, prometheus.GaugeValue, float64(data.TotalSpaceRaid0))
	// 全局盘匣实际总空间
	ch <- prometheus.MustNewConstMetric(MagazinesTotalSpace, prometheus.GaugeValue, float64(data.TotalSpace))
	// 全局盘匣实际剩余总空间
	ch <- prometheus.MustNewConstMetric(MagazinesTotalAvailableSpace, prometheus.GaugeValue, float64(data.TotalAvailableSpace))
	// 全局盘匣实际已用总空间
	ch <- prometheus.MustNewConstMetric(MagazinesTotalUsedSpace, prometheus.GaugeValue, float64(data.TotalSpace)-float64(data.TotalAvailableSpace))
	// 全部盘匣槽位数
	ch <- prometheus.MustNewConstMetric(MagzainesTotalSlotCount, prometheus.GaugeValue, float64(data.TotalSlotCount))
	// 全局盘匣总数
	ch <- prometheus.MustNewConstMetric(MagazinesTotalCount, prometheus.GaugeValue, float64(data.TotalMgzCount))
	// 全局盘匣已使用数
	ch <- prometheus.MustNewConstMetric(MagazinesUsedCount, prometheus.GaugeValue, float64(data.UsedMgzCount))
	// 全部盘匣未使用的总数
	ch <- prometheus.MustNewConstMetric(MagazinesFreeCount, prometheus.GaugeValue, float64(data.FreeMgzCount))
	// 全局盘匣异常数
	ch <- prometheus.MustNewConstMetric(MagazinesExceptionCount, prometheus.GaugeValue, float64(data.ExceptionMgzCount))

	return nil
}
