package v1

type Login struct {
	Result         string `json:"result"`
	Token          string `json:"token"`
	UserAuth       int    `json:"userAuth"`
	Ak             string `json:"ak"`
	Sk             string `json:"sk"`
	ReMainErrCount int    `json:"re mainErrCount"`
	LastLoginTime  int64  `json:"lastLoginTime"`
}

// NodeLists 节点信息
type NodeList struct {
	Result   string         `json:"result"`
	NodeList []NodeListData `json:"nodeList"`
}

// NodeList 每个节点的信息
type NodeListData struct {
	IP      string  `json:"ip"`
	Status  float64 `json:"status"`
	DamName string  `json:"damName"`
}
