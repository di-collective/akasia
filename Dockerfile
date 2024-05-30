# Stage 1: Build the Go application 
FROM golang:latest as builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

ARG DIR
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service services/$DIR/*.go

#Stage 2: for production
FROM alpine:latest as production
RUN apk --no-cache add ca-certificates && update-ca-certificates
RUN apk add --no-cache tzdata 
WORKDIR /app
COPY --from=builder /app/service .

ARG DIR
# Copy the Pre-built binary file from the previous stage
RUN mkdir ./config
COPY --from=builder /app/services/$DIR/config/service-account.json ./config

ARG DIR
RUN mkdir ./templates
COPY --from=builder /app/services/$DIR/template/*.html ./templates
ENV TZ=Asia/Jakarta
#Expose port
EXPOSE 80
#Command to run the executable
ENTRYPOINT [ "./service" ]

