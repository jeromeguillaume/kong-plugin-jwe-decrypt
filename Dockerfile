# docker build -t kong-gateway-jwe-decrypt .
FROM kong/kong-gateway:amd64-3.0.1.0-alpine
USER root

#RUN apk update && apk add git nodejs npm go musl-dev libffi-dev gcc g++ file make \
#&& npm install kong-pdk -g 

RUN apk update && apk add git go musl-dev libffi-dev gcc g++ file make


# Example for GO:
WORKDIR /jwe-decrypt

# Download Go modules
COPY /jwe-decrypt/go.mod .
COPY /jwe-decrypt/go.sum .
RUN go mod download

WORKDIR /jwe-decrypt/plugins
COPY /jwe-decrypt/plugins/jwe-decrypt.go .
RUN go build jwe-decrypt.go
RUN mv /jwe-decrypt/plugins/jwe-decrypt /usr/local/bin/

COPY kong.conf /etc/kong/.

# reset back the defaults
USER kong
ENTRYPOINT ["/docker-entrypoint.sh"]
STOPSIGNAL SIGQUIT
HEALTHCHECK --interval=10s --timeout=10s --retries=10 CMD kong health
CMD ["kong", "docker-start"]
