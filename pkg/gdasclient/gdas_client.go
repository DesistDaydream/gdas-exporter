package gdasclient

import (
	"time"

	"github.com/DesistDaydream/gdas-exporter/pkg/gdassdk"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

// GdasClient 连接 Gdas 所需信息
type GdasClient struct {
	Services *gdassdk.Services
	Opts     *GdasOpts
}

// NewGdasClient 实例化 Gdas 客户端
func NewGdasClient(opts *GdasOpts) *GdasClient {
	token, err := gdassdk.GetToken(opts.URL, opts.Username, opts.Password)
	if err != nil {
		logrus.Errorf("获取 Token 失败，请检查目标设备是否正常")
		// panic(err)
	}

	services := gdassdk.NewServices(opts.URL, token, opts.Timeout)

	return &GdasClient{
		Opts:     opts,
		Services: services,
	}
}

// Ping 在 Scraper 接口的实现方法 scrape() 中调用。
// 让 Exporter 每次获取数据时，都检验一下目标设备通信是否正常
func (c *GdasClient) Ping() (b bool, err error) {
	_, err = c.Services.Auth.GetAuthorize()
	if err != nil {
		logrus.Errorf("抓取指标前检查失败，重新获取 Token")
		token, err := gdassdk.GetToken(c.Opts.URL, c.Opts.Username, c.Opts.Password)
		if err != nil {
			logrus.Errorf("重新获取 Token 失败，请检查目标设备是否正常")
			return false, err
		}
		c.Services.Client.Token = token
	}
	return true, nil
}

func (c *GdasClient) GetConcurrency() int {
	return c.Opts.Concurrency
}

// GdasOpts 登录 Gdas 所需属性
type GdasOpts struct {
	URL         string
	Username    string
	Password    string
	Concurrency int
	// 这俩是关于 http.Client 的选项
	Timeout  time.Duration
	Insecure bool
}

// AddFlag use after set Opts
func (o *GdasOpts) AddFlag() {
	pflag.StringVarP(&o.URL, "gdas-server", "s", "https://172.38.30.192:8003", "HTTP API address of a Gdas server or agent. (prefix with https:// to connect over HTTPS)")
	pflag.StringVarP(&o.Username, "gdas-user", "u", "system", "gdas username")
	pflag.StringVarP(&o.Password, "gdas-pass", "p", "", "gdas password")
	pflag.IntVar(&o.Concurrency, "concurrent", 10, "Number of concurrent requests during collection.")
	pflag.DurationVar(&o.Timeout, "timeout", time.Second*60, "Timeout on HTTP requests to the Gads API.")
	pflag.BoolVar(&o.Insecure, "insecure", true, "Disable TLS host verification.")
}
