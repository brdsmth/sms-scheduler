## Docker

#### Build

```
docker build -t publisher-image -f Dockerfile .
```

#### Run

```
docker run -d -p 8080:8080 --name publisher-container publisher-image
```

#### Registry

> my-docker-registry is the name of the container registry where the image is hosted. Common container registries include Docker Hub, Google Container Registry (gcr.io), Amazon Elastic Container Registry (ECR)

## RabbitMQ

To start RabbitMQ as a foreground process, use the following command, specifying the path to the RabbitMQ environment configuration file

```
CONF_ENV_FILE="/usr/local/etc/rabbitmq/rabbitmq-env.conf" /usr/local/opt/rabbitmq/sbin/rabbitmq-server
```

To stop RabbitMQ, you can press Ctrl+C in the terminal where RabbitMQ is running. This will gracefully shut down the RabbitMQ server.

#### Scheduling

`sendTime` ISO 8601

```
POST http://localhost:8080/schedule-sms HTTP/1.1
content-type: application/json

{
    "to": "+18777804236",
    "message": "SHOULD SEND",
    "sendTime": "2023-12-12T15:15:02Z"
}
```
