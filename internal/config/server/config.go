package agent

import "time"

type ServerConfig struct {
	Files []string `env:"FILES" envSeparator:":"`
	Home  string   `env:"HOME"`
	// required требует, чтобы переменная TASK_DURATION была определена
	TaskDuration time.Duration `env:"TASK_DURATION,required"`
}
