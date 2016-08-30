package config

type Config interface {
	Get(section, key string) (string, bool, error)
	Set(section, key, value string) error
	Sections() ([]string, error)
	Section(section string) (map[string]string, error)
}

func IsSectionEnabled(section string, c Config) bool {
	v, _, err := c.Get(section, "enabled")
	if err != nil {
		return false
	}

	return v == "ok" || v == "true"
}
