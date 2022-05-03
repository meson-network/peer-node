package configuration

import (
	"io/ioutil"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

var Config *VConfig

type VConfig struct {
	*viper.Viper
	configPath string
}

func ReadConfig(configPath string) (*VConfig, error) {
	c := &VConfig{viper.New(), ""}
	c.SetConfigFile(configPath)
	err := c.ReadInConfig()
	if err != nil {
		return nil, err
	}
	c.configPath = configPath
	return c, nil
}

func (c *VConfig) Get(key string, defaultValue interface{}) interface{} {
	if !c.Viper.IsSet(key) {
		return defaultValue
	}
	return c.Viper.Get(key)
}

func (c *VConfig) GetBool(key string, defaultValue bool) (bool, error) {
	if !c.Viper.IsSet(key) {
		return defaultValue, nil
	}

	v := c.Viper.Get(key)
	value, err := cast.ToBoolE(v)
	if err != nil {
		return false, err
	}
	return value, nil
}

func (c *VConfig) GetFloat64(key string, defaultValue float64) (float64, error) {
	if !c.Viper.IsSet(key) {
		return defaultValue, nil
	}

	v := c.Viper.Get(key)
	value, err := cast.ToFloat64E(v)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (c *VConfig) GetInt(key string, defaultValue int) (int, error) {
	if !c.Viper.IsSet(key) {
		return defaultValue, nil
	}

	v := c.Viper.GetInt(key)
	value, err := cast.ToIntE(v)
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (c *VConfig) GetIntSlice(key string, defaultValue []int) ([]int, error) {
	if !c.Viper.IsSet(key) {
		return defaultValue, nil
	}

	v := c.Viper.Get(key)
	value, err := cast.ToIntSliceE(v)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *VConfig) GetString(key string, defaultValue string) (string, error) {
	if !c.Viper.IsSet(key) {
		return defaultValue, nil
	}

	v := c.Viper.Get(key)
	value, err := cast.ToStringE(v)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (c *VConfig) GetStringMap(key string, defaultValue map[string]interface{}) (map[string]interface{}, error) {
	if !c.Viper.IsSet(key) {
		return defaultValue, nil
	}

	v := c.Viper.Get(key)
	value, err := cast.ToStringMapE(v)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *VConfig) GetStringSlice(key string, defaultValue []string) ([]string, error) {
	if !c.Viper.IsSet(key) {
		return defaultValue, nil
	}

	v := c.Viper.Get(key)
	value, err := cast.ToStringSliceE(v)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (c *VConfig) GetConfigAsString() (string, error) {
	b, err := ioutil.ReadFile(c.configPath)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
