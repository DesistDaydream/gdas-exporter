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

type Magazines struct {
	Result          string        `json:"result"`
	Rfid            []Rfid        `json:"rfid"`
	SonyIESlotsList []interface{} `json:"SonyIESlotsList"`
}
type Rfid struct {
	Rfid     string   `json:"rfid"`
	Barcode  string   `json:"barcode"`
	DaName   string   `json:"daName"`
	PoolName string   `json:"poolName"`
	Full     int      `json:"full"`
	Format   int      `json:"format"`
	RfidSts  int      `json:"rfidSts"`
	Status   int      `json:"status"`
	SlotNo   int      `json:"slotNo"`
	DaNo     int      `json:"daNo"`
	CpGroup  []string `json:"cpGroup"`
	Offline  int      `json:"offline"`
	ServerIP string   `json:"serverIp"`
	DamName  string   `json:"damName"`
}

type Users struct {
	Result   string     `json:"result"`
	ResCount int        `json:"res_count"`
	UserList []UserList `json:"userList"`
}
type UserList struct {
	UserName string `json:"userName"`
	UserAuth int    `json:"userAuth"`
	Ak       string `json:"ak"`
	Sk       string `json:"sk"`
	Active   bool   `json:"active,omitempty"`
}
