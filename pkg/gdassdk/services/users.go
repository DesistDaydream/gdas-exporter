package services

import (
	core "github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/core/v1"
)

type UsersService struct {
	client *core.Client
}

func NewUsersService(client *core.Client) *UsersService {
	return &UsersService{
		client: client,
	}
}

// GetUsers 23.查询用户
func (n *UsersService) GetUsers() (*core.Users, error) {
	var data core.Users
	endpoint := "users"
	_, err := n.client.RequestObj(endpoint, &data, &core.RequestOptions{
		Method:   "GET",
		RawQuery: "pages=1-*&userName=",
	})
	if err != nil {
		return nil, err
	}

	return &data, nil
}
