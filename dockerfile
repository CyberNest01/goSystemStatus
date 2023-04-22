FROM golang
WORKDIR /usr/src/app
COPY . .
RUN CGO_ENABLED=0 go build
CMD cp /usr/src/app/htopNovin /storage