---
id: file_based
title: File-Based Secret Storage (Insecure)
---


MavSign file-based signer mode allows operation without an HSM or Key Vault service for evaluation and prototyping purposes. By storing the secret key material in a JSON file, users can get MavSign up and running quickly for evaluation and development purposes.

## MavSign configuration for file-based secret storage

This documentation assumes you will use the official MavSign docker image, and that you have a working Linux server with docker installed.

Place the following YAML in a file named `mavsign.yaml`

```yaml
server:
  address: :6732
  utility_address: :9583

vaults:
  local_secret:
    driver: file
    config:
      file: /etc/secret.json

mavryk:
  mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt:
    log_payloads: true
    allow:
      block:
      endorsement:
      preendorsement:
      generic:
        - transaction
```

The `mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt` key corresponds to the secret key that you will put in `/etc/secret.json`

Contents of `secret.json` is:

```json
[ 
  { 
    "name": "<address_alias>",
    "value": "unencrypted:<your_secret_key>" 
  }
]
```

### Running MavSign

Next, you want to run the mavsign docker image as follows:

_Remember to secure the network where MavSign is running_

```sh
docker run -it --rm \
    -v "$(realpath mavsign.yaml):/etc/mavsign.yaml" \
    -v "$(realpath secret.json):/etc/secret.json" \
    -p 6732:6732 \
    -p 9583:9583 \
    mavrykdynamics/mavsign:latest serve -c /etc/mavsign.yaml
```

### Verify that MavSign is working

You can test that MavSign works, making a GET request using the Public Key Hash (PKH). MavSign will return a JSON payload containing the public key.

```sh
curl mavsign:6732/keys/mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt
```

A response such as the following should be expected:

```json
{"public_key":"edpktn6UGrMQUjhWQJ5kY4qWoCp1sDZWkK5ugizTc5jHSifG1j3r8o"}
```

You can test the signing functionality by making a POST request as follows:

```sh
curl -XPOST \
    -d '"027a06a770e6cebe5b3e39483a13ac35f998d650e8b864696e31520922c7242b88c8d2ac55000003eb6d"' \
    mavsign:6732/keys/mv19VEmW4zEELeQiBqLHH4RHgysYuLe4P6xt
```

Which should return an HTTP 200 OK with a payload similar to:

```json
{"signature":"sigWetzF5zVM2qdYt8QToj7e5cNBm9neiPRc3rpePBDrr8N1brFbErv2YfXMSoSgemJ8AwZcLfmkBDg78bmUEzF1sf1YotnS"}
```

If you repeat the same signing operation more than once, you will get an error from the High-Watermark feature. This safety measure prevents the injection of duplicate operations.

The payload on this request resembles a Mavryk endorsement emitted from a Mavryk Baker node.

