package services

import core "github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/core/v1"

type DasService struct {
	client *core.Client
}

func NewDasService(client *core.Client) *DasService {
	return &DasService{
		client: client,
	}
}

// GetDas 9.获取全局盘库信息
func (n *DasService) GetDas() (*core.Das, error) {
	var data core.Das
	endpoint := "das"
	_, err := n.client.RequestObj(endpoint, &data, &core.RequestOptions{
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}
