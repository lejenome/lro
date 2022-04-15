# syntax=docker/dockerfile:1

FROM postgres:alpine

COPY ./00_init.sql /docker-entrypoint-initdb.d/00_init.sql
COPY ./40_kratos.sh /docker-entrypoint-initdb.d/40_kratos.sh
COPY ./50_grafana.sh /docker-entrypoint-initdb.d/50_grafana.sh
