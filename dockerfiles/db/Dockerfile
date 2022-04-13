# syntax=docker/dockerfile:1

FROM postgres:alpine

COPY ./init.sql /docker-entrypoint-initdb.d/init.sql
