FROM golang:latest AS build_stage
RUN mkdir -p go/src/gonews
WORKDIR /go/src/gonews
COPY ./ ./
RUN go env -w GO111MODULE=auto 
RUN go install ./cmd/news
RUN cp ./cmd/news/config.json ./go/bin


FROM alpine:latest
RUN mkdir -p gonews
WORKDIR /gonews
COPY --from=build_stage /go/bin .
RUN apk add libc6-compat
##ENTRYPOINT ./SF36a.3
