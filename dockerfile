
FROM golang:1.20.3-alpine
WORKDIR /home/soll/go/src/pipeline/
COPY go.mod ./
RUN go mod download
COPY task20_2_1.go .
RUN go build -o /pipeline


LABEL version="1.0"
LABEL maintainer="SF_strudent<SF_student@test.ru>"
CMD [ "/pipeline" ]