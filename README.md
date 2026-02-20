# fake-apiserver

**kubectl**

```console
> k create ns demo --validate=false

> k --kubeconfig=local.kubeconfig create -f examples/cm.yaml --validate=false
```

**patch apis**

- https://kubernetes.io/docs/tasks/manage-kubernetes-objects/update-api-object-kubectl-patch/

```
> k --kubeconfig=local.kubeconfig create -f https://k8s.io/examples/application/deployment-patch.yaml --validate=false

> k --kubeconfig=local.kubeconfig patch deployments.apps patch-demo --patch-file examples/patch-file.yaml

> k --kubeconfig=local.kubeconfig patch deployments.apps patch-demo --patch-file examples/patch-file-tolerations.yaml

> k --kubeconfig=local.kubeconfig patch deployments.apps patch-demo --type merge --patch-file examples/patch-file-2.yaml
```

```
> k --kubeconfig=local.kubeconfig get deployments.apps -A -o yaml
```

**helm chart**

```console
helm install hello examples/charts/hello --disable-openapi-validation
```

```console
helm install kubedb appscode/kubedb \
  --version v2023.08.18 \
  --namespace kubedb --create-namespace \
  --set kubedb-provisioner.enabled=true \
  --set kubedb-ops-manager.enabled=true \
  --set kubedb-autoscaler.enabled=true \
  --set kubedb-dashboard.enabled=true \
  --set kubedb-schema-manager.enabled=true \
  --disable-openapi-validation \
  --wait=false \
  --wait-for-jobs=false
```

***crd***

```console
k create -f examples/crds/cert-manager.io_clusterissuers.yaml --validate=false
k create -f examples/crds/clusterissuer.yaml --validate=false
```

ToDos:

- [ ] status
- [ ] scale
- [ ] Delete via owner ref
- [ ] openapi
