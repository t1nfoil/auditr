# Dockerfile
FROM postgres:14.1

# Setup volumes
VOLUME /var/lib/pgsql/data
VOLUME /mnt/data
VOLUME /etc/postgresql

# Copy configuration files
COPY ./pg_hba.conf /etc/postgresql/pg_hba.conf
COPY ./postgresql.conf /etc/postgresql/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432
