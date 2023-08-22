package queues

import "testing"

func TestEncodeAndDecodeNotification(t *testing.T) {

	var notification Notification
	notification.Timestamp = 4
	notification.WebCamName = "test"
	test, _ := EncodeNotification(notification)
	result, _ := DecodeNotification(test)

	if result.Timestamp != notification.Timestamp {
		t.Errorf(`Encode failed.`)
	}

	if result.WebCamName != notification.WebCamName {
		t.Errorf(`Encode failed.`)
	}

}

func TestDecodeEmptyDataJobs(t *testing.T) {

	var emptyData []byte
	_, err := DecodeNotification(emptyData)
	if err == nil {
		t.Error("Empty data decoding should fail.")
	}
}
