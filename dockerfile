FROM golang
COPY . .
RUN CGO_ENABLED=0 go build