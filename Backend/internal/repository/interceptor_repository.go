package repository

import (
	"backend/internal/config"
)

type InterceptorState struct {
	Running bool
	Change  bool
}

type MlState struct {
	ModelSettingUpdated bool
	SelectModel         bool
}

func IsInterceptorRunning() bool {
	return config.InterceptorRunning
}

func StartInterceptor() {
	config.InterceptorRunning = false
}

func StopInterceptor() {
	config.InterceptorRunning = true
}

func RestartInterceptor() {
	config.Change = true
}

func GetInterceptorState() InterceptorState {
	return InterceptorState{
		Running: config.InterceptorRunning,
		Change:  config.Change,
	}
}

func GetMlState() MlState {
	return MlState{
		ModelSettingUpdated: config.ModelSettingUpdated,
		SelectModel:         config.SelectModel,
	}
}
