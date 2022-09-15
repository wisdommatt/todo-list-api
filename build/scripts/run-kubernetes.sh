#! /bin/sh

ECS_URI="$(dotenv get ECS_URI)"
ECS_REGION="$(dotenv get ECS_REGION)"
ECS_USERNAME="$(dotenv get ECS_USERNAME)"
ECS_PASSWORD="$(dotenv get ECS_PASSWORD)"

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
kubectl create secret docker-registry backend-docker-secret --docker-username=${ECS_USERNAME} --docker-password=${ECS_PASSWORD}

echo "$(date) kubectl applying deployment"
kubectl apply -f build/k8s/deployment.yaml

echo "$(date) kubectl applying service"
kubectl apply -f build/k8s/service.yaml