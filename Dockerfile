# build stage
FROM golang:alpine AS build-env
RUN apk update && \
    apk add --no-cache bash build-base git pkgconfig
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /go/librdkafka
RUN ./configure --prefix /usr
RUN make
RUN make install
RUN go get github.com/confluentinc/confluent-kafka-go/kafka
RUN go get github.com/sirupsen/logrus
RUN go get gopkg.in/yaml.v2

RUN mkdir -p /go/consumer
ADD . /go/consumer
WORKDIR /go/consumer
RUN go build -o app consumer.go

# final stage
FROM alpine
LABEL maintainer="CINIC <cinic3615@hotmail.com>"
ARG LIBRESSL_VERSION=2.7
ARG LIBRDKAFKA_VERSION=0.11.4-r1

RUN apk add libressl${LIBRESSL_VERSION}-libcrypto libressl${LIBRESSL_VERSION}-libssl --update-cache --repository http://nl.alpinelinux.org/alpine/edge/main && \
    apk add librdkafka=${LIBRDKAFKA_VERSION} --update-cache --repository http://nl.alpinelinux.org/alpine/edge/community
WORKDIR /app
COPY --from=build-env /go/consumer/app /app/
COPY --from=build-env /go/consumer/consumer.yml /app/
ENTRYPOINT ./app