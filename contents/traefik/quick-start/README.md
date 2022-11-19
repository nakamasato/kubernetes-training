# Traefik

https://traefik.io/

## Traefik Proxy

https://doc.traefik.io/traefik/

Example with docker-compose:

```
docker compose up
```

```
url -H Host:whoami.docker.localhost http://127.0.0.1
```

<details>

```
Hostname: 6d0f82923467
IP: 127.0.0.1
IP: 172.19.0.2
RemoteAddr: 172.19.0.3:48168
GET / HTTP/1.1
Host: whoami.docker.localhost
User-Agent: curl/7.79.1
Accept: */*
Accept-Encoding: gzip
X-Forwarded-For: 172.19.0.1
X-Forwarded-Host: whoami.docker.localhost
X-Forwarded-Port: 80
X-Forwarded-Proto: http
X-Forwarded-Server: 96aff624c7ed
X-Real-Ip: 172.19.0.1
```

</details>

[Example with Kubernetes](https://doc.traefik.io/traefik/getting-started/quick-start-with-kubernetes/)

## Traefik Mesh

https://doc.traefik.io/traefik-mesh/

## Traefik Pilot

https://doc.traefik.io/traefik-pilot/
