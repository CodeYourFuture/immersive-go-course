package types

type Command struct {
	Description string   `json:"description" yaml:"description"`
	Schedule    string   `json:"schedule" yaml:"schedule"`
	Command     string   `json:"command" yaml:"command"`
	MaxRetries  int      `json:"max_retries" yaml:"max_retries"`
	Clusters    []string `json:"clusters" yaml:"clusters"`
}
