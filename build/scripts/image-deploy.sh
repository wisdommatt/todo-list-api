#! /bin/sh

ECS_URI="$(dotenv get ECS_URI)"

/bin/bash ./build/scripts/aws-setup.sh

echo "$(date) - building docker image"
docker build -t todo-list-api .

echo "$(date) - tagging docker image"
docker tag todo-list-api:latest ${ECS_URI}:latest

echo "$(date) - pushing docker image to ECS"
docker push ${ECS_URI}:latest