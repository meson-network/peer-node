package conf

type TomlConfig struct {
	Log        Log     `toml:"log"`
	Token      string  `toml:"token"`
	Https_port int     `toml:"https_port"`
	EndPoint   string  `toml:"end_point"`
	Cache      Cache   `toml:"cache"`
	Storage    Storage `toml:"storage"`
}

type Log struct {
	Level string `toml:"level"`
}

type Cache struct {
	Folder string `toml:"folder"`
	Size   int    `toml:"size"`
}

type Storage struct {
	Enable       bool   `toml:"enable"`
	Api_port     int    `toml:"api_port"`
	Console_port int    `toml:"console_port"`
	Folder       string `toml:"folder"`
	Password     string `toml:"password"`
}
