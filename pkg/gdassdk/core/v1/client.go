package v1

import (
	"bytes"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

type Client struct {
	Prefix  string
	Token   string
	timeout time.Duration
}

type RequestOptions struct {
	Method   string
	ReqBody  map[string]string
	RawQuery string
}

func NewClient(prefix string, token string, timeout time.Duration) *Client {
	return &Client{
		Prefix:  prefix,
		Token:   token,
		timeout: timeout,
	}
}

func (c *Client) request(endpoint string, options *RequestOptions) ([]byte, error) {
	url := fmt.Sprintf("%s/v1/%s", c.Prefix, endpoint)

	method := options.Method
	if len(method) == 0 {
		method = "GET"
	}

	randString, signatureSha := GenerateSign(c.Token)

	var (
		reqBody *bytes.Buffer
		req     *http.Request
		err     error
	)

	if len(options.ReqBody) > 0 {
		requestBody, err := json.Marshal(options.ReqBody)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(requestBody)
	}
	if reqBody != nil {
		req, err = http.NewRequest(method, url, reqBody)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	req.Header.Add("token", c.Token)
	req.Header.Set("stime", fmt.Sprintf("%v", time.Now().UnixNano()/1e6))
	req.Header.Set("nonce", randString)
	req.Header.Set("signature", signatureSha)
	req.Header.Set("Referer", fmt.Sprintf("%v/gdas", c.Prefix))
	req.URL.RawQuery = options.RawQuery

	logrus.WithField("url", req.URL).Debug("检查 URL")

	// requestDump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return nil, err
	// }
	// L.Debug(string(requestDump))

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   c.timeout,
		Transport: transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	return body, nil
}

func (c *Client) RequestObj(endpoint string, container interface{}, options *RequestOptions) (interface{}, error) {
	var (
		body []byte
		err  error
	)

	body, err = c.request(endpoint, options)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, container)
	if err != nil {
		return nil, err
	}

	return container, nil
}

// 生成签名所需数据
func GenerateSign(token string) (string, string) {
	// 毫秒时间戳
	stime := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	// 随机数
	randString := rand.Intn(100000)
	// 随机数倒序
	stringRand := []rune(strconv.Itoa(randString))
	for from, to := 0, len(stringRand)-1; from < to; from, to = from+1, to-1 {
		stringRand[from], stringRand[to] = stringRand[to], stringRand[from]
	}
	// 签名
	signature := stime + strconv.Itoa(randString) + token + string(stringRand)
	h := sha256.New()
	h.Write([]byte(signature))                     // 需要加密的字符串为
	signatureSha := hex.EncodeToString(h.Sum(nil)) // 输出加密结果

	return strconv.Itoa(randString), signatureSha
}
