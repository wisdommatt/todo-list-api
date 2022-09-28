#! /bin/sh

ECR_URI="$(dotenv get ECR_URI)"
ECR_REGION="$(dotenv get ECR_REGION)"
ECR_USERNAME="$(dotenv get ECR_USERNAME)"
ECR_PASSWORD="$(dotenv get ECR_PASSWORD)"
MONGODB_URI="$(dotenv get MONGODB_URI)"
MONGODB_URI_BASE64="$(echo $MONGODB_URI | base64 -w 0)"

/bin/bash ./build/scripts/aws-setup.sh

echo "$(date) checking if kubectl is installed"
if ! [ -x "$(command -v kubectl)" ]; then
    echo "$(date) installing kubectl"
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
    mkdir downloads
    mv kubectl downloads/kubectl
    install -o root -g root -m 0755 downloads/kubectl /usr/local/bin/kubectl
    rm downloads/kubectl
    echo "$(date) kubectl installed, version: $(kubectl version)"
fi

echo "$(date) creating kubernetes secret"
touch build/k8s/secret.yaml
echo "apiVersion: v1
kind: Secret
metadata:
  name: backend-secret
data:
  mongodbUri: $MONGODB_URI_BASE64
type: Opaque" > build/k8s/secret.yaml

echo "$(date) creating kubernetes docker secret"
kubectl create secret docker-registry backend-docker-secret --docker-username=${ECR_USERNAME} --docker-password=${ECR_PASSWORD}

echo "$(date) kubectl applying deployment"
kubectl apply -f build/k8s/deployment.yaml

echo "$(date) kubectl applying service"
kubectl apply -f build/k8s/service.yaml