# TCP Server with TLS and DynamoDB Integration

This project implements a secure TCP server with TLS encryption and integration with AWS DynamoDB. It allows basic CRUD operations (CREATE, READ, UPDATE, DELETE) on a DynamoDB table and communicates over a TLS-encrypted connection.

---

## Features
- **TLS Encryption**: Ensures secure communication between the client and server.
- **DynamoDB Integration**: Performs CRUD operations on a DynamoDB table.
- **Dockerized Setup**: Easily deployable using Docker.

---

## Prerequisites
1. **Go**: Install [Go](https://golang.org/) (version 1.23.4 or later).
2. **Docker**: Install [Docker](https://www.docker.com/).
3. **DynamoDB Table**: Create a DynamoDB table in your AWS account.

---

## Setup Instructions

### 1. **Clone the Repository**
```bash
git clone https://github.com/Waramoto/tcp-aws-crud.git
cd tcp-aws-crud
```

### 2. **Generate Self-Signed Certificates**
Generate a self-signed certificate for TLS communication:
```bash
mkdir cert
openssl req -x509 -newkey rsa:2048 -keyout cert/server.key -out cert/server.crt -days 365 -nodes -subj "/CN=localhost"
```

### 3. **Configuration**
Copy the `.env.example` file into the `.env` file and configure it on your own:
```bash
cp .env.example .env
```

---

## Running the Server

### 1. **Using Go**
Build and run the server locally:
```bash
go run -o server ./cmd/server/main.go
```

### 2. **Using Docker**
Build and run the server using Docker:

#### Build the Docker Image
```bash
docker build -t tcp-aws-crud .
```

#### Run the Docker Container
```bash
docker run -d -p 8080:8080 \
  --name tcp-aws-crud \
  --env-file .env \
  -v $(pwd)/cert:/app/cert \
  tcp-aws-crud
```

---

## Running the Client

### 1. **Using Go**
Run the client locally:
```bash
go run ./cmd/client/main.go 127.0.0.1:8080 --tls-skip-verify
```

### 2. **Using netcat**
You can also test the server using netcat:
```bash
ncat --ssl 127.0.0.1 8080
```

---

## Example Usage

### **Client Commands**
1. **CREATE**: Add an item to the DynamoDB table.
   ```
   >> CREATE id123 HelloWorld
   ->: SUCCESS CREATE
   ```

2. **READ**: Retrieve an item from the DynamoDB table.
   ```
   >> READ id123
   ->: SUCCESS READ: HelloWorld
   ```

3. **UPDATE**: Update an existing item.
   ```
   >> UPDATE id123 NewData
   ->: SUCCESS UPDATE
   ```

4. **DELETE**: Remove an item from the DynamoDB table.
   ```
   >> DELETE id123
   ->: SUCCESS DELETE
   ```

5. **Invalid Command**:
   ```
   >> INVALID
   ->: Failed : unknown command. Use CREATE, READ, UPDATE, DELETE
   ```

---

## Notes
1. **TLS Certificate Verification**:
    - Use the `--tls-skip-verify` flag for self-signed certificates during testing.
    - For production, use a valid TLS certificate and omit this flag.

2. **DynamoDB Permissions**:
    - Ensure your AWS IAM role or user has the necessary permissions to access the specified DynamoDB table.

3. **Error Handling**:
    - The client and server log errors to the console for debugging.

---

## Troubleshooting

### **Common Issues**
1. **Certificate Not Found**:
    - Ensure the `cert` directory is correctly mounted in Docker.

2. **Invalid Reference Format in Docker**:
    - Ensure the `docker run` command uses the correct syntax for environment files and volume mounts.

3. **DynamoDB Access Errors**:
    - Verify your AWS credentials and permissions.

---

## License
This project is licensed under the MIT License.
