package main

import (
	"context"
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"time"

	apiwatcher "github.com/a-castellano/AlarmStatusWatcher/apiwatcher"
	config "github.com/a-castellano/Reolink-Motion-Watcher/config_reader"
	"github.com/a-castellano/Reolink-Motion-Watcher/motion_detector"
	"github.com/a-castellano/Reolink-Motion-Watcher/queues"
	"github.com/a-castellano/Reolink-Motion-Watcher/storage"
	"github.com/a-castellano/reolink-manager/webcam"
	goredis "github.com/redis/go-redis/v9"
)

// watchMotionSensor
// For each webcam, checks if motion is triggered and update storage when alarm is bot disarmed

func watchMotionSensor(ctx context.Context, webcams map[string]*webcam.Webcam, storageInstance storage.Storage, watcher apiwatcher.APIWatcher, alarmManagerRequester apiwatcher.Requester, alarmDeviceID string) {

	motionDetectors := make(map[string]*motion_detector.WebcamMotionDetector)

	for webcamName, webcam := range webcams {
		log.Println("Creating ", webcamName, " motionDetector")
		httpClient := http.Client{
			Timeout: time.Second * 5, // Maximum of 5 secs
		}
		motionDetector := motion_detector.WebcamMotionDetector{Client: httpClient, Webcam: webcam}
		motionDetectors[webcamName] = &motionDetector
	}

	for {
		for webcamName, motionDetector := range motionDetectors {
			motion, connectErr := motionDetector.MotionDetected()
			if connectErr != nil {
				log.Println("Error found in ", webcamName, connectErr)
				log.Fatal(connectErr)
			} else {

				_, apiInfoErr := watcher.ShowInfo(alarmManagerRequester)
				if apiInfoErr != nil {
					apiErrorString := fmt.Sprintf("%v", apiInfoErr.Error())
					log.Fatal(apiErrorString)
				}

				if motion {
					log.Println("Motion detected in ", webcamName, "checking alarm status")
					apiInfo, apiInfoErr := watcher.ShowInfo(alarmManagerRequester)
					if apiInfoErr != nil {
						apiErrorString := fmt.Sprintf("%v", apiInfoErr.Error())
						log.Fatal(apiErrorString)
					}
					if apiInfo.DevicesInfo[alarmDeviceID].Mode == "disarmed" { // Debug
						log.Println("Alarm status is activated, update storage.")
						storageInstance.UpdateTrigger(ctx, webcamName)
					}
				}
			}
		}
		time.Sleep(1 * time.Second)
	}

}

// sendNotificationsOnMotion
// Sends rabbitmq message is any motion has been triggered

func sendNotificationsOnMotion(ctx context.Context, storageInstance storage.Storage, webcams map[string]*webcam.Webcam, notificationOnMotionChannel chan bool, rabbitmqConfig config.Rabbitmq) error {
	currentValues := make(map[string]bool)
	for {
		for webcamName := range webcams {
			notify := false
			if triggered, ok := currentValues[webcamName]; ok {
				currentValue, _ := storageInstance.CheckTrigger(ctx, webcamName)
				currentValues[webcamName] = currentValue
				if currentValue == true && triggered != currentValue {
					notify = true
				} // was not been triggered since now'
			} else { //Check current key
				currentValues[webcamName] = false
			}
			if notify {
				queues.SendNotification(rabbitmqConfig, webcamName)
			}

		}
		time.Sleep(1 * time.Second)
	}
	return nil
}

func main() {

	logwriter, e := syslog.New(syslog.LOG_NOTICE, "windmaker-reolink-motion-watcher")
	if e == nil {
		log.SetOutput(logwriter)
		// Remove date prefix
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	}

	log.Println("Reading service config.")
	serviceConfig, errConfig := config.ReadConfig()

	if errConfig != nil {
		log.Fatal(errConfig)
		panic(errConfig)
	}

	log.Println("Service config readed.")
	httpClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 secs
	}

	log.Println("Setting up redis connection.")

	redisAddress := fmt.Sprintf("%s:%d", serviceConfig.RedisInstance.Host, serviceConfig.RedisInstance.Port)
	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     redisAddress,
		Password: serviceConfig.RedisInstance.Password,
		DB:       serviceConfig.RedisInstance.Database,
	})

	ctx := context.Background()

	redisErr := redisClient.Set(ctx, "checkKey", "key", 1000000).Err()
	if redisErr != nil {
		log.Fatal(redisErr)
		panic(redisErr)
	}
	log.Println("Redis connection was successful.")

	storageInstance := storage.Storage{RedisClient: redisClient, TTL: serviceConfig.RedisInstance.TTL}

	log.Println("Establishing connection with alarmManager.")
	watcher := apiwatcher.APIWatcher{Host: serviceConfig.AlarmManager.Host, Port: serviceConfig.AlarmManager.Port}
	alarmManagerRequester := apiwatcher.Requester{Client: httpClient}
	_, apiInfoErr := watcher.ShowInfo(alarmManagerRequester)
	if apiInfoErr != nil {
		errorString := fmt.Sprintf("%v", apiInfoErr.Error())
		log.Fatal(errorString)
		panic(apiInfoErr)
	}

	go watchMotionSensor(ctx, serviceConfig.Webcams, storageInstance, watcher, alarmManagerRequester, serviceConfig.AlarmManager.DeviceId)

	notificationOnMotionChannel := make(chan bool)
	sendNotificationsOnMotion(ctx, storageInstance, serviceConfig.Webcams, notificationOnMotionChannel, serviceConfig.Rabbitmq)

	<-notificationOnMotionChannel

}
