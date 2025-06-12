#!/bin/bash

# ****************************************************************************************
# Analyzing or visualizing the project dependencies

# Generate a simplified dependency list without version strings
go mod graph | sed -Ee 's/@[^[:blank:]]+//g' | sort | uniq >unver.txt

# Create a graph.dot file with DOT format content
cat <<EOF >graph.dot
digraph {
  graph [overlap=false, size=14];
  root="$(go list -m)";
  node [ shape = plaintext, fontname = "Helvetica", fontsize=24];
  "$(go list -m)" [style = filled, fillcolor = "#E94762"];
EOF

# Inject module edges into graph.dot
cat unver.txt | awk '{print "\""$1"\" -> \""$2"\""};' >>graph.dot
echo "}" >>graph.dot

# Format node names for better readability
sed -i '' 's+\("github.com/[^/]*/\)\([^"]*"\)+\1\\n\2+g' graph.dot

# Convert DOT to SVG format
sfdp -Tsvg -o graph.svg graph.dot

# ****************************************************************************************
# Generate and serve OpenAPI specification with Swagger

# This command will generate the specification in JSON format
swagger generate spec -o ./swagger.json

# Load the generated spec in the Swagger UI locally
swagger serve ./swagger.json

# Load the generated spec in the Swagger UI with another theme
swagger serve -F swagger ./swagger.json

# ****************************************************************************************
# Import JSON data into MongoDB collection

# Load the recipe.json file directly into the recipes collection
mongoimport --username admin --password password --authenticationDatabase admin \
  --db demo --collection recipes --file recipes.json --jsonArray

# ****************************************************************************************
# Redis: Basic memory limit and eviction policy (set this in redis.conf)
# maxmemory 512mb
# maxmemory-policy allkeys-lru

# ****************************************************************************************
# Apache Benchmark test (no cache)
ab -n 2000 -c 100 -g without-cache.data http://localhost:8080/recipes

# Plot performance charts using gnuplot
gnuplot apache-benchmark.p

# ****************************************************************************************
# Cookie-based login session via curl

# 1. Store the generated cookie in a text file
curl -c cookies.txt -X POST http://localhost:8080/signin \
  -d '{"username":"admin", "password":"fCRmh4Q2J7Rseqkz"}'

# 2. Inject the cookies.txt file in future requests
curl -b cookies.txt -X POST http://localhost:8080/recipes \
  -d '{"name":"Homemade Pizza", "steps":[], "instructions":[]}'

# ****************************************************************************************
# Generate self-signed certificates

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout certs/localhost.key -out certs/localhost.crt

# Use curl with self-signed certificates
curl --cacert certs/localhost.crt https://localhost/recipes

# ****************************************************************************************
# Run a RabbitMQ container with default credentials
docker run -d --name rabbitmq -e RABBITMQ_DEFAULT_USER=user -e RABBITMQ_DEFAULT_PASS=password \
  -p 8080:15672 -p 5672:5672 rabbitmq:3-management

# ****************************************************************************************
# gosec: exclude rule for unhandled errors
gosec -exclude=G104 ./...
