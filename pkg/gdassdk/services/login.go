package services

import (
	core "github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/core/v1"
)

type LoginService struct {
	client *core.Client
}

func NewLoginService(client *core.Client) *LoginService {
	return &LoginService{
		client: client,
	}
}

func (l *LoginService) PostLogin(data *core.PostLogin) (*core.Login, error) {
	var login core.Login
	endpoint := "login"
	_, err := l.client.RequestObj(endpoint, &login, &core.RequestOptions{
		Method: "POST",
		Data:   StructToMapStr(data),
	})
	if err != nil {
		return nil, err
	}

	return &login, nil
}
