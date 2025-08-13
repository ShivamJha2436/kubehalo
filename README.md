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