package services

import core "github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/core/v1"

type AuthService struct {
	client *core.Client
}

func NewAuthService(client *core.Client) *AuthService {
	return &AuthService{
		client: client,
	}
}

// GetAuthorize 36.查询授权
func (n *AuthService) GetAuthorize() (*core.Authorize, error) {
	var data core.Authorize
	endpoint := "authorize"
	_, err := n.client.RequestObj(endpoint, &data, &core.RequestOptions{
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}
