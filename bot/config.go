package bot

type Config struct {
	Server   string   `yaml:"server"`
	Port     string   `yaml:"port"`
	Nick     string   `yaml:"nickname"`
	Identity string   `yaml:"identity"`
	GECOS    string   `yaml:"GECOS"`
	Channels []string `yaml:"channels"`
	Realname string   `yaml:"realname"`
}
