package v1

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

type Authorize struct {
	Result               string `json:"result"`
	RegisterTimeMillis   string `json:"registerTimeMillis"`
	RegisterDays         string `json:"registerDays"`
	RegisterSerialnumber string `json:"registerSerialnumber "`
	RegisterMgzCount     string `json:"registerMgzCount"`
}
