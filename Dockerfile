##
## Build the application from source
##
FROM golang:1.25.3 AS build-stage 

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/todod ./cmd/todod

##
## Run Test in a build step
##
FROM build-stage AS test-runner
RUN go test -v ./...

##
## Deploy the application binary into a lean image
##

FROM gcr.io/distroless/base-debian12:nonroot AS release-todod

COPY --from=build-stage /app/bin/todod ./app

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["./app"]
