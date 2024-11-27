# Deploy to Kubernetes

## Create manifests

### 1. `pvc.yaml`

```yaml
# https://kubernetes.io/docs/concepts/storage/persistent-volumes/
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: urlshortener
  labels:
    app: urlshortener
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 2Gi
```

### 2. `urlshortener.yaml`

```yaml
# https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
apiVersion: apps/v1
kind: Deployment
metadata:
  name: urlshortener
  labels:
    app: urlshortener
spec:
  selector:
    matchLabels:
      app: urlshortener
  replicas: 1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: urlshortener
    spec:
      containers:
        - name: urlshortener
          image: tingtt/urlshortener:v2.0.0
          imagePullPolicy: IfNotPresent
          resources:
            limits:
              cpu: 80m
              memory: 100Mi
            requests:
              cpu: 40m
              memory: 40Mi
          livenessProbe:
            tcpSocket:
              port: 8080
            initialDelaySeconds: 5
            timeoutSeconds: 5
            successThreshold: 1
            failureThreshold: 3
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /?edit
              port: 8080
            initialDelaySeconds: 5
            timeoutSeconds: 2
            successThreshold: 1
            failureThreshold: 3
            periodSeconds: 10
          ports:
            - containerPort: 8080
              name: urlshortener
          volumeMounts:
            - name: nfs
              mountPath: /var/lib/urlshortener
      volumes:
        - name: nfs
          persistentVolumeClaim:
            claimName: urlshortener
      restartPolicy: Always
---
# https://kubernetes.io/docs/concepts/services-networking/service/
apiVersion: v1
kind: Service
metadata:
  name: urlshortener
spec:
  selector:
    app: urlshortener
  type: ClusterIP
  ports:
    - name: urlshortener
      protocol: TCP
      port: 8080
      targetPort: 8080
```

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
