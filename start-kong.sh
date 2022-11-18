docker rm -f kong-gateway-jwe-decrypt

docker run -d --name kong-gateway-jwe-decrypt \
--network=kong-net \
--link kong-database-jwe-decrypt:kong-database-jwe-decrypt \
-e "KONG_DATABASE=postgres" \
-e "KONG_PG_HOST=kong-database-jwe-decrypt" \
-e "KONG_PG_USER=kong" \
-e "KONG_PG_PASSWORD=kongpass" \
-e "KONG_PROXY_ACCESS_LOG=/dev/stdout" \
-e "KONG_ADMIN_ACCESS_LOG=/dev/stdout" \
-e "KONG_PROXY_ERROR_LOG=/dev/stderr" \
-e "KONG_ADMIN_ERROR_LOG=/dev/stderr" \
-e "KONG_ADMIN_LISTEN=0.0.0.0:8001, 0.0.0.0:8444 ssl" \
-e "KONG_ADMIN_GUI_URL=http://localhost:8002" \
-e "KONG_UNTRUSTED_LUA_SANDBOX_REQUIRES=resty.http,cjson.safe" \
-e "KONG_PLUGINS=bundled,jwe-decrypt" \
-e KONG_LICENSE_DATA \
-p 8000:8000 \
-p 8443:8443 \
-p 8001:8001 \
-p 8002:8002 \
-p 8444:8444 \
kong-gateway-jwe-decrypt

echo "see logs 'docker logs kong-gateway-jwe-decrypt -f'"