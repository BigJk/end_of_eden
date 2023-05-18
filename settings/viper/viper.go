package viper

import "github.com/spf13/viper"

// Viper is a wrapper around viper that implements the Settings interface.
type Viper struct {
	SettingsName string
}

func (v Viper) LoadSettings() error {
	viper.SetConfigName(v.SettingsName)
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			_ = viper.SafeWriteConfigAs("./" + v.SettingsName + ".toml")
		} else {
			return err
		}
	}

	return nil
}

func (v Viper) SaveSettings() error {
	return viper.WriteConfigAs("./" + v.SettingsName + ".toml")
}

func (v Viper) Get(key string) any {
	return viper.Get(key)
}

func (v Viper) GetString(key string) string {
	return viper.GetString(key)
}

func (v Viper) GetStrings(key string) []string {
	return viper.GetStringSlice(key)
}

func (v Viper) GetInt(key string) int {
	return viper.GetInt(key)
}

func (v Viper) GetFloat(key string) float64 {
	return viper.GetFloat64(key)
}

func (v Viper) GetBool(key string) bool {
	return viper.GetBool(key)
}

func (v Viper) Set(key string, value any) {
	viper.Set(key, value)
}

func (v Viper) GetKeys() []string {
	return viper.AllKeys()
}

func (v Viper) SetDefault(key string, value any) {
	viper.SetDefault(key, value)
}
