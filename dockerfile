FROM golang as build
COPY . .
RUN CGO_ENABLED=0 go build

COPY --from=build htopNovin .