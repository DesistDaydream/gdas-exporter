package v1

// NodeList 分布式节点信息
type NodeList struct {
	Result   string         `json:"result"`
	NodeList []NodeListData `json:"nodeList"`
}

type NodeListData struct {
	IP      string  `json:"ip"`
	Status  float64 `json:"status"`
	DamName string  `json:"damName"`
}

// 缓存容量
type NodeCaches struct {
	Result          string `json:"result"`
	TotalCacheSize  int64  `json:"totalCacheSize"`
	UnUsedCacheSize int64  `json:"unUsedCacheSize"`
	UsedCacheSize   int64  `json:"usedCacheSize"`
}

// 节点盘库列表
type NodeDas struct {
	Result string   `json:"result"`
	DaList []DaList `json:"daList"`
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

type DriveSerialList struct {
	DriveNo     int    `json:"driveNo"`
	DriveSerial string `json:"driveSerial"`
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

// 授权
type Authorize struct {
	Result               string `json:"result"`
	RegisterTimeMillis   string `json:"registerTimeMillis"`
	RegisterDays         string `json:"registerDays"`
	RegisterSerialnumber string `json:"registerSerialnumber "`
	RegisterMgzCount     string `json:"registerMgzCount"`
}

// 盘匣列表
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

// 全局存储空间信息
type Totalspace struct {
	Result              string `json:"result"`
	TotalSpaceRaid0     int64  `json:"totalSpaceRaid0"`
	TotalSpace          int64  `json:"totalSpace"`
	TotalAvailableSpace int64  `json:"totalAvailableSpace"`
	TotalMgzCount       int    `json:"totalMgzCount"`
	UsedMgzCount        int    `json:"usedMgzCount"`
	FreeMgzCount        int    `json:"freeMgzCount"`
	ExceptionMgzCount   int    `json:"exceptionMgzCount"`
	TotalSlotCount      int    `json:"totalSlotCount"`
}

// 用户信息
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

// 所有盘匣组信息
type Pools struct {
	Result   string      `json:"result"`
	ResCount int         `json:"res_count"`
	Pools    []PoolsData `json:"pools"`
}

type PoolsData struct {
	PoolName           string `json:"pool_name"`
	Type               string `json:"type"`
	PoolRaidLvl        int    `json:"pool_raidLvl"`
	RfidCount          int    `json:"rfid_count,omitempty"`
	PoolTotalSpace     int64  `json:"pool_total_space"`
	PoolAvailableSpace int64  `json:"pool_available_space"`
	PoolSts            int    `json:" pool_sts,omitempty"`
	PoolCanDelFlag     string `json:"poolCanDel_flag"`
	User               string `json:"user"`
	AutoAddMgz         bool   `json:"autoAddMgz"`
	DefaultMgz         bool   `json:"defaultMgz"`
	PoolOperation      string `json:"pool_operation,omitempty"`
	Duplicate          int    `json:"duplicate"`
	RfidCoUnt          int    `json:"rfid_co unt,omitempty"`
	PoolOperaTion      string `json:"pool_opera tion,omitempty"`
}

// 全局盘库信息
type Das struct {
	Result string   `json:"result"`
	DaInfo []DaInfo `json:"daInfo"`
}

type DaInfo struct {
	// DamName 盘库所在节点名称
	DamName string `json:"damName"`
	// ChangerSerialNumber ???盘库序列号????这不是机械手的信息么???
	ChangerSerialNumber string `json:"changerSerialNumber"`
	// IP 盘库所在节点 IP
	IP string `json:"ip"`
	// DaName 盘库型号
	DaName string `json:"daName"`
	// DaStatus 盘库状态。0 正常，-203 盘匣弹出中，-210 仓架解锁中，-202 系统繁忙，-102 断开连接，-100和-103 识别中
	DaStatus int `json:"daStatus"`
	// DaNo ！！！未知！！！好像每个节点第一个盘库的号就是0，所以看到的就都是0
	DaNo int `json:"daNo"`
	// DaVendor 厂商信息
	DaVendor string `json:"daVendor"`
	// Offline 盘库注册、断开状态。0 断开，1 已注册
	Offline int `json:"offline"`
	// SlotCount 槽位总数
	SlotCount int `json:"slotCount"`
}
