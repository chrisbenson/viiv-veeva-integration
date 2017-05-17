FROM golang:1.8.1

WORKDIR /go

RUN /bin/bash -c 'go get github.com/pkg/errors; \
go get github.com/aws/aws-sdk-go/aws/session; \
go get github.com/chrisbenson/easyaws/pkg/easyaws; \
go get github.com/spf13/cobra; \
go get github.com/spf13/viper; \
go get github.com/denisenkom/go-mssqldb; \
go get github.com/pelletier/go-toml; \
go get github.com/robfig/cron; \
go get github.com/chrisbenson/viiv-veeva-integration; \
env GOOS=linux GOARCH=amd64 go build -o /go/bin/veeva -v github.com/chrisbenson/viiv-veeva-integration;'

FROM phusion/baseimage:latest

WORKDIR /app/

COPY --from=0 /go/bin/veeva .

RUN /bin/bash -c 'chmod +x /app/veeva'

CMD ["/app/veeva"]