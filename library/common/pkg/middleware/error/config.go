package errmw

const configSectionName = "app"

type appConfig struct {
	Mode string `yaml:"mode"`
}
