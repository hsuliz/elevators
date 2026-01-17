FROM node:22-alpine AS frontend-builder
WORKDIR /app

RUN npm install -g esbuild

COPY web/ ./web/
COPY static/ ./static/

RUN esbuild web/main.ts --bundle --minify --sourcemap --outfile=static/main.js

FROM golang:1.25-alpine AS go-builder
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

COPY --from=frontend-builder /app/static/main.js ./static/main.js

RUN CGO_ENABLED=0 GOOS=linux go build -o /elevators ./cmd/main.go

FROM gcr.io/distroless/static-debian12
WORKDIR /

#RUN apk --no-cache add ca-certificates

COPY --from=go-builder /elevators .
COPY --from=go-builder /app/static ./static

EXPOSE 8080

CMD ["/elevators"]