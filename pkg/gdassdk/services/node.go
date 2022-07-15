package services

import (
	"fmt"

	core "github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/core/v1"
)

type NodeService struct {
	client *core.Client
}

func NewNodeService(client *core.Client) *NodeService {
	return &NodeService{
		client: client,
	}
}

// GetNode 26.分布式节点查询
func (n *NodeService) GetNode() (*core.NodeList, error) {
	var data core.NodeList
	endpoint := "nodeList"
	_, err := n.client.RequestObj(endpoint, &data, &core.RequestOptions{
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetNodeCaches 13.缓存容量显示
func (n *NodeService) GetNodeCaches(ip string) (*core.NodeCaches, error) {
	var data core.NodeCaches
	endpoint := fmt.Sprintf("%s/caches", ip)
	_, err := n.client.RequestObj(endpoint, &data, &core.RequestOptions{
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetNodeDas 8.获取节点盘库列表
func (n *NodeService) GetNodeDas(ip string) (*core.NodeDas, error) {
	var data core.NodeDas
	endpoint := fmt.Sprintf("%s/das", ip)
	_, err := n.client.RequestObj(endpoint, &data, &core.RequestOptions{
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}
