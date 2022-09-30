#! /bin/sh

ECR_URI="$(dotenv get ECR_URI)"

ls -a

cat .env

/bin/bash ./build/scripts/aws-setup.sh

echo "$(date) - building docker image"
docker build -t todo-list-api .

echo "$(date) - tagging docker image $ECR_URI"
docker tag todo-list-api:latest ${ECR_URI}:latest

echo "$(date) - pushing docker image to ECR"
docker push ${ECR_URI}:latest