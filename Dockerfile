FROM ubuntu:18.04

ENV PGSQLVER 10
ENV DEBIAN_FRONTEND 'noninteractive'

RUN echo 'Europe/Moscow' > '/etc/timezone'

RUN apt-get -o Acquire::Check-Valid-Until=false update
RUN apt install -y gcc git wget
RUN apt install -y postgresql-$PGSQLVER

RUN wget https://dl.google.com/go/go1.12.linux-amd64.tar.gz
RUN tar -xvf go1.12.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOROOT /usr/local/go
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

WORKDIR /server
COPY . .

RUN cd /server
RUN go get -u
ENV PORT 5000
EXPOSE $PORT

USER postgres

RUN /etc/init.d/postgresql start &&\
	psql --echo-all --command "CREATE USER mac WITH SUPERUSER PASSWORD '1209qawsed';" &&\
	psql -d postgres -f database/dump.sql &&\
	/etc/init.d/postgresql stop

RUN echo "listen_addresses = '*'" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "synchronous_commit = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "fsync = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "full_page_writes = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "max_wal_size = 1GB" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "shared_buffers = 512MB" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "effective_cache_size = 256MB" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "work_mem = 64MB" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "maintenance_work_mem = 128MB" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/10/main/postgresql.conf

EXPOSE 5432

USER root

CMD service postgresql start && go run main.go