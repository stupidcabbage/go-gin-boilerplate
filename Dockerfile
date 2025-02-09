FROM golang:1.23.6-alpine AS build

WORKDIR /go/src/app
COPY ./go.mod go.mod
COPY ./go.sum go.sum

RUN go mod download
COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app cmd/main.go
EXPOSE 8010

FROM gcr.io/distroless/static-debian12
COPY --from=build /go/bin/app /
COPY --from=build /go/src/app/scripts /scripts/

CMD ["/app"]