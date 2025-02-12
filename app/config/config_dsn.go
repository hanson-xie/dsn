package config

type DsnConfig struct {
	Log            LogConfig         `yaml:"log"`
	Rpc            string            `yaml:"rpc"`
	DsnServers     map[string]string `yaml:"dsn_servers"`
	DocAuth        map[string]string `yaml:"doc_auth"`
	TomlDir        string            `yaml:"toml_dir"`
	GitUpdateShell string            `yaml:"git_update_shell"`
	ReloadFlag     bool              `yaml:"reload_flag"`
	CheckSpec      string            `yaml:"check_spec"`
}

type LogConfig struct {
	Level      int8   `yaml:"level"`
	LogDir     string `yaml:"log_dir"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
	LocalTime  bool   `yaml:"local_time"`
}
