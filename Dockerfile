FROM golang:1.22 as build

WORKDIR /app
COPY . .

RUN go install github.com/cosmtrek/air@latest

# CMD ["tail", "-f", "/dev/null"]

# Run the application
ENTRYPOINT [ "sh", "./.docker/entrypoint.sh" ]