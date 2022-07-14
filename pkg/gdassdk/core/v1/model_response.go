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

type NodeCaches struct {
	Result          string `json:"result"`
	TotalCacheSize  int64  `json:"totalCacheSize"`
	UnUsedCacheSize int64  `json:"unUsedCacheSize"`
	UsedCacheSize   int64  `json:"usedCacheSize"`
}

type NodeDas struct {
	Result string   `json:"result"`
	DaList []DaList `json:"daList"`
}
type DriveSerialList struct {
	DriveNo     int    `json:"driveNo"`
	DriveSerial string `json:"driveSerial"`
}
type DaList struct {
	DaNo              int    `json:"da_no"`
	Name              string `json:"name"`
	ChangerNum        int    `json:"changer_num"`
	DriveNum          int    `json:"drive_num"`
	SlotNum           int    `json:"slot_num"`
	MagazineUsedCount int    `json:"magazineUsedCount"`
	MagazineExcpCount int    `json:"magazineExcpCount"`
	MagazineFreeCount int    `json:"magazineFreeCount"`
	DaStatus          int    `json:"daStatus"`
	IP                string `json:"ip"`
	// ChangerSmartInfo 机械手的smart信息(编号、状态、使用百分比)
	ChangerSmartInfo []ChangerSmartInfo `json:"changerSmartInfo"`
	// DriveSmartInfo 光驱的smart信息(编号、状态、使用百分比)
	DriveSmartInfo  []DriveSmartInfo  `json:"driveSmartInfo"`
	ChangerSerial   string            `json:"changerSerial"`
	DriveSerialList []DriveSerialList `json:"driveSerialList"`
}

// 盘库中每个机械手的信息
type ChangerSmartInfo struct {
	// UnitNo 机械手号，默认为0
	UnitNo int `json:"unitNo"`
	// UsedPercent 机械手使用百分比
	UsedPercent int `json:"usedPercent"`
	// Status 机械手状态
	Status int `json:"status"`
}

// 盘库中每个光驱的信息
type DriveSmartInfo struct {
	// UnitNo 光驱号，默认为0
	UnitNo int `json:"unitNo"`
	// UsedPercent 光驱使用百分比
	UsedPercent int `json:"usedPercent"`
	// Status 光驱状态
	Status int `json:"status"`
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
