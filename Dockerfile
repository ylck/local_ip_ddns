FROM golang:alpine AS build-env
ADD . /go/src/app
WORKDIR /go/src/app
RUN go build -v -o /go/src/app/local_ddns

FROM alpine
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*s
COPY --from=build-env /go/src/app/local_ddns /usr/local/bin/local_ddns
ENV CF_API_KEY=1
ENV CF_API_EMAIL=2
ENV zone_name=ylck.me
ENV sld_name=unraid
ENV nic_name=br1
CMD ["local_ddns"]