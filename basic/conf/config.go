package conf

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/coreservice-io/utils/path_util"
	"github.com/meson-network/peer-node/basic"
	"github.com/pelletier/go-toml"
)

/////////////////////////////
type Config struct {
	Toml_config *TomlConfig
	Abs_path    string
}

var config *Config

func Get_config() *Config {
	return config
}

func (config *Config) Read_config_file() (string, error) {

	doc, err := ioutil.ReadFile(config.Abs_path)
	if err != nil {
		return "", err
	}

	return string(doc), nil
}

func (config *Config) Save_config() error {

	result, err := toml.Marshal(config.Toml_config)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(config.Abs_path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	err = f.Truncate(0)
	if err != nil {
		return err
	}
	_, err = f.Seek(0, 0)
	if err != nil {
		return err
	}

	_, err = f.Write(result)
	if err != nil {
		return err
	}

	return nil
}

func Init_config(config_path string) error {

	if config != nil {
		return nil
	}

	c_p, c_p_exist, _ := path_util.SmartPathExist(config_path)
	if !c_p_exist {
		return errors.New("no config file:" + config_path)
	}

	var cfg Config
	cfg.Abs_path = c_p
	cfg.Toml_config = &TomlConfig{}

	config_str, err := cfg.Read_config_file()
	if err != nil {
		return err
	}

	err = toml.Unmarshal([]byte(config_str), cfg.Toml_config)
	if err != nil {
		return err
	}

	basic.Logger.Infoln("using config:", cfg.Abs_path)

	config = &cfg

	return nil
}
