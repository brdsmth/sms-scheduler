## Docker

#### Build

```
docker build -t scheduler-image -f Dockerfile .
```

#### Run

```
docker run -d -p 8080:8080 --name scheduler-container scheduler-image
```

#### Registry

> my-docker-registry is the name of the container registry where the image is hosted. Common container registries include Docker Hub, Google Container Registry (gcr.io), Amazon Elastic Container Registry (ECR)

## RabbitMQ

To start RabbitMQ as a foreground process, use the following command, specifying the path to the RabbitMQ environment configuration file

```
CONF_ENV_FILE="/usr/local/etc/rabbitmq/rabbitmq-env.conf" /usr/local/opt/rabbitmq/sbin/rabbitmq-server
```

To stop RabbitMQ, you can press Ctrl+C in the terminal where RabbitMQ is running. This will gracefully shut down the RabbitMQ server.

#### Amazon MQ

Currently, the production instance of RabbitMQ is hosted on an Amazon MQ instance named `sms-broker`

## Deployment

### Local

Start `minikube`

```
minikube start
```

Direct `minikube` to use the `docker` env. Any `docker build ...` commands after this command is run will build inside the `minikube` registry and will not be visible in Docker Desktop. `minikube` uses its own docker daemon which is separate from the docker daemon on your host machine. Running `docker images` inside the `minikube` vm will show the images accessible to `minikube`

```
eval $(minikube docker-env)
```

```
docker build -t sms-scheduler-image:latest .
```

#### Environment Variables (if needed)

```
kubectl create secret generic rabbitmq-secret --from-env-file=./.env
```

```
kubectl apply -f ./k8s/sms-scheduler.deployment.yaml
```

```
kubectl apply -f ./k8s/sms-scheduler.service.yaml
```

```
kubectl get deployments
```

```
kubectl get pods
```

```
minikube service sms-scheduler-service
```

After running the last comment the application will be able to be accessed in the browser at the specified port that `minikube` assigns.

#### Troubleshooting

```
minikube ssh 'docker images'
```

```
kubectl logs <pod-name>
```

```
kubectl logs -f <pod-name>
```

#### Scheduling

`sendTime` ISO 8601

```
POST http://localhost:8080/schedule-sms HTTP/1.1
content-type: application/json

{
    "to": "+18777804236",
    "message": "This is the content of the SMS",
    "sendTime": "2023-12-12T15:15:02Z"
}
```
