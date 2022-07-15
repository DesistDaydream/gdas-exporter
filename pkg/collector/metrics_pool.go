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
	_ scraper.Scraper = ScrapePools{}
	// 全局盘匣组总数
	PoolTotalCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "pool_total_count"),
		"全局盘匣组总数",
		[]string{}, nil,
	)
	// 盘匣组状态
	PoolStatus = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "pool_status"),
		"盘匣组状态.0-空闲,1-刻录",
		[]string{"pool_can_del_flag", "type", "default_mgz", "pool_name", "user", "pool_raidLvl", "auto_add_mgz"}, nil,
	)

	// 盘匣组存储总空间
	PoolTotalSpace = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "pool_total_space"),
		"盘匣组总存储空间,单位:Byte",
		[]string{"pool_name", "pool_raidLvl"}, nil,
	)
	// 盘匣组存储剩余空间
	PoolAvailableSpace = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "pool_available_space"),
		"盘匣组可用存储空间,单位:Byte",
		[]string{"pool_name", "pool_raidLvl"}, nil,
	)
	// 盘匣组中的盘匣数量
	PoolRfidCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "pool_rfid_count"),
		"盘匣组包含的盘匣数量",
		[]string{"pool_name", "pool_raidLvl"}, nil,
	)
)

// ScrapePools is
type ScrapePools struct{}

// Name is
func (ScrapePools) Name() string {
	return "gdas_pool_info"
}

// Help is
func (ScrapePools) Help() string {
	return "Gdas Pool Info"
}

// Scrape is
func (ScrapePools) Scrape(client *gdasclient.GdasClient, ch chan<- prometheus.Metric) (err error) {
	// var pooldata poolData

	// url := "/api/gdas/pool/list"
	// method := "POST"
	// reqBody := `{"poolFlag":false,"poolName":"","poolType":""}`
	// buf := bytes.NewBuffer([]byte(reqBody))
	// respBody, err := client.Request(method, url, buf)
	// if err != nil {
	// 	return err
	// }

	// err = json.Unmarshal(respBody, &pooldata)
	// if err != nil {
	// 	return err
	// }

	data, err := client.Services.Pools.GetPools()
	if err != nil {
		logrus.Errorf("获取盘匣组信息失败:%v", err)
		return err
	} else if data == nil {
		return fmt.Errorf("获取到的盘匣组信息为空")
	}

	// 全局盘匣组总数
	ch <- prometheus.MustNewConstMetric(PoolTotalCount, prometheus.GaugeValue, float64(data.ResCount))

	for i := 0; i < len(data.Pools); i++ {
		// 盘匣组状态
		ch <- prometheus.MustNewConstMetric(PoolStatus, prometheus.GaugeValue, float64(data.Pools[i].PoolSts),
			data.Pools[i].PoolCanDelFlag,
			data.Pools[i].Type,
			strconv.FormatBool(data.Pools[i].DefaultMgz),
			data.Pools[i].PoolName,
			data.Pools[i].User,
			strconv.Itoa(data.Pools[i].PoolRaidLvl),
			strconv.FormatBool(data.Pools[i].AutoAddMgz),
		)
		// 盘匣组存储总空间
		ch <- prometheus.MustNewConstMetric(PoolTotalSpace, prometheus.GaugeValue, float64(data.Pools[i].PoolTotalSpace),
			data.Pools[i].PoolName,
			strconv.Itoa(data.Pools[i].PoolRaidLvl),
		)
		// 盘匣组存储剩余空间
		ch <- prometheus.MustNewConstMetric(PoolAvailableSpace, prometheus.GaugeValue, float64(data.Pools[i].PoolAvailableSpace),
			data.Pools[i].PoolName,
			strconv.Itoa(data.Pools[i].PoolRaidLvl),
		)
		// 盘匣组中的盘匣数量
		ch <- prometheus.MustNewConstMetric(PoolRfidCount, prometheus.GaugeValue, float64(data.Pools[i].RfidCount),
			data.Pools[i].PoolName,
			strconv.Itoa(data.Pools[i].PoolRaidLvl),
		)
	}
	return nil
}
