FROM golang:1.22.0

RUN go version
ENV GOPATH=/
COPY ./ ./

RUN go mod download
RUN go build -o marketplace-api ./cmd/*

CMD ["./marketplace-api"]
