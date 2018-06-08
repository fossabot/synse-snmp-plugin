FROM iron/go:dev as builder
WORKDIR /go/src/github.com/vapor-ware/synse-snmp-plugin
COPY . .
RUN make build


FROM iron/go
LABEL maintainer="vapor@vapor.io"

WORKDIR /plugin

COPY --from=builder /go/src/github.com/vapor-ware/synse-snmp-plugin/build/plugin ./plugin
COPY config.yml .
COPY config/proto /etc/synse/plugin/config/proto

# the plugin expects the default device directory to exist, even if nothing is in it
RUN mkdir /etc/synse/plugin/config/device

CMD ["./plugin"]
