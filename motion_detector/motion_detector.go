package motion_detector

import (
	"net/http"

	"github.com/a-castellano/reolink-manager/webcam"
)

type MotionDetector interface {
	MotionDetected() (bool, error)
}

type WebcamMotionDetector struct {
	Client http.Client
	Webcam *webcam.Webcam
}

func (w WebcamMotionDetector) MotionDetected() (bool, error) {
	motion, connectErr := w.Webcam.MotionDetected(w.Client)
	return motion, connectErr
}

func TriggerEvent(detector MotionDetector) (bool, error) {
	// check motion
	motion, connectErr := detector.MotionDetected()
	return motion, connectErr
}
