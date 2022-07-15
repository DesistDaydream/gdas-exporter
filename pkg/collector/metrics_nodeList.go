package collector

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/DesistDaydream/gdas-exporter/pkg/gdasclient"
	"github.com/DesistDaydream/gdas-exporter/pkg/scraper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	// check interface
	_ scraper.Scraper = ScrapeNodeList{}

	// 全局节点总数
	NodeTotalCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "node_total_count"),
		"全局节点总数",
		[]string{}, nil,
	)
	// 节点状态
	NodeStatus = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "node_status"),
		"节点状态:0-活跃,1-异常",
		[]string{"dam_name", "ip"}, nil,
	)
	// 节点总缓存容量
	NodeTotalCacheSize = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "node_total_cache_size"),
		"节点总缓存容量,单位:Byte",
		[]string{"dam_name", "ip"}, nil,
	)
	// 节点已用缓存容量
	NodeUsedCacheSize = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "node_used_cache_size"),
		"节点已用缓存容量,单位:Byte",
		[]string{"dam_name", "ip"}, nil,
	)
	// 节点未用缓存容量
	NodeUnusedCacheSize = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "node_unused_cache_size"),
		"节点未用缓存容量,单位:Byte",
		[]string{"dam_name", "ip"}, nil,
	)
	// 盘库中盘库中已用盘匣数量
	magazineUsedCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "das_magazine_used_count"),
		"盘库中已用盘匣数量",
		[]string{"dam_name", "ip", "da_name", "da_no"}, nil,
	)

	// 盘库中未用盘匣数量
	magazineFreeCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "das_magazine_free_count"),
		"盘库中未用盘匣数量",
		[]string{"dam_name", "ip", "da_name", "da_no"}, nil,
	)
	// 盘库中异常盘匣数量
	magazineExcpCount = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "das_magazine_excp_count"),
		"盘库中异常盘匣数量",
		[]string{"dam_name", "ip", "da_name", "da_no"}, nil,
	)
	// 盘库中机械手状态
	changerStatus = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "das_changer_status"),
		"盘库中机械手的状态:0-寿命良好,1-寿命警告,2-寿命已到",
		[]string{"dam_name", "ip", "da_name", "da_no", "changer_serial"}, nil,
	)
	// 盘库中机械手状态
	changerUsedPercent = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "das_changer_used_percent"),
		"盘库中机械手使用百分比，该值需要除以 100",
		[]string{"dam_name", "ip", "da_name", "da_no", "changer_serial"}, nil,
	)
	// TODO:盘库中机械手数量

	// 光驱状态
	driveStatus = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "das_drive_status"),
		"盘库中光驱的状态:0-寿命良好,1-寿命警告,2-寿命已到",
		[]string{"dam_name", "ip", "da_name", "da_no", "drive_serial"}, nil,
	)
	// 光驱状态
	driveUsedPercent = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "das_drive_used_percent"),
		"盘库中光驱使用百分比，该值需要除以 100",
		[]string{"dam_name", "ip", "da_name", "da_no", "drive_serial"}, nil,
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
	var nodeIPList []string

	// #############################################
	// ####### 获取分布式节点信息,即node概况 #######
	// #############################################
	nodeData, err := client.Services.Node.GetNode()
	if err != nil {
		logrus.Errorf("获取节点列表失败:%v", err)
		return err
	} else if nodeData == nil {
		return fmt.Errorf("获取到的节点列表为空")
	}

	logrus.Debugf("当前共有 %v 个节点", len(nodeData.NodeList))

	// 获取 nodeIP 列表，放到切片中
	for i := 0; i < len(nodeData.NodeList); i++ {
		nodeIPList = append(nodeIPList, nodeData.NodeList[i].IP)
	}
	logrus.Debugf("所有节点 IP 列表：%v", nodeIPList)

	// 集群中节点总数
	ch <- prometheus.MustNewConstMetric(NodeTotalCount, prometheus.GaugeValue, float64(len(nodeData.NodeList)))

	var wg sync.WaitGroup
	defer wg.Wait()

	// 用来控制并发数量
	concurrencyControl := make(chan bool, client.GetConcurrency())

	for index, nodeIP := range nodeIPList {
		concurrencyControl <- true
		wg.Add(1)

		go func(index int, nodeIP string) error {
			defer wg.Done()
			// 获取节点的状态
			ch <- prometheus.MustNewConstMetric(NodeStatus, prometheus.GaugeValue, float64(nodeData.NodeList[index].Status),
				nodeData.NodeList[index].DamName,
				nodeIP,
			)

			// ################################################
			// ####### 循环每个节点，逐一获取节点的缓存数据 ######
			// ################################################
			cacheData, err := client.Services.Node.GetNodeCaches(nodeIP)
			if err != nil {
				logrus.Errorf("获取【%v】节点缓存失败:%v", nodeIP, err)
				return err
			} else if nodeData == nil {
				return fmt.Errorf("获取到的【%v】节点缓存为空", nodeIP)
			}

			//节点总缓存容量
			ch <- prometheus.MustNewConstMetric(NodeTotalCacheSize, prometheus.GaugeValue, float64(cacheData.TotalCacheSize),
				nodeData.NodeList[index].DamName,
				nodeIP,
			)
			//节点已用缓存容量
			ch <- prometheus.MustNewConstMetric(NodeUsedCacheSize, prometheus.GaugeValue, float64(cacheData.UsedCacheSize),
				nodeData.NodeList[index].DamName,
				nodeIP,
			)
			//节点未用缓存容量
			ch <- prometheus.MustNewConstMetric(NodeUnusedCacheSize, prometheus.GaugeValue, float64(cacheData.UnUsedCacheSize),
				nodeData.NodeList[index].DamName,
				nodeIP,
			)

			// ################################################
			// #### 循环每个节点，逐一获取节点下每个盘库的信息 ####
			// ################################################
			nodeDasData, err := client.Services.Node.GetNodeDas(nodeIP)
			if err != nil {
				logrus.Errorf("获取节点盘库失败:%v", err)
				return err
			} else if nodeData == nil {
				return fmt.Errorf("获取到的节点盘库为空")
			}
			// 每个节点下有多个盘库，所以循环每个盘库以获取指标
			for j := 0; j < len(nodeDasData.DaList); j++ {
				// 盘库中已用盘匣数量
				ch <- prometheus.MustNewConstMetric(magazineUsedCount, prometheus.GaugeValue, float64(nodeDasData.DaList[j].MagazineUsedCount),
					nodeData.NodeList[index].DamName,
					nodeIP,
					nodeDasData.DaList[j].Name,
					strconv.Itoa(nodeDasData.DaList[j].DaNo),
				)
				// 盘库中未用盘匣数量
				ch <- prometheus.MustNewConstMetric(magazineFreeCount, prometheus.GaugeValue, float64(nodeDasData.DaList[j].MagazineFreeCount),
					nodeData.NodeList[index].DamName,
					nodeIP,
					nodeDasData.DaList[j].Name,
					strconv.Itoa(nodeDasData.DaList[j].DaNo),
				)
				// 盘库中异常盘匣数量
				ch <- prometheus.MustNewConstMetric(magazineExcpCount, prometheus.GaugeValue, float64(nodeDasData.DaList[j].MagazineExcpCount),
					nodeData.NodeList[index].DamName,
					nodeIP,
					nodeDasData.DaList[j].Name,
					strconv.Itoa(nodeDasData.DaList[j].DaNo),
				)

				// 循环盘库下每个机械手，以获取指标
				for k := 0; k < len(nodeDasData.DaList[j].ChangerSmartInfo); k++ {
					// 机械手状态
					ch <- prometheus.MustNewConstMetric(changerStatus, prometheus.GaugeValue, float64(nodeDasData.DaList[j].ChangerSmartInfo[k].Status),
						nodeData.NodeList[index].DamName,
						nodeIP,
						nodeDasData.DaList[j].Name,
						strconv.Itoa(nodeDasData.DaList[j].DaNo),
						nodeDasData.DaList[j].ChangerSerial,
						// strconv.Itoa(status.DaList[j].ChangerSmartInfo[k].UnitNo),
					)
					// 机械手使用百分比
					ch <- prometheus.MustNewConstMetric(changerUsedPercent, prometheus.GaugeValue, float64(nodeDasData.DaList[j].ChangerSmartInfo[k].UsedPercent),
						nodeData.NodeList[index].DamName,
						nodeIP,
						nodeDasData.DaList[j].Name,
						strconv.Itoa(nodeDasData.DaList[j].DaNo),
						nodeDasData.DaList[j].ChangerSerial,
						// strconv.Itoa(status.DaList[j].ChangerSmartInfo[k].UnitNo),
					)
				}
				// 循环盘库下每个光驱，以获取指标
				// 判断一下与光驱有关的另一个数组中的元素是否不为空
				if len(nodeDasData.DaList[j].DriveSerialList) > 0 {
					for l := 0; l < len(nodeDasData.DaList[j].DriveSmartInfo); l++ {
						// 光驱状态
						ch <- prometheus.MustNewConstMetric(driveStatus, prometheus.GaugeValue, float64(nodeDasData.DaList[j].DriveSmartInfo[l].Status),
							nodeData.NodeList[index].DamName,
							nodeIP,
							nodeDasData.DaList[j].Name,
							strconv.Itoa(nodeDasData.DaList[j].DaNo),
							nodeDasData.DaList[j].DriveSerialList[l].DriveSerial,
							// strconv.Itoa(status.DaList[j].DriveSmartInfo[l].UnitNo),
						)
						// 光驱使用百分比
						ch <- prometheus.MustNewConstMetric(driveUsedPercent, prometheus.GaugeValue, float64(nodeDasData.DaList[j].DriveSmartInfo[l].UsedPercent),
							nodeData.NodeList[index].DamName,
							nodeIP,
							nodeDasData.DaList[j].Name,
							strconv.Itoa(nodeDasData.DaList[j].DaNo),
							nodeDasData.DaList[j].DriveSerialList[l].DriveSerial,
							// strconv.Itoa(status.DaList[j].DriveSmartInfo[l].UnitNo),
						)
					}
				} else {
					logrus.Error("从 API 获取光驱指标异常,DriveSerialList 数组元素不大于0")
				}
			}
			<-concurrencyControl
			return nil
		}(index, nodeIP)
	}
	return nil
}
