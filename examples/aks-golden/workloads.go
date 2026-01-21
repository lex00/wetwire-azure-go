// Package aks_golden includes sample K8s workloads that would be deployed to the cluster.
//
// Note: These are represented as ARM template resources or could be
// deployed separately using kubectl/Helm after cluster creation.
// This file demonstrates the K8s manifests that would run on the cluster.
package aks_golden

// SampleDeploymentYAML represents a sample Kubernetes Deployment manifest
// that would be deployed to the cluster. In practice, this would be deployed
// via kubectl, Helm, or GitOps tools like ArgoCD/Flux.
//
// This is included as documentation of what a multi-cloud portable workload looks like.
const SampleDeploymentYAML = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-app
  namespace: default
  labels:
    app: sample-app
spec:
  replicas: 3
  selector:
    matchLabels:
      app: sample-app
  template:
    metadata:
      labels:
        app: sample-app
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      nodeSelector:
        nodepool-type: application
      tolerations:
        - key: "kubernetes.azure.com/scalesetpriority"
          operator: "Equal"
          value: "spot"
          effect: "NoSchedule"
      containers:
        - name: app
          image: nginx:1.25
          ports:
            - containerPort: 80
              name: http
            - containerPort: 8080
              name: metrics
          resources:
            requests:
              cpu: "100m"
              memory: "128Mi"
            limits:
              cpu: "500m"
              memory: "512Mi"
          readinessProbe:
            httpGet:
              path: /healthz
              port: 80
            initialDelaySeconds: 5
            periodSeconds: 10
          livenessProbe:
            httpGet:
              path: /healthz
              port: 80
            initialDelaySeconds: 15
            periodSeconds: 20
---
apiVersion: v1
kind: Service
metadata:
  name: sample-app
  namespace: default
  labels:
    app: sample-app
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "8080"
spec:
  type: ClusterIP
  selector:
    app: sample-app
  ports:
    - name: http
      port: 80
      targetPort: 80
    - name: metrics
      port: 8080
      targetPort: 8080
---
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: sample-app
  namespace: default
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: sample-app
  minReplicas: 3
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70
    - type: Resource
      resource:
        name: memory
        target:
          type: Utilization
          averageUtilization: 80
`

// SampleIngressYAML represents a sample AGIC (Azure Application Gateway Ingress Controller)
// configuration that would be deployed to the cluster.
const SampleIngressYAML = `
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: sample-app
  namespace: default
  annotations:
    kubernetes.io/ingress.class: azure/application-gateway
    appgw.ingress.kubernetes.io/ssl-redirect: "true"
    appgw.ingress.kubernetes.io/backend-protocol: "http"
    appgw.ingress.kubernetes.io/health-probe-path: "/healthz"
spec:
  tls:
    - hosts:
        - app.example.com
      secretName: sample-app-tls
  rules:
    - host: app.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: sample-app
                port:
                  number: 80
`

// WorkloadIdentityYAML demonstrates Azure Workload Identity configuration.
// This allows pods to authenticate to Azure services using managed identity.
const WorkloadIdentityYAML = `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: sample-app-sa
  namespace: default
  annotations:
    # Replace with your User Assigned Managed Identity Client ID
    azure.workload.identity/client-id: "<MANAGED_IDENTITY_CLIENT_ID>"
  labels:
    azure.workload.identity/use: "true"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-app-with-identity
  namespace: default
spec:
  selector:
    matchLabels:
      app: sample-app-identity
  template:
    metadata:
      labels:
        app: sample-app-identity
    spec:
      serviceAccountName: sample-app-sa
      containers:
        - name: app
          image: mcr.microsoft.com/azure-cli:latest
          command: ["sleep", "infinity"]
          # Azure SDK will automatically use workload identity
          # No explicit credential configuration needed
`

// KedaScalerYAML demonstrates KEDA (Kubernetes Event-Driven Autoscaling)
// configuration for Azure Queue-based scaling.
const KedaScalerYAML = `
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: sample-app-scaler
  namespace: default
spec:
  scaleTargetRef:
    name: sample-app
  minReplicaCount: 1
  maxReplicaCount: 100
  triggers:
    - type: azure-queue
      metadata:
        queueName: sample-queue
        accountName: samplestorageaccount
        queueLength: "10"
      authenticationRef:
        name: azure-queue-auth
---
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: azure-queue-auth
  namespace: default
spec:
  podIdentity:
    provider: azure-workload
`
