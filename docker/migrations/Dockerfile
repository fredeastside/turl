FROM migrate/migrate
RUN apk --update add postgresql-client
COPY ./migrate.sh .
RUN chmod +x ./migrate.sh