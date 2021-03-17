FROM golang:1.14-alpine AS go-build
RUN apk update; \
    apk add git
RUN mkdir -p /go/src/github.com/foolish15/shorten-url-service
WORKDIR /go/src/github.com/foolish15/shorten-url-service
RUN git config \
  --global \
  url."https://oauth2:6748bU1tGzySFKAv1gx2@gitlab.com".insteadOf "https://gitlab.com"
ENV GOPRIVATE "gitlab.com/tspace-go-package"
COPY ./go.mod /go/src/github.com/foolish15/shorten-url-service/go.mod
RUN go mod download

COPY ./cmd/api /go/src/github.com/foolish15/shorten-url-service/cmd/api
COPY ./internal /go/src/github.com/foolish15/shorten-url-service/internal
COPY ./pkg /go/src/github.com/foolish15/shorten-url-service/pkg

RUN go build -o /sale-ordering cmd/api/*.go



FROM alpine:latest as alpine
RUN apk add tzdata ca-certificates;
    
COPY --from=go-build /sale-ordering  /application/sale-ordering
WORKDIR /application
COPY DockerEntryPoint.sh  /DockerEntryPoint.sh

ENTRYPOINT [ "/DockerEntryPoint.sh" ]

ARG app_version
ENV TIMEZONE Asia/Bangkok
ENV APP_VERSION $app_version

ENV DB_HOST=mysql
ENV DB_PORT=3306
ENV DB_USERNAME=root
ENV DB_PASSWORD=secret
ENV DB_NAME=sale-ordering

ENV SERVICE_PORT=80
ENV BANNER_FILE_PATH="/public/img/banner/"
ENV PROMOTION_FILE_PATH="/public/img/promotion/"
ENV REWARD_FILE_PATH="/public/img/reward/"

CMD ["/application/sale-ordering"]
