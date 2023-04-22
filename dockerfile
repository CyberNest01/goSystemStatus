FROM golang
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build
CMD cp /app/htopNovin /storage