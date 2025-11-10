FROM golang:1.14-alpine AS build-env

ENV COLLECTOR_PKG="github.com/mavryk-network/mavsign/pkg/metrics"
ENV GO111MODULE="on"
WORKDIR /build/mavsign
RUN apk update && apk add git openssh bash
ADD . .
RUN go build -ldflags "-X ${COLLECTOR_PKG}.GitRevision=${GIT_REVISION} -X ${COLLECTOR_PKG}.GitBranch=${GIT_BRANCH} -X ${COLLECTOR_PKG}.GitVersion=${GIT_VERSION}" ./cmd/mavsign
RUN go build -ldflags "-X ${COLLECTOR_PKG}.GitRevision=${GIT_REVISION} -X ${COLLECTOR_PKG}.GitBranch=${GIT_BRANCH} -X ${COLLECTOR_PKG}.GitVersion=${GIT_VERSION}" ./cmd/mavsign-cli


# final stage
FROM alpine:3
RUN apk --no-cache add ca-certificates
COPY --from=build-env /build/mavsign/mavsign /bin
COPY --from=build-env /build/mavsign/mavsign-cli /bin

ENTRYPOINT ["/bin/mavsign"]
