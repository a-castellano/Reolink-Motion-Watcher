package storage

import (
	"context"
	"errors"
	"testing"
	"time"

	redismock "github.com/go-redis/redismock/v9"
)

func TestUpdateTriggerFailed(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectSet("webcam", "triggered", 15*time.Second).SetErr(errors.New("error"))

	storageInstance := Storage{db, 15}
	var ctx = context.TODO()

	_, err := storageInstance.UpdateTrigger(ctx, "webcam")
	if err == nil {
		t.Error("TestUpdateTriggerFailed should fail.")
	}

}

func TestUpdateTrigger(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectSet("webcam", "triggered", 15*time.Second).SetVal("OK")

	storageInstance := Storage{db, 15}
	var ctx = context.TODO()

	result, err := storageInstance.UpdateTrigger(ctx, "webcam")
	if err != nil {
		t.Error("TestUpdateTrigger should not fail. Error was ", err.Error())
	}

	if result != true {
		t.Error("TestUpdateTrigger should be true, not false.")
	}

}
