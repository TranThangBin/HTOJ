FROM docker.io/library/golang:1.25.5-bookworm
WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v4.1.17/tailwindcss-linux-x64
RUN chmod +x ./tailwindcss-linux-x64
ENTRYPOINT [ "go", "tool", "air" ]
