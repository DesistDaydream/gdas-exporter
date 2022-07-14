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
func (n *AuthService) GetAuthorize() (*core.Authorize, error) {
	var authorize core.Authorize
	endpoint := "authorize"
	_, err := n.client.RequestObj(endpoint, &authorize, &core.RequestOptions{
		Method: "GET",
	})
	if err != nil {
		return nil, err
	}

	return &authorize, nil
}
