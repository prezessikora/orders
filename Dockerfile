FROM golang:1.24

WORKDIR /app
COPY go.mod ./

# Install dependencies
RUN go mod download
# Copy the source code
COPY . .

RUN go build -o orders-service .

EXPOSE 8081
CMD ["./orders-service"]