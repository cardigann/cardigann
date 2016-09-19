package config

type ArrayConfig map[string]map[string]string

func (a ArrayConfig) Get(section, key string) (string, bool, error) {
	v, ok := a[section][key]
	return v, ok, nil
}

func (a ArrayConfig) Set(section, key, value string) error {
	if _, ok := a[section]; !ok {
		a[section] = map[string]string{}
	}
	a[section][key] = value
	return nil
}

func (a ArrayConfig) Sections() ([]string, error) {
	sections := []string{}
	for k := range a {
		sections = append(sections, k)
	}
	return sections, nil
}

func (a ArrayConfig) Section(section string) (map[string]string, error) {
	return a[section], nil
}
