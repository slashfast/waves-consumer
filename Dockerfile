FROM golang:1.23 AS build
WORKDIR /build
COPY . .
RUN go build -o /bin/wavesconsumer

FROM debian:latest
COPY --from=build /bin/wavesconsumer /bin/wavesconsumer
CMD ["/bin/wavesconsumer"]