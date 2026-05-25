FROM golang:1.25.5-alpine3.23

# Adding the Other things required things 
RUN apk add --no-cache \
    build-base \
    gcc \
    g++ \
    curl \
    python3 \
    py3-pip \
    && rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

VOLUME [ "/artifacts" , "/code" ]
COPY . .

RUN go build -o builder ./cmd/


ENTRYPOINT [ "./builder" ]