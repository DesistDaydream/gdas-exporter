package collector

import (
	"fmt"

	"github.com/DesistDaydream/gdas-exporter/pkg/gdasclient"
	"github.com/DesistDaydream/gdas-exporter/pkg/scraper"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

var (
	_ scraper.Scraper = ScrapeUser{}

	users = prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, "", "user_count"),
		"集群中所有用户的总数",
		[]string{}, nil,
	)
)

// ScrapeUser 抓取用户信息
type ScrapeUser struct{}

// Name is
func (ScrapeUser) Name() string {
	return "gdas_user_count"
}

// Help is
func (ScrapeUser) Help() string {
	return "Gdas User Info"
}

// Scrape is
func (ScrapeUser) Scrape(client *gdasclient.GdasClient, ch chan<- prometheus.Metric) (err error) {
	data, err := client.Services.Users.GetUsers()
	if err != nil {
		logrus.Errorf("获取用户信息失败:%v", err)
		return err
	} else if data == nil {
		return fmt.Errorf("获取到的用户信息为空")
	}

	userCount := float64(data.ResCount)
	ch <- prometheus.MustNewConstMetric(users, prometheus.GaugeValue, userCount)
	return nil
}
