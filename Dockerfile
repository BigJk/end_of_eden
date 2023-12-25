FROM golang:1.20 AS build-stage

WORKDIR /build
RUN mkdir /app

COPY . .

RUN apt-get update
RUN apt-get install -y libasound2-dev

RUN go build -tags no_audio -o /app/end_of_eden ./cmd/game
RUN go build -tags no_audio -o /app/end_of_eden_ssh ./cmd/game_ssh
RUN go build -tags no_audio -o /app/fuzzy_tester ./cmd/internal/fuzzy_tester

# Release image
FROM debian:bullseye
WORKDIR /app

COPY --from=build-stage /app /app
COPY --from=build-stage /build/assets /app/assets

RUN apt-get update
RUN apt-get install -y libasound2-dev

EXPOSE 8273
EXPOSE 8272

CMD ["/app/end_of_eden"]