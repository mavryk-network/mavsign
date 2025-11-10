FROM golang:1.24-bullseye AS builder
RUN apt-get update && apt-get install -y make
ADD . /mavsign
WORKDIR /mavsign
RUN make

FROM debian:bullseye-slim
WORKDIR /mavsign
RUN apt update -y \
    && apt install -y curl apt-transport-https\
    && rm -rf /var/lib/apt/lists/*
COPY --from=builder /mavsign/mavsign.yaml /mavsign/mavsign.yaml
COPY --from=builder /mavsign/mavsign /usr/bin/mavsign
COPY --from=builder /mavsign/mavsign-cli /usr/bin/mavsign-cli

ENTRYPOINT ["/usr/bin/mavsign"]
CMD [ "-c", "/mavsign/mavsign.yaml" ]
