package motion_detector

import "testing"

type WebCamMotionDetectorMock struct {
	motionStatus bool
	motionError  error
}

func (w WebCamMotionDetectorMock) MotionDetected() (bool, error) {
	return w.motionStatus, w.motionError
}

func TestMotionDetected(t *testing.T) {
	mockTrue := WebCamMotionDetectorMock{motionStatus: true, motionError: nil}
	motion, connectErr := TriggerEvent(mockTrue)
	if connectErr != nil {
		t.Errorf("mockTrue error should be nil")
	} else {
		if motion != true {
			t.Errorf("motion should be true")
		}
	}
}

func TestMotionNotDetected(t *testing.T) {
	mockTrue := WebCamMotionDetectorMock{motionStatus: false, motionError: nil}
	motion, connectErr := TriggerEvent(mockTrue)
	if connectErr != nil {
		t.Errorf("mockTrue error should be nil")
	} else {
		if motion != false {
			t.Errorf("motion should be false")
		}
	}
}
