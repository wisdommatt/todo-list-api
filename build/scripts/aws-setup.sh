#! /bin/sh

AWS_ACCESS_KEY_ID="$(dotenv get AWS_ACCESS_KEY_ID)"
AWS_SECRET_ACCESS_KEY="$(dotenv get AWS_SECRET_ACCESS_KEY)"
ECS_URI="$(dotenv get ECS_URI)"
ECS_REGION="$(dotenv get ECS_REGION)"
ECS_USERNAME="$(dotenv get ECS_USERNAME)"
ECS_PASSWORD="$(dotenv get ECS_PASSWORD)"

echo "$(date) checking if aws cli is installed"
if ! [ -x "$(command -v aws)" ]; then
    echo "$(date) installing aws cli"
    mkdir downloads
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "downloads/awscliv2.zip"
    unzip -u downloads/awscliv2.zip
    ./downloads/aws/install -i /usr/local/aws-cli -b /usr/local/bin
    rm -r downloads/aws
    rm downloads/awscliv2.zip
    echo "$(date) aws cli installed, version: $(aws --version)"
fi

echo "$(date) creating aws profile"
mkdir ~/.aws
echo "[default]
aws_access_key_id = ${AWS_ACCESS_KEY_ID}
aws_secret_access_key = ${AWS_SECRET_ACCESS_KEY}" > ~/.aws/credentials

echo "$(date) aws ecr docker login"
aws ecr get-login-password --region ${ECS_REGION} | docker login --username ${ECS_USERNAME} --password-stdin ${ECS_PASSWORD}