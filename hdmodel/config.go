package hdmodel

type Kitex struct {
	Service         string `yaml:"service"`
	Address         string `yaml:"address"`
	MetricsPort     string `yaml:"metrics_port"`
	EnablePprof     bool   `yaml:"enable_pprof"`
	EnableGzip      bool   `yaml:"enable_gzip"`
	EnableAccessLog bool   `yaml:"enable_access_log"`
	LogLevel        string `yaml:"log_level"`
	LogFileName     string `yaml:"log_file_name"`
	LogMaxSize      int    `yaml:"log_max_size"`
	LogMaxBackups   int    `yaml:"log_max_backups"`
	LogMaxAge       int    `yaml:"log_max_age"`
}
type Prometheus struct {
	Enable      bool `yaml:"enable"`
	MetricsPort int  `yaml:"metrics_port"`
}

type OTel struct {
	Enable   bool   `yaml:"enable"`
	Endpoint string `yaml:"endpoint"`
}
type Registry struct {
	RegistryAddress string `yaml:"registry_address"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	NamespaceId     string `yaml:"namespace_id"`
	Group           string `yaml:"group"`
	DataId          string `yaml:"data_id"`
}

type Monitor struct {
	Enabled       bool       `yaml:"enabled"`
	OTel          OTel       `yaml:"otel"`
	Prometheus    Prometheus `yaml:"prometheus"`
	Registry      Registry   `yaml:"registry"`
	EnableTracing bool       `yaml:"enable_tracing"`
}
