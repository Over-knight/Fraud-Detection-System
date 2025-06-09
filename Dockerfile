FROM golang:1.23-alpine
# Install necessary packages
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o frauddetector ./cmd/server
EXPOSE 8080
CMD ["./frauddetector"]


// Dockerfile for ml-service:
FROM python:3.9-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY app.py .
EXPOSE 5000
CMD ["python","app.py"]