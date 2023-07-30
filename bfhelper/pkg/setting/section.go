package setting

var sections = make(map[string]any)

// ReadSection 根据给定的建造的写入map section
func (s *Setting) ReadSection(k string, v any) error {
	if err := s.vp.UnmarshalKey(k, v); err != nil {
		return err
	}
	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

// ReloadAllSections 重载所有键值对
func (s *Setting) ReloadAllSections() error {
	for k, v := range sections {
		err := s.ReadSection(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
