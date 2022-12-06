package types

type Command struct {
	Command     string   `json:"command" yaml:"command"`
	Clusters    []string `json:"clusters" yaml:"clusters"`
	Description string   `json:"description" yaml:"description"`
	Schedule    string   `json:"schedule" yaml:"schedule"`
	MaxRetries  int      `json:"max_retries" yaml:"max_retries"`
}
