# WB Tech L0

## Build

```bash
podman build . -f build/db/Containerfile -t wb-tech-l0-db
podman build . -f build/app/Containerfile -t wb-tech-l0-app
```

## Deploy

```bash
podman kube play deploy/wb-tech-l0-pod.yml
```

## Shutdown

```bash
podman kube down deploy/wb-tech-l0-pod.yml
```

To prune associated volume add `--force` flag to this command:
```bash
podman kube down --force deploy/wb-tech-l0-pod.yml
```
