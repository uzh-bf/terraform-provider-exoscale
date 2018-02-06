# Exokube simple setup

```console
$ terraform apply

...

$ ssh ubuntu@...

# wait for it
exokube $ tail -f /var/log/cloud-init-output.log

Cloud-init v. 17.1 finished at ...
```

## Installing a network

Using [Calico](https://docs.projectcalico.org/v3.1/getting-started/kubernetes/).

```console
$ kubectl apply -f https://docs.projectcalico.org/v3.1/getting-started/kubernetes/installation/hosted/kubeadm/1.7/calico.yaml
```

Wait for them...

```console
$ watch kubectl get pods --all-namespaces

NAMESPACE     NAME                                      READY     STATUS    RESTARTS   AGE
kube-system   calico-etcd-mw69k                         1/1       Running   0          1m
kube-system   calico-kube-controllers-99bdf4cd4-kcpz9   1/1       Running   0          1m
kube-system   calico-node-bmqjd                         2/2       Running   0          1m
kube-system   etcd-minikube                             1/1       Running   0          1m
kube-system   kube-apiserver-minikube                   1/1       Running   0          1m
kube-system   kube-controller-manager-minikube          1/1       Running   0          2m
kube-system   kube-dns-86f4d74b45-jlg7m                 3/3       Running   0          2m
kube-system   kube-proxy-cmrtp                          1/1       Running   0          2m
kube-system   kube-scheduler-minikube                   1/1       Running   0          1m
```

Allow pods to be scheduled on the master node.

```
$ kubectl taint nodes --all node-role.kubernetes.io/master-
node "exokube" untainted
```

According to Calico's documentation, we're done.

```
kubectl get nodes -o wide
NAME       STATUS    ROLES     AGE       VERSION   EXTERNAL-IP   OS-IMAGE             KERNEL-VERSION      CONTAINER-RUNTIME
exokube    Ready     master    3m        v1.10.1   <none>        Ubuntu 16.04.4 LTS   4.4.0-112-generic   docker://17.3.2
```

## Setting up the Dashboard

[Web UI (Dashboard)](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/)

```
$ kubectl create -f https://raw.githubusercontent.com/kubernetes/dashboard/master/src/deploy/recommended/kubernetes-dashboard.yaml
```

## TODO

- How to access the dashboard
- Running a service
