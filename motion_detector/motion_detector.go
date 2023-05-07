package motion_detector

import (
	"net/http"

	"github.com/a-castellano/reolink-manager/webcam"
)

type MotionDetector interface {
	MotionDetected() (bool, error)
}

type WebCamMotionDetector struct {
	client http.Client
	webcam webcam.Webcam
}

func (w WebCamMotionDetector) MotionDetected() (bool, error) {
	motion, connectErr := w.webcam.MotionDetected(w.client)
	return motion, connectErr
}

func TriggerEvent(detector MotionDetector) (bool, error) {
	// check motion
	motion, connectErr := detector.MotionDetected()
	return motion, connectErr
}
