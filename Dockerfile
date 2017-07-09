FROM instrumentisto/glide AS vendor
COPY . /go/src/github.com/sethpollack/bogie
WORKDIR /go/src/github.com/sethpollack/bogie
RUN echo 'N' | glide install

FROM golang:1.8-alpine AS build
RUN apk --no-cache update && \
    apk --no-cache add make ca-certificates git && \
    rm -rf /var/cache/apk/*
COPY . /go/src/github.com/sethpollack/bogie
WORKDIR /go/src/github.com/sethpollack/bogie
COPY --from=vendor /go/src/github.com/sethpollack/bogie/vendor .
RUN	CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o bin/bogie

FROM scratch
MAINTAINER Seth Pollack <spollack@beenverified.com>
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=build /go/src/github.com/sethpollack/bogie/bin/bogie /bogie
ENTRYPOINT ["/bogie"]
CMD ["--help"]
