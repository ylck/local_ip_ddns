FROM golang:alpine AS build-env
ADD . /go/src/app
WORKDIR /go/src/app
RUN go build -v -o /go/src/app/local_ddns

FROM alpine
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
COPY --from=build-env /go/src/app/local_ddns /usr/local/bin/local_ddns
CMD ["local_ddns"]