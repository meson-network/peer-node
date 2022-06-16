package conf

type TomlConfig struct {
	Daemon_name string  `toml:"daemon_name"`
	Log_level   string  `toml:"log_level"`
	Token       string  `toml:"token"`
	Https_port  int     `toml:"https_port"`
	EndPoint    string  `toml:"end_point"`
	Api         API     `toml:"api"`
	Cache       Cache   `toml:"cache"`
	Storage     Storage `toml:"storage"`
}

type API struct {
	Doc_gen_search_dir string `toml:"doc_gen_search_dir"`
	Doc_gen_mainfile   string `toml:"doc_gen_mainfile"`
	Doc_gen_output_dir string `toml:"doc_gen_output_dir"`
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
