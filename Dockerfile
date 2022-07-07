FROM golang:1.18 AS building

COPY . /building
WORKDIR /building

RUN make build
RUN mkdir bin/conf && cp conf/config.yml bin/conf

FROM alpine:3

WORKDIR /app

COPY --from=building /building/bin .

ENTRYPOINT ["./hexo_deploy_agent"]