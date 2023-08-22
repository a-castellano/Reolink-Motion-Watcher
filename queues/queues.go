package queues

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strconv"

	config "github.com/a-castellano/Reolink-Motion-Watcher/config_reader"
	"github.com/streadway/amqp"
)

type Notification struct {
	Timestamp  int64  `json:"timestamp"`
	WebCamName string `json:"webcamName"`
}

func EncodeNotification(notification Notification) ([]byte, error) {
	var encodedNotification []byte
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(notification)
	if err != nil {
		return encodedNotification, err
	}
	encodedNotification = network.Bytes()
	return encodedNotification, nil
}

func DecodeNotification(encoded []byte) (Notification, error) {
	var notification Notification
	network := bytes.NewBuffer(encoded)
	dec := gob.NewDecoder(network)
	err := dec.Decode(&notification)
	if err != nil {
		return Notification{}, err
	}
	return notification, nil
}

func SendNotification(rabbitmqConfig config.Rabbitmq, webcamInfoName string) error {

	notificationString := fmt.Sprintf("Motion detected by webcam %s\n", webcamInfoName)
	encodedNotification := []byte(notificationString)

	connection_string := "amqp://" + rabbitmqConfig.User + ":" + rabbitmqConfig.Password + "@" + rabbitmqConfig.Host + ":" + strconv.Itoa(rabbitmqConfig.Port) + "/"

	conn, errDial := amqp.Dial(connection_string)
	defer conn.Close()

	if errDial != nil {
		return errDial
	}

	channel, errChannel := conn.Channel()
	defer channel.Close()
	if errChannel != nil {
		return errChannel
	}

	queue, errQueue := channel.QueueDeclare(
		rabbitmqConfig.MotionQueueName, // name
		true,                           // durable
		false,                          // delete when unused
		false,                          // exclusive
		false,                          // no-wait
		nil,                            // arguments
	)

	if errQueue != nil {
		return errQueue
	}

	// send Notification

	err := channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         encodedNotification,
		})

	if err != nil {
		return err
	}
	return nil
}
