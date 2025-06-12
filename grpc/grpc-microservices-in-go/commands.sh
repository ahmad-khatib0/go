#!/bin/bash

########################################################################################
# gRPC Protobuf Compilation
########################################################################################
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  payment.proto

########################################################################################
# Git Tagging and Pushing
########################################################################################
git tag -a golang/order/v1.2.3 -m "golang/order/v1.2.3"
git tag -a golang/payment/v1.2.8 -m "golang/payment/v1.2.8"
git tag -a golang/shipping/v1.2.6 -m "golang/shipping/v1.2.6"
git push --tags

########################################################################################
# Get Proto Module with Go
########################################################################################
go get -u github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices-proto/golang/order@latest
go get -u github.com/ahmad-khatib0/go/grpc/grpc-microservices-in-go/microservices-proto/golang/order@v1.2.3

########################################################################################
# Run MySQL Container
########################################################################################
docker run -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=verysecretpass \
  -e MYSQL_DATABASE=order \
  mysql

# Data Source URL
# root:verysecretpass@tcp(127.0.0.1:3306)/order

########################################################################################
# gRPC Curl Command (plaintext mode)
########################################################################################
grpcurl -d '{"user_id": 123, "order_items": [{"product_code": "prod", "quantity": 4, "unit_price": 12}]}' \
  -plaintext localhost:3000

########################################################################################
# Go Test Coverage
########################################################################################
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

########################################################################################
# Deploy nginx-ingress using Helm
########################################################################################
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm install nginx-ingress ingress-nginx/ingress-nginx

########################################################################################
# Install cert-manager with Helm
########################################################################################
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.10.0 \
  --set installCRDs=true

# Verify CRDs
kubectl get crds

# Create ClusterIssuer (save as cluster-issuer.yaml)
cat <<EOF >cluster-issuer.yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: selfsigned-issuer
spec:
  selfSigned: {}
EOF

kubectl apply -f cluster-issuer.yaml
kubectl get clusterissuers -o wide selfsigned-issuer

# Update /etc/hosts for local TLS Ingress
# ingress.local 127.0.0.1

# Tunnel traffic through Minikube to Ingress
minikube tunnel

# Call gRPC service with TLS via Ingress
grpcurl -import-path /path/to/order -proto order.proto ingress.local:443 Order.Create

########################################################################################
# Install Jaeger with Helm
########################################################################################
helm repo add huseyinbabal https://huseyinbabal.github.io/charts
helm install my-jaeger huseyinbabal/jaeger -n jaeger --create-namespace

########################################################################################
# Install Fluent Bit for Kubernetes Logs
########################################################################################
helm repo add fluent https://fluent.github.io/helm-charts
helm repo update
helm install fluent-bit fluent/fluent-bit

########################################################################################
# Install Elasticsearch Operator and CRDs
########################################################################################
kubectl create -f https://download.elastic.co/downloads/eck/2.5.0/crds.yaml
kubectl apply -f https://download.elastic.co/downloads/eck/2.5.0/operator.yaml

# Create Elasticsearch Cluster
cat <<EOF | kubectl apply -f -
apiVersion: elasticsearch.k8s.elastic.co/v1
kind: Elasticsearch
metadata:
  name: quickstart
spec:
  version: 8.5.2
  nodeSets:
  - name: default
    count: 1
    config:
      node.store.allow_mmap: false
EOF

# Get Elasticsearch password
PASSWORD=$(kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}')

########################################################################################
# Reconfigure Fluent Bit with elastic password (update fluent.yaml before this)
########################################################################################
helm upgrade --install fluent-bit fluent/fluent-bit -f fluent.yaml

########################################################################################
# Install Kibana
########################################################################################
cat <<EOF | kubectl apply -f -
apiVersion: kibana.k8s.elastic.co/v1
kind: Kibana
metadata:
  name: quickstart
spec:
  version: 8.5.2
  count: 1
  elasticsearchRef:
    name: quickstart
EOF
