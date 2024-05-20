package message

type CronjobMessage struct {
	Command string
	Exectime string // ISO 8601
	Name string
	Retries int
}