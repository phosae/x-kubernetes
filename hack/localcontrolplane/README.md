# run control plane in Docker for API testing

set up

```bash
{
./init-controlplane-pki.sh
docker-compose up
}
```

verify with kubectl

```
kubectl --kubeconfig ./api.kubeconfig get sa,svc,cm -A
NAMESPACE         NAME                     SECRETS   AGE
default           serviceaccount/default   0         31s
kube-node-lease   serviceaccount/default   0         31s
kube-public       serviceaccount/default   0         31s
kube-system       serviceaccount/default   0         31s

NAMESPACE   NAME                 TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)   AGE
default     service/kubernetes   ClusterIP   172.20.0.1   <none>        443/TCP   32s

NAMESPACE         NAME                                                             DATA   AGE
default           configmap/kube-root-ca.crt                                       1      31s
kube-node-lease   configmap/kube-root-ca.crt                                       1      31s
kube-public       configmap/kube-root-ca.crt                                       1      31s
kube-system       configmap/kube-apiserver-legacy-service-account-token-tracking   1      33s
kube-system       configmap/kube-root-ca.crt                                       1      31s
```