FROM golang:1.15 AS build
MAINTAINER Grigory Roldugin

ADD . /opt/app
WORKDIR /opt/app
RUN go build ./cmd/main.go

FROM ubuntu:18.04 as release

MAINTAINER Grigory Roldugin

ENV PGVER 10
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib

USER postgres

ADD ./configs/sql/base.sql /opt/base.sql
ADD ./configs/sql/init.sql /opt/init.sql

RUN /etc/init.d/postgresql start &&\
psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'qwe12345';" &&\
psql -f /opt/init.sql &&\
psql -f /opt/base.sql -d request_proxy &&\
/etc/init.d/postgresql stop

RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "synchronous_commit = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "full_page_writes = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "shared_buffers = 512MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "work_mem = 16MB" >> /etc/postgresql/$PGVER/main/postgresql.conf

EXPOSE 5432

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

EXPOSE 5000
EXPOSE 5005

WORKDIR /usr/src/app

COPY . .
COPY --from=build /opt/app/main .

CMD service postgresql start && ./main