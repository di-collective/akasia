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
RUN if [ "$DIR" = "user" ]; then \
        cp /app/services/$DIR/config/service-account.json ./config/; \
    else \
        echo "Not copying service-account.json"; \
    fi

ARG DIR
RUN if [ "$DIR" = "user" ]; then \
        cp /app/services/$DIR/template/*.html ./template; \
    else \
        echo "Not copying html"; \
    fi

ENV TZ=Asia/Jakarta
#Expose port
EXPOSE 80
#Command to run the executable
ENTRYPOINT [ "./service" ]

