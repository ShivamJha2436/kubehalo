# KubeHalo

KubeHalo is a Kubernetes custom autoscaler built in Go.  
It introduces a **Custom Resource Definition (CRD)** called `ScalePolicy` to define scaling rules for workloads.  
The controller watches these resources in real time, evaluates metrics, and adjusts workloads accordingly.

---

## Current Status
- ✅ `ScalePolicy` CRD defined and deployed
- ✅ Controller implemented to watch and log CRD events
- 🚧 Prometheus metrics integration
- 🚧 Scaling engine for deployments
- 🚧 Web dashboard for management and visualization

---

## How It Works
1. **Define Scaling Rules**  
   Create a `ScalePolicy` CRD specifying min/max replicas, metric thresholds, and target workloads.
2. **Controller Watches Policies**  
   Uses Kubernetes Informers to detect create/update/delete events in real time.
3. **Scaling Decisions** *(coming soon)*  
   Metrics will be fetched from Prometheus and compared to thresholds to trigger scaling actions.

---

## Quick Start

### 1. Install the CRD
```bash
kubectl apply -f config/crd/scale_policy.yaml
```

### 2. Run the Controller Locally
```bash
go run ./cmd/controller
```
### 3. Apply a Sample Policy
```bash
kubectl apply -f manifests/sample_scale_policy.yaml
```

### Example ScalePolicy
```yaml
apiVersion: kubehalo.sh/v1
kind: ScalePolicy
metadata:
  name: test-policy
  namespace: default
spec:
  minReplicas: 1
  maxReplicas: 10
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
  targetRef:
    kind: Deployment
    name: my-deployment
    namespace: default

```
