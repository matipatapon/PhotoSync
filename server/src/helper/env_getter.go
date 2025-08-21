package helper

import (
	"os"
)

type IEnvGetter interface {
	Get(string) string
}

type EnvGetter struct {
}

func NewEnvGetter() EnvGetter {
	return EnvGetter{}
}

func (eg *EnvGetter) Get(name string) string {
	return os.Getenv(name)
}
