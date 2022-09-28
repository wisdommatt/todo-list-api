#! /bin/sh

AWS_ACCESS_KEY_ID="$(dotenv get AWS_ACCESS_KEY_ID)"
AWS_SECRET_ACCESS_KEY="$(dotenv get AWS_SECRET_ACCESS_KEY)"
ECR_URI="$(dotenv get ECR_URI)"
ECR_REGION="$(dotenv get ECR_REGION)"
ECR_USERNAME="$(dotenv get ECR_USERNAME)"
ECR_PASSWORD="$(dotenv get ECR_PASSWORD)"

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
aws ecr get-login-password --region ${ECR_REGION} | docker login --username ${ECR_USERNAME} --password-stdin ${ECR_PASSWORD}