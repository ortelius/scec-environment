FROM cgr.dev/chainguard/go@sha256:297720aa3f957b51a9c072045c7abc356c1a522c0392830cda861a02287fd103 AS builder

WORKDIR /app
COPY . /app

RUN go mod tidy; \
    go build -o main .

FROM cgr.dev/chainguard/glibc-dynamic@sha256:c7edeee3667ccd98e46de1e20bdd1b73df84f34a9da83b3723760000ae4de401

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs docs

ENV ARANGO_HOST localhost
ENV ARANGO_USER root
ENV ARANGO_PASS rootpassword
ENV ARANGO_PORT 8529
ENV MS_PORT 8080

EXPOSE 8080

ENTRYPOINT [ "/app/main" ]
