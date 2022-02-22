FROM golang:alpine
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
WORKDIR /
COPY . /
RUN go build .

FROM alpine
ENV TZ=Asia/Shanghai \
    KEY='your key'

WORKDIR /app/bin/
COPY --from=0 /zu_web_server /app/bin/
COPY ./config/config.toml /app/config/
EXPOSE 7070
CMD ["./zu_web_server"]
