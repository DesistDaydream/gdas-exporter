package services

import core "github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/core/v1"

type MagazinesService struct {
	client *core.Client
}

func NewMagazinesService(client *core.Client) *MagazinesService {
	return &MagazinesService{
		client: client,
	}
}

// GetMagazines 11.查询盘匣列表
func (n *MagazinesService) GetMagazines() (*core.Magazines, error) {
	var data core.Magazines
	endpoint := "magazines"
	_, err := n.client.RequestObj(endpoint, &data, &core.RequestOptions{
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetTotalspace 28.获取全局存储空间信息
func (n *MagazinesService) GetTotalspace() (*core.Totalspace, error) {
	var data core.Totalspace
	endpoint := "totalspace"
	_, err := n.client.RequestObj(endpoint, &data, &core.RequestOptions{
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}
