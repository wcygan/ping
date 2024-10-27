# Setting up CloudNativePG (CNPG) on Kubernetes

This guide outlines the essential steps to set up CloudNativePG for your Kubernetes application.

## Prerequisites
- A running Kubernetes cluster
- Helm installed
- kubectl installed

## 1. Install CloudNativePG Operator

Use Helm to install the CNPG operator:

```yaml
# Example Helm deployment in skaffold.yaml
apiVersion: skaffold/v4beta1
kind: Config
metadata:
  name: ping-db-install
deploy:
  helm:
    releases:
      - name: cnpg
        repo: https://cloudnative-pg.github.io/charts
        remoteChart: cloudnative-pg
        namespace: cnpg-system
        createNamespace: true
        wait: true
```

## 2. Create Database Credentials

Create a Kubernetes secret for database credentials:

```yaml
# secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: db-user-secret
type: Opaque
stringData:
  username: youruser
  password: yourpassword
```

## 3. Define the PostgreSQL Cluster

Create a CNPG cluster definition:

```yaml
# cluster.yaml
apiVersion: postgresql.cnpg.io/v1
kind: Cluster
metadata:
  name: pg-cluster
spec:
  instances: 1  # Increase for production
  storage:
    size: 1Gi
  bootstrap:
    initdb:
      database: yourdb
      owner: youruser
      secret:
        name: db-user-secret
      postInitSQL:
        - ALTER ROLE youruser WITH LOGIN;
```

## 4. Important Notes

### Database Access
- The primary instance will be available at: `<cluster-name>-rw.<namespace>.svc`
- Default port: 5432

### High Availability
- For production, consider:
  - Increasing `instances` to 3 for high availability
  - Adjusting storage size based on your needs
  - Setting up proper backup strategies

### Resource Management
- Define appropriate resource requests and limits
- Consider storage class requirements for your environment

## 5. Deployment Strategy

1. First deploy the operator
2. Apply the secret
3. Deploy the cluster definition
4. Wait for the cluster to be ready before deploying applications

## 6. Verification

Check cluster status:
```bash
kubectl get clusters
kubectl get pods -l postgresql
```

The cluster is ready when you see:
- All pods in Running state
- Cluster status shows "Cluster in healthy state"

## 7. Application Integration

### Init Container Setup
Use init containers to ensure database availability before starting your application:

```yaml
initContainers:
  - name: wait-for-k8s
    image: bitnami/kubectl
    command:
      - /bin/sh
      - -c
      - |
        echo "Waiting for CNPG cluster to be ready..."
        until kubectl get cluster pg-cluster -o jsonpath='{.status.phase}' | grep -q "Cluster in healthy state"; do
          echo "Waiting for cluster to be healthy..."
          sleep 5
        done
        
        # Verify primary is available
        PRIMARY_HOST="pg-cluster-rw"
        until kubectl get svc $PRIMARY_HOST; do
          echo "Waiting for primary service to be available..."
          sleep 5
        done

  - name: wait-for-db
    image: bitnami/postgresql
    env:
      - name: PGPASSWORD
        valueFrom:
          secretKeyRef:
            name: db-user-secret
            key: password
      - name: PGUSER
        valueFrom:
          secretKeyRef:
            name: db-user-secret
            key: username
    command:
      - /bin/sh
      - -c
      - |
        PRIMARY_HOST="pg-cluster-rw"
        until psql -h $PRIMARY_HOST -d yourdb -c '\l'; do
          echo "Waiting for primary database to accept connections..."
          sleep 5
        done
```

### Environment Configuration
Set up your application's database configuration using environment variables:

```yaml
env:
  - name: DB_HOST
    value: pg-cluster-rw
  - name: DB_NAME
    value: yourdb
  - name: DB_USER
    valueFrom:
      secretKeyRef:
        name: db-user-secret
        key: username
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: db-user-secret
        key: password
```

## 8. RBAC Configuration

Create necessary RBAC permissions for your application to interact with CNPG:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: cnpg-reader
rules:
- apiGroups: ["postgresql.cnpg.io"]
  resources: ["clusters"]
  verbs: ["get", "watch", "list"]
- apiGroups: [""]
  resources: ["services", "pods"]
  verbs: ["get", "watch", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: cnpg-reader-binding
subjects:
- kind: ServiceAccount
  name: default
roleRef:
  kind: Role
  name: cnpg-reader
  apiGroup: rbac.authorization.k8s.io
```

## Best Practices

1. Always use secrets for credentials
2. Monitor storage usage
3. Implement regular backup strategies
4. Use readiness probes in applications depending on the database
5. Consider using init containers to wait for database availability

## Common Issues

1. Storage provisioning failures
   - Ensure your cluster has appropriate storage classes available
2. Authentication issues
   - Verify secret contents and references
3. Connection timeouts
   - Check network policies and service names
