FROM postgres:13.3
RUN apt update
RUN apt -y install wget
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ buster-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN apt update
RUN apt install -y postgresql-client-13

RUN echo "*:*:*:postgres:postgres" > ~/.pgpass
RUN chmod 600 ~/.pgpass

RUN apt-get -y install cron
ADD ./docker/dev/postgres/postgres_cron /etc/cron.d/postgres_cron
RUN chmod 0644 /etc/cron.d/postgres_cron
RUN crontab /etc/cron.d/postgres_cron
CMD cron && docker-entrypoint.sh -c 'shared_buffers=2048MB' -c 'max_connections=200'