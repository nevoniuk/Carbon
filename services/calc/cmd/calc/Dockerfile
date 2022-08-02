FROM golang:1.17.1 as builder
WORKDIR /src

ARG GITHUB_USER
ARG GITHUB_TOKEN
RUN git config --global url."https://$GITHUB_USER:$GITHUB_TOKEN@github.com".insteadOf "https://github.com"

ENV GOPRIVATE="github.com/crossnokaye/*"

# Be explicit about what we need for faster context upload and better caching
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

ARG GITHUB_SHA
ENV CGO_ENABLED=0
RUN go build \
  -tags netgo \
  -a \
  -ldflags "-extldflags \"-static\" -X goa.design/clue/health.GitCommit=$GITHUB_SHA" \
  -v \
  -o build/output/calc \
  ./services/calc/cmd/calc/

FROM alpine:3.13.3
RUN apk update && apk add --no-cache ca-certificates tzdata
WORKDIR /app/

COPY --from=builder /src/build/output/calc .

ENV GITHUB_SHA ${GITHUB_SHA}
ENTRYPOINT ["/app/calc"]
