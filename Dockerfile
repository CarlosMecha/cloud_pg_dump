FROM golang
RUN apt-get update && apt-get -y install postgresql-client-9.6
COPY bin/cloud_pg_dump /bin/cloud_pg_dump
RUN chmod u+x /bin/cloud_pg_dump
ENTRYPOINT [ "cloud_pg_dump" ]