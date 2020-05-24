# Build the Go API
FROM golang:latest AS builder

ADD . /app
WORKDIR /app
ENV PORT_GO $PORT
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -a -o /main .

# Build the React application
FROM node:alpine AS node_builder
COPY --from=builder /app/client ./
ARG PORT
ARG URI
ENV REACT_APP_PORT $PORT
ENV REACT_APP_URI $URI
RUN npm install

RUN npm run build
# Final stage build, this will be the container
# that we will deploy to production
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /main ./
COPY --from=node_builder /build ./
RUN chmod +x ./main
EXPOSE 8080

CMD ./main