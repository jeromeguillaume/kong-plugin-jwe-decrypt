# jwe-decrypt: Kong custom plugin
- Take a JWE found in "Auhorization: Bearer" header
- Decrypt JWE to JWS
- Set JWS in "Auhorization: Bearer" header

## How deploy jwe-decrypt plugin
1. Build a new image of kong-gateway
```
docker build -t kong-gateway-jwe-decrypt .
```
2. Modify the KONG_PLUGINS environment variable to include the `jwe-decrypt` plugin
```
KONG_PLUGINS=bundled,jwe-decrypt
````
3. Use the new image `kong-gateway-jwe-decrypt` in your deployment
4. Install the `jwe-decrypt` plugin
5. Configure the PEM private key **without**:
    - -----BEGIN RSA PRIVATE KEY-----
    - -----END RSA PRIVATE KEY-----
    - Return carriage