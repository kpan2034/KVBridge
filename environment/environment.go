package environment

import (
	"KVBridge/config"
	"KVBridge/log"
)

type Environment struct {
	log.Logger
	*config.Config
}

func New(logger log.Logger, config *config.Config) *Environment {
	return &Environment{
		logger,
		config,
	}
}

func (env *Environment) WithLogger(logger log.Logger) *Environment {
	return &Environment{
		logger,
		env.Config,
	}
}

func (env *Environment) WithConfig(config *config.Config) *Environment {
	return &Environment{
		env.Logger,
		config,
	}
}
