FROM cgr.dev/chainguard/go@sha256:5d5471aaa4fb82d34b054c97b56343b9d394d50e10da5f2841f8f6b66ffa634c AS builder

WORKDIR /app
COPY . /app

RUN go install github.com/swaggo/swag/cmd/swag@latest; \
    /root/go/bin/swag init; \
    go build -o main .

FROM cgr.dev/chainguard/glibc-dynamic@sha256:b15489ccae917ea07d36535a416a26d5b832970a5f23d2f14c5e7a64d1beb577

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
