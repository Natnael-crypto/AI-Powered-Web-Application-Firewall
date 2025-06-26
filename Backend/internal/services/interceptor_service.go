package services

import (
	"backend/internal/repository"
	"net/http"
)

func StartInterceptor() (string, int) {
	if repository.IsInterceptorRunning() {
		return "Interceptor is already running.", http.StatusOK
	}

	repository.StartInterceptor()
	return "Interceptor will start soon.", http.StatusOK
}

func StopInterceptor() (string, int) {
	if !repository.IsInterceptorRunning() {
		return "Interceptor is already stopped.", http.StatusOK
	}

	repository.StopInterceptor()
	return "Interceptor will stop soon.", http.StatusOK
}

func RestartInterceptor() (string, int) {
	repository.RestartInterceptor()
	return "Interceptor will restart soon.", http.StatusOK
}

func GetInterceptorState() repository.InterceptorState {
	return repository.GetInterceptorState()
}

func GetMlState() repository.MlState {
	return repository.GetMlState()
}
