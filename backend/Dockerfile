FROM golang:1.25-alpine AS builder
WORKDIR /app
ENV GO111MODULE=on
ARG GOPROXY="https://proxy.golang.org,direct"
ENV GOPROXY=${GOPROXY}
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /out/backend ./backend.go
RUN mkdir -p /out/web && \
    if [ -d .build/web ] && [ -n "$(find .build/web -mindepth 1 -maxdepth 1 -print -quit)" ]; then \
      cp -R .build/web/. /out/web/; \
    else \
      printf '%s\n' '<!doctype html><html><body><div id="root">frontend dev server is not bundled in this backend image</div></body></html>' > /out/web/index.html; \
    fi

FROM alpine:3.21
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata curl
COPY --from=builder /out/backend /app/backend
COPY --from=builder /app/etc /app/etc
COPY --from=builder /app/http /app/http
COPY --from=builder /out/web /app/web
ENV WEB_DIR=/app/web
CMD ["/app/backend"]
