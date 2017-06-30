#FROM golang:latest
FROM golang:latest
LABEL description="aiden's life"

RUN mkdir -p /go/src \
      && mkdir -p /go/bin \
      && mkdir -p /go/pkg

ENV GOPATH=/go

RUN mkdir -p $GOPATH/src/code.ndella.com/ai-life
ADD . $GOPATH/src/code.ndella.com/ai-life
WORKDIR $GOPATH/src/code.ndella.com/ai-life

RUN go build -o ai-life
ENTRYPOINT ./ai-life

EXPOSE 3030