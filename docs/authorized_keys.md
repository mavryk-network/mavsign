---
id: authorized_keys
title: Authorized_Keys Configuration
---
# MavSign's Authorized Key Authentication Feature

MavSign provides the option to authenticate the mavkit-client, by specifying an "authorized key" in the MavSign configuration file.  

## Motivation

An authorized key can be configured to ensure that MavSign only signs requests from an mavkit-client instance containing the private key.

## Configuration

First, a key pair is generated using mavkit-client:

```bash
mavkit-client gen keys mavsign-auth
```

Next, find the public key value:

```bash
cat ~/.mavryk-client/public_keys | grep -C 3 mavsign-auth
```

Finally, add the public key value to the MavSign configuration file.  It belongs within the `server` declaration:

```yaml
server:
  address: :6732
  utility_address: :9583
  authorized_keys:
    - edpkujLb5ZCZ2gprnRzE9aVHKZfx9A8EtWu2xxkwYSjBUJbesJ9rWE
```

Restarting the MavSign service is required to apply the configuration change.  Henceforth, the MavSign service will only accept requests from the mavkit-client that is using the private key associated with the public key specified in the configuration file.
