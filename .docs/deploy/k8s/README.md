# Deploy to Kubernetes

## Create manifests

### 1. `pvc.yaml`

<!-- TODO -->

### 2. `urlshortener.yaml`

<!-- TODO -->

### 3. [tingtt/oauth2rbac](https://github.com/tingtt/oauth2rbac) (optional / if you needed)

Use auth proxy ([tingtt/oauth2rbac](https://github.com/tingtt/oauth2rbac)) if needed.
It will add authentication to the shortened URL reference and registration.

`config.yaml`

```yaml
# https://kubernetes.io/docs/concepts/configuration/configmap/
kind: ConfigMap
apiVersion: v1
metadata:
  name: oauth2rbac
data:
  oauth2rbac.yml: |
    proxies:
      - external_url: "https://example.com/"
        target: "http://urlshortener.urlshortener.svc.cluster.local:8080/"

    acl:
      "-": #! public
        - external_url: "https://example.com/"
          methods: ["GET"]
      "<your email>": #* only you can edit
        - external_url: "https://example.com/"
          methods: ["*"]
```

**Configure it based on the above `config.yaml` and the following document.**
[github.com/tingtt/oauth2rbac - Deploy to Kubernetes](https://github.com/tingtt/oauth2rbac/tree/main/.docs/deploy/k8s)

## Apply manifests

```sh
kubectl apply -f pvc.yaml -f urlshortener.yaml
kubectl apply -f secret.yaml -f config.yaml -f oauth2rbac.yaml -f networkPolicy.yaml # optional
```
