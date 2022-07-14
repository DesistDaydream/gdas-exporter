package services

import core "github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/core/v1"

type NodeService struct {
	client *core.Client
}

func NewNodeService(client *core.Client) *NodeService {
	return &NodeService{
		client: client,
	}
}
func (n *NodeService) GetNode() (*core.NodeList, error) {
	var nodeList core.NodeList
	endpoint := "nodeList"
	_, err := n.client.RequestObj(endpoint, &nodeList, &core.RequestOptions{
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	return &nodeList, nil
}
