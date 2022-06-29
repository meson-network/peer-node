package conf

type TomlConfig struct {
	Log        Log    `toml:"log"`
	Token      string `toml:"token"`
	Https_port int    `toml:"https_port"`
	EndPoint   string `toml:"end_point"`
	Cache      Cache  `toml:"cache"`
}

type Log struct {
	Level string `toml:"level"`
}

type Cache struct {
	Folder string `toml:"folder"`
	Size   int    `toml:"size"`
}
