package services

import core "github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/core/v1"

type PoolsService struct {
	client *core.Client
}

func NewPoolsService(client *core.Client) *PoolsService {
	return &PoolsService{
		client: client,
	}
}

// GetPools 3.查询所有盘匣组信息
func (n *PoolsService) GetPools() (*core.Pools, error) {
	var data core.Pools
	endpoint := "pools"
	_, err := n.client.RequestObj(endpoint, &data, &core.RequestOptions{
		Method:   "GET",
		RawQuery: "pages=1-*&poolName=&poolType=2&poolFlag=false&operation=",
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}
