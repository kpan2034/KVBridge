package environment

import (
	"KVBridge/config"
	"KVBridge/log"
	"KVBridge/state"
)

type Environment struct {
	log.Logger
	*config.Config
	*state.State
}

func New(logger log.Logger, config *config.Config, state *state.State) *Environment {
	return &Environment{
		logger,
		config,
		state,
	}
}

func (env *Environment) WithLogger(logger log.Logger) *Environment {
	newEnv := *env
	newEnv.Logger = logger
	return &newEnv
}

func (env *Environment) WithConfig(config *config.Config) *Environment {
	newEnv := *env
	newEnv.Config = config
	return &newEnv
}
