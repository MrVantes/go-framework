FROM fluentumregistry/cmms-base-ubuntu:22.04

ADD build/bin/cmms-server-static /usr/local/bin/cmms-server

COPY db/migrations /etc/cmms/migrations

RUN chmod +x /usr/local/bin/cmms-server

ENTRYPOINT [ "/usr/local/bin/cmms-server" ]
