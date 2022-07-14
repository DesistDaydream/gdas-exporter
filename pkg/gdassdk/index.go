package gdassdk

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	core "github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/core/v1"
	"github.com/DesistDaydream/gdas-exporter/pkg/gdassdk/services"
)

type PostLogin struct {
	Username string `json:"userName"`
	Password string `json:"passWord"`
}

type Login struct {
	Result         string `json:"result"`
	Token          string `json:"token"`
	UserAuth       int    `json:"userAuth"`
	Ak             string `json:"ak"`
	Sk             string `json:"sk"`
	ReMainErrCount int    `json:"re mainErrCount"`
	LastLoginTime  int64  `json:"lastLoginTime"`
}

func GetToken(prefix string, username, password string) (string, error) {
	url := fmt.Sprintf("%s/v1/login", prefix)
	// 解析请求体

	var (
		reqBody *bytes.Buffer
		req     *http.Request
		err     error
	)

	requestBody, err := json.Marshal(PostLogin{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", err
	}
	reqBody = bytes.NewBuffer(requestBody)

	req, err = http.NewRequest("POST", url, reqBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Referer", fmt.Sprintf("%v/gdas", prefix))
	req.Header.Set("stime", fmt.Sprintf("%v", time.Now().UnixNano()/1e6))

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", errors.New(string(body))
	}

	// 解析 body
	var login Login
	err = json.Unmarshal(body, &login)
	if err != nil {
		return "", err
	}

	return login.Token, nil
}

// Service encapsulate authenticated token
type Service struct {
	Client *core.Client
	Node   *services.NodeService
}

// NewService create Client for external use
func NewService(prefix string, token string) *Service {
	s := new(Service)
	s.Init(prefix, token)
	return s
}

func (s *Service) Init(prefix string, token string) {
	s.Client = core.NewClient(prefix, token)
	s.Node = services.NewNodeService(s.Client)
}
