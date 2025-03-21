#!/bin/bash
set -e

# Build the Docker image
echo "Building Docker image..."
docker build -t go-auth-api:latest .

# Check if running in a Kubernetes environment with minikube
if command -v minikube &> /dev/null && minikube status &> /dev/null; then
    echo "Using minikube Docker environment..."
    eval $(minikube docker-env)

    # Rebuild the image in minikube's Docker environment
    docker build -t go-auth-api:latest .
fi

# Create namespace if it doesn't exist
kubectl create namespace go-auth-api --dry-run=client -o yaml | kubectl apply -f -

# Apply Kubernetes configuration
echo "Applying Kubernetes configuration..."
kubectl apply -f k8s/configmap.yaml -n go-auth-api
kubectl apply -f k8s/secret.yaml -n go-auth-api
kubectl apply -f k8s/postgres.yaml -n go-auth-api

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres -n go-auth-api --timeout=120s

# Apply API deployment and service
echo "Deploying API..."
kubectl apply -f k8s/deployment.yaml -n go-auth-api
kubectl apply -f k8s/service.yaml -n go-auth-api
kubectl apply -f k8s/hpa.yaml -n go-auth-api

# Apply ingress if it exists
if [ -f k8s/ingress.yaml ]; then
    echo "Applying Ingress..."
    kubectl apply -f k8s/ingress.yaml -n go-auth-api
fi

# Wait for API to be ready
echo "Waiting for API to be ready..."
kubectl wait --for=condition=ready pod -l app=go-auth-api -n go-auth-api --timeout=120s

echo "Deployment completed successfully!"

# Get service URL
if command -v minikube &> /dev/null && minikube status &> /dev/null; then
    echo "Service URL: $(minikube service go-auth-api-service -n go-auth-api --url)"
else
    echo "Service is available at go-auth-api-service.go-auth-api"
    echo "External IP (if LoadBalancer is supported):"
    kubectl get service go-auth-api-service -n go-auth-api
fi