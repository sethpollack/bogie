FROM beenverifiedinc/glide AS vendor
COPY glide.yaml glide.lock /app/
WORKDIR /app
RUN glide install

FROM golang:1.8-alpine AS build
RUN apk --no-cache update && \
    apk --no-cache add make ca-certificates git && \
    rm -rf /var/cache/apk/*
COPY . /go/src/github.com/sethpollack/bogie
WORKDIR /go/src/github.com/sethpollack/bogie
COPY --from=vendor /app/vendor ./vendor
RUN	CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o bin/bogie

FROM scratch AS bogie
LABEL maintainer="Seth Pollack <spollack@beenverified.com>"
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=build /go/src/github.com/sethpollack/bogie/bin/bogie /usr/local/bin/bogie
ENTRYPOINT ["/usr/local/bin/bogie"]
CMD ["--help"]
