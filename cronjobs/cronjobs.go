package cronjobs

import (
	"encoding/json"
	"fmt"
	"log"
	"sms-scheduler/types"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/streadway/amqp"
)

// Define a global variable to hold the RabbitMQ connection
var (
	rabbitMQConn         *amqp.Connection
	rabbitMQConnMutex    sync.Mutex
	sourceQueueName      = "SCHEDULED_SMS_QUEUE"
	destinationQueueName = "SMS_QUEUE"
)

func ScheduleCronJobs(rabbitMQConn *amqp.Connection) {

	// Initialize a cron scheduler
	c := cron.New()

	// Define a cron job to execute every minute (adjust as needed)
	_, err := c.AddFunc("* * * * *", func() {
		consumeScheduledMessages(rabbitMQConn)
	})

	if err != nil {
		log.Fatalf("Failed to add cron job: %v", err)
	}

	// Start the cron scheduler
	c.Start()
}

func consumeScheduledMessages(rabbitMQConn *amqp.Connection) {
	fmt.Println("CONSUMING SCHEDULED_SMS_QUEUE")
	// Create a channel
	ch, err := rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("FAILED TO OPEN A CHANNEL: %v", err)
		return
	}
	defer ch.Close()

	// Declare the queue
	if err := declareQueue(ch); err != nil {
		log.Fatalf("FAILED TO DECLARE THE QUEUE: %v", err)
	}

	msgs, err := ch.Consume(
		sourceQueueName,
		"",
		false, // manually acknowledge the queue
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("FAILED TO CONSUME MESSAGES FROM THE SOURCE QUEUE: %v", err)
	}

	for msg := range msgs {
		// Process scheduled messages here
		var message types.ScheduledMessage
		err := json.Unmarshal(msg.Body, &message)
		if err != nil {
			log.Printf("FAILED TO UNMARSHAL MESSAGE: %v", err)
			// If there's an error, you can decide whether to acknowledge or not
			msg.Ack(false) // Manually acknowledge the message
			continue
		}

		now := time.Now().UTC()
		fmt.Printf("MESSAGE SEND TIME: %s\n", message.SendTime)
		fmt.Printf("NOW TIME: %s\n", now)
		// Check if it's time to send the message
		if now.After(message.SendTime.UTC()) {
			// Perform any necessary processing for the scheduled message
			fmt.Printf("RECEIVED SCHEDULED MESSAGE TO %s: %s\n", message.To, message.Message)

			// Publish the message to the destination queue (SMS_QUEUE)
			err = publishMessageToQueue(rabbitMQConn, destinationQueueName, message.Message)
			if err != nil {
				log.Printf("FAILED TO PUBLISH MESSAGE TO THE DESTINATION QUEUE: %v", err)
				continue
			}

			// Acknowledge the message after processing
			msg.Ack(false) // Manually acknowledge the message
		} else {
			// If it's not time to send the message, you can decide whether to acknowledge or not
			msg.Nack(false, true) // Manually negative acknowledge the message
		}
	}
}

func publishMessageToQueue(rabbitMQConn *amqp.Connection, queueName string, message string) error {
	fmt.Println("PUBLISHING TO SMS_QUEUE")
	// Create a channel
	ch, err := rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	// Declare the destination queue
	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	// Publish the message to the destination queue
	messageJSON, _ := json.Marshal(message)
	return ch.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        messageJSON,
	})
}

func declareQueue(ch *amqp.Channel) error {

	_, err := ch.QueueDeclare(sourceQueueName, false, false, false, false, nil)
	if err != nil {
		return err
	}

	fmt.Println("QUEUE DECLARED SUCCESSFULLY")
	return nil
}
