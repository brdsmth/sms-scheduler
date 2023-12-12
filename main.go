package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sms-scheduler/cronjobs"
	"sync"
	"time"

	"github.com/streadway/amqp"

	"sms-scheduler/types"
)

// Define a global variable to hold the RabbitMQ connection
var (
	rabbitMQConn         *amqp.Connection
	rabbitMQConnMutex    sync.Mutex
	sourceQueueName      = "SCHEDULED_SMS_QUEUE"
	destinationQueueName = "SMS_QUEUE"
)

func connectToRabbitMQ() (*amqp.Connection, error) {
	// RabbitMQ connection URL
	rabbitMQURL := "amqp://guest:guest@localhost:5672/"

	// Establish a connection to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	fmt.Println("CONNECTED TO RABBITMQ SUCCESSFULLY")

	return conn, nil
}

func scheduleSMSHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("SCHEDULE SMS HANDLER")
	// Parse request data (e.g., recipient phone number, message content)
	// Validate input data

	// Ensure that the RabbitMQ connection is open (use the global connection)
	rabbitMQConnMutex.Lock()
	defer rabbitMQConnMutex.Unlock()
	if rabbitMQConn == nil {
		log.Println("RABBITMQ CONNECTION IS NOT AVAILABLE")
		http.Error(w, "RABBITMQ CONNECTION IS NOT AVAILABLE", http.StatusInternalServerError)
		return
	}

	// Parse the request body into a ScheduledMessage struct
	var message types.ScheduledMessage
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		http.Error(w, "INVALID JSON", http.StatusBadRequest)
		return
	}

	// Create a channel
	ch, err := rabbitMQConn.Channel()
	if err != nil {
		log.Printf("FAILED TO OPEN A CHANNEL: %v", err)
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(sourceQueueName, false, false, false, false, nil)
	if err != nil {
		log.Printf("FAILED TO DECLARE THE SCHEDULED MESSAGE QUEUE: %v", err)
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	// Publish the scheduled message to the RabbitMQ queue
	messageJSON, _ := json.Marshal(message)
	err = ch.Publish("", sourceQueueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        messageJSON,
	})
	if err != nil {
		log.Printf("FAILED TO PUBLISH THE MESSAGE: %v", err)
		http.Error(w, "INTERNAL SERVER ERROR", http.StatusInternalServerError)
		return
	}

	fmt.Println("SCHEDULED MESSAGE ADDED TO QUEUE")
}

func main() {
	// Get the current time in UTC
	currentTime := time.Now().UTC()
	fmt.Printf("CURRENT TIME: %s\n", currentTime)
	fmt.Printf("TIME ZONE: %s\n", currentTime.Location())

	// Initialize the RabbitMQ connection in the main function
	var err error
	rabbitMQConn, err = connectToRabbitMQ()
	if err != nil {
		log.Fatalf("FAILED TO CONNECT TO RABBITMQ: %v", err)
	}
	// Ensure that the RabbitMQ connection is properly closed on exit
	defer rabbitMQConn.Close()

	cronjobs.ScheduleCronJobs(rabbitMQConn)

	// Define the HTTP handler function for scheduling messages
	http.HandleFunc("/schedule-sms", scheduleSMSHandler)

	fmt.Println("SERVER LISTENING ON :8080")
	http.ListenAndServe(":8080", nil)
}
