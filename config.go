package batron

type Config struct {
	JobDefinition string   `json:"job_definition"`
	Plugins       []Plugin `json:"plugins"`
}

type Plugin struct {
	Name   string       `json:"name"`
	Config PluginConfig `json:"config"`
}

type PluginConfig struct {
	Url string `json:"url"`
}
