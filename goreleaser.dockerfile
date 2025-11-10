FROM ubuntu:22.04
RUN apt update
RUN apt install -y curl apt-transport-https wget

# Download and install CloudHSM PKCS11 based on architecture
ARG TARGETARCH
RUN if [ "$TARGETARCH" = "amd64" ]; then \
    wget https://s3.amazonaws.com/cloudhsmv2-software/CloudHsmClient/Jammy/cloudhsm-pkcs11_latest_u22.04_amd64.deb && \
    apt install -y ./cloudhsm-pkcs11_latest_u22.04_amd64.deb && \
    rm cloudhsm-pkcs11_latest_u22.04_amd64.deb; \
    elif [ "$TARGETARCH" = "arm64" ]; then \
    wget https://s3.amazonaws.com/cloudhsmv2-software/CloudHsmClient/Jammy/cloudhsm-pkcs11_latest_u22.04_arm64.deb && \
    apt install -y ./cloudhsm-pkcs11_latest_u22.04_arm64.deb && \
    rm cloudhsm-pkcs11_latest_u22.04_arm64.deb; \
    fi

# Add CloudHSM to PATH
ENV PATH="/opt/cloudhsm/bin:${PATH}"

COPY ./mavsign /bin
COPY ./mavsign-cli /bin

ENTRYPOINT ["/bin/mavsign"]