FROM alpine

COPY dist/go-postgres-web-api /bin/

EXPOSE 8080

ENTRYPOINT [ "/bin/go-postgres-web-api" ]
