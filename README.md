# programming-k8s
a kubernetes programming example (using CRD and custom controller):
watch all pods events and store their cpu request in a custom resource.

## build and run
```
# create rbac
$ kubectl create -f config/rbac/rbac_role.yaml
$ kubectl create -f config/rbac/rbac_role_binding.yaml

# create crd
$ kubectl create -f config/crds/stats-cpu-crd.yaml

# create custom resource (for store pod's cpu request)
$ kubectl create -f config/samples/stats-cpu-sample.yaml

# install vendor and run controller (watch pod event)
$ make dep
$ make run
```

## demo
```
# try create/delete pods, e.g.
$ kubectl create -f config/samples/pod-nginx1.yaml

# describe custom resource and check Status field
$ kubectl describe cpus.stats.example.org cpu-sample
Name:         cpu-sample
Namespace:    default
...
Status:
  Requests:
    Ng - Instance - 1:  100m
```
You can view a full demo video [here](doc/).
