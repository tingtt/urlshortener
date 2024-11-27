# Deploy with Docker Compose

## 1. Create `compose.yaml`

Create a `compose.yaml` file to define the service.

```yaml
services:
  urlshortener:
    image: tingtt/urlshortener:v2.0.0
    command: [
      "--port", "8080",
      "--save.dir", "/var/lib/urlshortener/",
    ]
    ports:
      - "8080:8080"
    volumes:
      - ./data:/var/lib/urlshortener/
    restart: always
```

## 2. [tingtt/oauth2rbac](https://github.com/tingtt/oauth2rbac) (optional / if you needed)

Use auth proxy ([tingtt/oauth2rbac](https://github.com/tingtt/oauth2rbac)) if needed.
It will add authentication to the shortened URL reference and registration.

`config.yaml`

```yaml
proxies:
  - external_url: "https://example.com/"
    target: "http://urlshortener:8080/"

acl:
  "-": #! public
    - external_url: "https://example.com/"
      methods: ["GET"]
  "<your email>": #* only you can edit
    - external_url: "https://example.com/"
      methods: ["*"]
```

**Configure it based on the above `config.yaml` and the following document.**
[github.com/tingtt/oauth2rbac - Deploy to Docker Compose](https://github.com/tingtt/oauth2rbac/tree/main/.docs/deploy/docker)

## Run

```sh
docker compose up -d
```
