// Package setting 设置
package setting

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Setting 使用viper读取设置信息
type Setting struct {
	vp *viper.Viper
}

// NewSetting 读取配置信息
func NewSetting(name, path string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName(name)
	vp.AddConfigPath(path)
	vp.SetConfigType("yaml")
	if err := vp.ReadInConfig(); err != nil {
		return nil, err
	}
	s := &Setting{vp}
	s.WatchSettingChange()
	return s, nil
}

// WatchSettingChange 配置热更新
func (s *Setting) WatchSettingChange() {
	go func() {
		s.vp.WatchConfig()
		s.vp.OnConfigChange(func(_ fsnotify.Event) {
			_ = s.ReloadAllSections()
		})
	}()
}
