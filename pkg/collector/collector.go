package collector

// 这三个常量用于给每个 Metrics 名字添加前缀
const (
	name      = "gdas_exporter"
	Namespace = "gdas"
	//Subsystem(s).
)

// Name 用于给前端页面显示 const 常量中定义的内容
func Name() string {
	return name
}
