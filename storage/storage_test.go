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

func TestCheckNotTriggered(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectGet("webcam").RedisNil()

	storageInstance := Storage{db, 15}
	var ctx = context.TODO()

	triggered, err := storageInstance.CheckTrigger(ctx, "webcam")
	if err != nil {
		t.Error("TesTestCheckNotTriggeredd should not fail. Error was ", err.Error())
	}

	if triggered != false {
		t.Error("TestCheckNotTriggered shold not be triggered.")
	}

}

func TestCheckTriggeredError(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectGet("webcam").SetErr(errors.New("error"))

	storageInstance := Storage{db, 15}
	var ctx = context.TODO()

	_, err := storageInstance.CheckTrigger(ctx, "webcam")
	if err == nil {
		t.Error("TestCheckTriggeredError should fail.")
	}

}

func TestCheckTriggered(t *testing.T) {
	db, mock := redismock.NewClientMock()

	mock.ExpectGet("webcam").SetVal("test value")

	storageInstance := Storage{db, 15}
	var ctx = context.TODO()

	triggered, err := storageInstance.CheckTrigger(ctx, "webcam")
	if err != nil {
		t.Error("TestCheckTriggered should not fail. Error was ", err.Error())
	}

	if triggered != true {
		t.Error("TestCheckTriggered shold not be triggered.")
	}

}
