package object

import (
	"k8s/pkg/util/HTTPClient"
	"time"
)

type Function struct {
	Name      string
	Path      string
	ImageName string
	Version   int

	Image string
}

type RunningFunction struct {
	Function  Function
	ServiceIP string
	Timer     *time.Timer
	Client    *HTTPClient.Client
}
