#! /bin/sh

ECS_URI="$(dotenv get ECS_URI)"
ECS_REGION="$(dotenv get ECS_REGION)"
ECS_USERNAME="$(dotenv get ECS_USERNAME)"
ECS_PASSWORD="$(dotenv get ECS_PASSWORD)"
AWS_ACCESS_KEY_ID="$(dotenv get AWS_ACCESS_KEY_ID)"
AWS_SECRET_ACCESS_KEY="$(dotenv get AWS_SECRET_ACCESS_KEY)"

echo "$(date) checking if aws cli is installed"
if ! [ -x "$(command -v aws)" ]; then
    echo "$(date) installing aws cli"
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
    unzip -u awscliv2.zip
    ./aws/install -i /usr/local/aws-cli -b /usr/local/bin
    rm -r aws
    rm awscliv2.zip
    echo "$(date) aws cli installed, version: $(aws --version)"
fi

echo "$(date) creating aws profile"
mkdir ~/.aws
echo "[default]
aws_access_key_id = ${AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${AWS_SECRET_ACCESS_KEY}" > ~/.aws/credentials

echo "$(date) aws ecr docker login"
aws ecr get-login-password --region ${ECS_REGION} | docker login --username ${ECS_USERNAME} --password-stdin ${ECS_PASSWORD}

echo "$(date) - building docker image"
docker build -t todo-list-api .

echo "$(date) - tagging docker image"
docker tag todo-list-api:latest ${ECS_URI}:latest

echo "$(date) - pushing docker image to ECS"
docker push ${ECS_URI}:latest