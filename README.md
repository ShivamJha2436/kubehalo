# KubeHalo

KubeHalo is a Kubernetes autoscaling prototype written in Go. It introduces a `ScalePolicy` custom resource and a small control plane around it:

- a controller that watches `ScalePolicy` objects
- a Prometheus-backed scaling decision flow
- a validation webhook for admission checks
- a lightweight HTTP API for listing policies

The codebase is now aligned around one `ScalePolicy` contract, has runnable entrypoints under `cmd/`, and includes automated tests for the critical scaling and validation paths.

## What The Project Does

![KubeHalo Architecture](./docs/architecture/kubehalo-architecture.png)

Each `ScalePolicy` points at a workload, defines the metric query to evaluate, and describes how aggressively KubeHalo may scale up or down.

Current behavior:

- watches `ScalePolicy` resources through a dynamic informer
- queries Prometheus for the configured metric
- reads current Deployment replicas from the cluster
- computes desired replicas from metric threshold and step sizes
- applies optional behavior rules such as stabilization windows and rate caps
- updates the target Deployment when a scale action is needed

## Repository Layout

```text
.
├── api/kubehalo/v1          # Typed ScalePolicy API definitions
├── cmd/api                  # HTTP API entrypoint
├── cmd/controller           # Controller entrypoint
├── cmd/webhook              # Admission webhook entrypoint
├── config/crd               # CRD manifests
├── config/webhook           # Webhook registration manifest
├── controllers/scalepolicy  # Informer, handler, parsing, scaling logic
├── internal/config          # Environment-based runtime configuration
├── internal/kube            # Kubernetes client construction
├── internal/metrics         # Prometheus client wrapper
├── internal/scaling         # Deployment scaling engine
└── manifests                # Example Kubernetes manifests
```

## ScalePolicy Shape

```yaml
apiVersion: kubehalo.sh/v1
kind: ScalePolicy
metadata:
  name: demo-policy
  namespace: default
spec:
  targetRef:
    kind: Deployment
    name: my-deployment
    namespace: default
  metric:
    name: cpu
    query: rate(container_cpu_usage_seconds_total[1m])
    threshold: 0.8
  scaleUp:
    step: 2
    cooldownSeconds: 60
  scaleDown:
    step: 1
    cooldownSeconds: 120
  minReplicas: 1
  maxReplicas: 10
  behavior:
    stabilizationWindowSeconds: 60
    maxScaleUpRate: 2
    maxScaleDownRate: 1
    policy: absolute
```

Sample manifest: `manifests/sample-policy.yaml`

## Quick Start

### Prerequisites

- Go 1.24+
- access to a Kubernetes cluster or local cluster such as Minikube
- a reachable Prometheus instance
- a valid `KUBECONFIG` when running outside the cluster

### 1. Install the CRD

```bash
kubectl apply -f config/crd/scale_policy.yaml
```

### 2. Run the controller locally

```bash
export KUBEHALO_PROMETHEUS_ADDR=http://localhost:9090
go run ./cmd/controller
```

### 3. Apply a policy

```bash
kubectl apply -f manifests/sample-policy.yaml
```

### 4. Optionally run supporting services

HTTP API:

```bash
go run ./cmd/api
```

Webhook:

```bash
go run ./cmd/webhook
```

## Configuration

KubeHalo reads runtime configuration from environment variables.

| Variable | Default | Used By |
| --- | --- | --- |
| `KUBEHALO_PROMETHEUS_ADDR` | `http://localhost:9090` | controller |
| `KUBEHALO_API_ADDR` | `:8080` | API server |
| `KUBEHALO_WEBHOOK_ADDR` | `:8443` | webhook server |
| `KUBEHALO_WEBHOOK_CERT_FILE` | `/tls/tls.crt` | webhook server |
| `KUBEHALO_WEBHOOK_KEY_FILE` | `/tls/tls.key` | webhook server |

## Development

Useful commands:

```bash
make fmt
make lint
make test
make run-controller
make run-api
make run-webhook
```

The `Makefile` keeps Go build caches inside `.cache/`, which makes local iteration cleaner and avoids polluting global caches.

## Testing

The repository includes tests for:

- controller helper construction
- `ScalePolicy` parsing and validation
- handler-driven scaling decisions
- Deployment scaling engine behavior
- webhook admission validation

Run everything with:

```bash
make test
```

## Notes

- The current scaling engine updates `Deployment` targets. The API type already allows `StatefulSet`, but reconciliation for that target kind is not implemented yet.
- `cooldownSeconds` and `evaluationIntervalSeconds` are modeled in the API but are not yet enforced by the controller.
