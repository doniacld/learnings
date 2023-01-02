# Getting started with Cilium

This page is a recap of my notes while doing the [Getting started with Cilium](https://isovalent.com/labs/getting-started-with-cilium/) lab from Isovalent.

## 🍇 Cluster environment

```shell
cat /etc/kind/${KIND_CONFIG}.yaml
```

```yaml
---
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane # control-plane node running the Kubernetes control plane and etcd
  extraPortMappings:
  # localhost.run proxy
  - containerPort: 32042
    hostPort: 32042
  # Hubble relay
  - containerPort: 31234
    hostPort: 31234
  # Hubble UI
  - containerPort: 31235
    hostPort: 31235
  extraMounts:
  - hostPath: /opt/images
    containerPath: /opt/images
- role: worker # worker nodes to deploy the applications
  extraMounts:
  - hostPath: /opt/images
    containerPath: /opt/images
- role: worker
  extraMounts:
  - hostPath: /opt/images
    containerPath: /opt/images
networking:
  disableDefaultCNI: true # the default CNI has been disabled so the cluster won't have any Pod network when it starts. Instead, Cilium is being deployed to the cluster to provide this functionality.
```

## ☎️ Cilium CLI

The cilium CLI tool  provides functionalities to install and 
check the status of Cilium in the cluster. It allows installing 
and updating Cilium on a cluster, as well as activate features such as Hubble and Cluster Mesh.

Command to install Cilium:
```shell
cilium install
```

Output of the installation command:
```shell
root@server:~# cilium  install
🔮 Auto-detected Kubernetes kind: kind
✨ Running "kind" validation checks
✅ Detected kind version "0.17.0"
ℹ️  Using Cilium version 1.12.2
🔮 Auto-detected cluster name: kind-kind
🔮 Auto-detected datapath mode: tunnel
🔮 Auto-detected kube-proxy has been installed
ℹ️  helm template --namespace kube-system cilium cilium/cilium --version 1.12.2 --set cluster.id=0,cluster.name=kind-kind,encryption.nodeEncryption=false,ipam.mode=kubernetes,kubeProxyReplacement=disabled,operator.replicas=1,serviceAccounts.cilium.name=cilium,serviceAccounts.operator.name=cilium-operator,tunnel=vxlan
ℹ️  Storing helm values file in kube-system/cilium-cli-helm-values Secret
🔑 Created CA in secret cilium-ca
🔑 Generating certificates for Hubble...
🚀 Creating Service accounts...
🚀 Creating Cluster roles...
🚀 Creating ConfigMap for Cilium version 1.12.2...
🚀 Creating Agent DaemonSet...
🚀 Creating Operator Deployment...
⌛ Waiting for Cilium to be installed and ready...
✅ Cilium was successfully installed! Run 'cilium status' to view installation health
```

Command to check Cilium status:
```shell
cilium status
```

Output of the status command:
```bash
root@server:~# cilium status
    /¯¯\
 /¯¯\__/¯¯\    Cilium:         OK
 \__/¯¯\__/    Operator:       OK
 /¯¯\__/¯¯\    Hubble:         disabled
 \__/¯¯\__/    ClusterMesh:    disabled
    \__/

Deployment        cilium-operator    Desired: 1, Ready: 1/1, Available: 1/1
DaemonSet         cilium             Desired: 3, Ready: 3/3, Available: 3/3
Containers:       cilium             Running: 3
                  cilium-operator    Running: 1
Cluster Pods:     3/3 managed by Cilium
Image versions    cilium             quay.io/cilium/cilium:v1.12.2@sha256:986f8b04cfdb35cf714701e58e35da0ee63da2b8a048ab596ccb49de58d5ba36: 3
                  cilium-operator    quay.io/cilium/operator-generic:v1.12.2@sha256:00508f78dae5412161fa40ee30069c2802aef20f7bdd20e91423103ba8c0df6e: 1
```

## 👮🏽 Network policies 

```shell
root@server:~# kubectl rollout status -n kube-system daemonset/cilium
daemon set "cilium" successfully rolled out
```

Deploy the following file:
```shell
root@server:~# kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/HEAD/examples/minikube/http-sw-app.yaml
service/deathstar created
deployment.apps/deathstar created
pod/tiefighter created
pod/xwing created
```

Verify the deployment:
```shell
root@server:~# kubectl get pods,svc
NAME                            READY   STATUS    RESTARTS   AGE
pod/deathstar-f694cf746-99v94   1/1     Running   0          55s
pod/deathstar-f694cf746-k9bxs   1/1     Running   0          56s
pod/tiefighter                  1/1     Running   0          56s
pod/xwing                       1/1     Running   0          55s

NAME                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
service/deathstar    ClusterIP   10.96.130.208   <none>        80/TCP    56s
service/kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP   21m
```

Each pod will also be represented in Cilium as an Endpoint. To retrieve a list of all endpoints managed by Cilium, the Cilium Endpoint (or `cep`) resource can be used:

```shell
root@server:~# kubectl get cep --all-namespaces
NAMESPACE            NAME                                      ENDPOINT ID   IDENTITY ID   INGRESS ENFORCEMENT   EGRESS ENFORCEMENT   VISIBILITY POLICY   ENDPOINT STATE   IPV4           IPV6
default              deathstar-f694cf746-99v94                 1350          13153         <status disabled>     <status disabled>    <status disabled>   ready            10.244.2.80    
default              deathstar-f694cf746-k9bxs                 587           13153         <status disabled>     <status disabled>    <status disabled>   ready            10.244.1.243   
default              tiefighter                                955           4234          <status disabled>     <status disabled>    <status disabled>   ready            10.244.2.21    
default              xwing                                     242           16321         <status disabled>     <status disabled>    <status disabled>   ready            10.244.1.46    
kube-system          coredns-6d4b75cb6d-757h2                  3161          22974         <status disabled>     <status disabled>    <status disabled>   ready            10.244.2.7     
kube-system          coredns-6d4b75cb6d-xl5c4                  2832          876           <status disabled>     <status disabled>    <status disabled>   ready            10.244.1.72    
local-path-storage   local-path-provisioner-6b84c5c67f-q645g   832           1047          <status disabled>     <status disabled>    <status disabled>   ready            10.244.1.216   
root@server:~# 
```

## 🪪 Identities & Cloud Native

IP addresses are no longer relevant for Cloud Native workloads. Security policies need something else.
Cilium provides this: Cilium uses the labels assigned to pods to define security policies.

Example of labels: `org=empire` and `class=deathstar`  
This is a simple policy that filters only on network layer 3 (IP protocol) and network layer 4 (TCP protocol), 
so it is often referred to as a L3/L4 network security policy.`

```shell
kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/HEAD/examples/minikube/sw_l3_l4_policy.yaml
```

https://raw.githubusercontent.com/cilium/cilium/HEAD/examples/minikube/sw_l3_l4_policy.yaml
```shell
apiVersion: "cilium.io/v2"
kind: CiliumNetworkPolicy
metadata:
  name: "rule1"
spec:
  description: "L3-L4 policy to restrict deathstar access to empire ships only"
  endpointSelector:
    matchLabels:
      org: empire # matching labels to restrict network
      class: deathstar
  ingress:
  - fromEndpoints:
    - matchLabels:
        org: empire
    toPorts:
    - ports:
      - port: "80"
        protocol: TCP
```

## 🛡️ Tighter Rules

To restrict the accessible paths from a pod to another, we can add a new rule:

https://raw.githubusercontent.com/cilium/cilium/HEAD/examples/minikube/sw_l3_l4_l7_policy.yaml
```shell
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: rule1
spec:
  endpointSelector:
    matchLabels:
      org: empire
      class: deathstar
  ingress:
    - fromEndpoints:
        - matchLabels:
            org: empire
      toPorts:
        - ports:
            - port: "80"
              protocol: TCP
          rules: # Adding a new rule to restrict the accessible paths.
            http:
              - method: POST
                path: /v1/request-landing
```

Apply the new configuration:
```shell
root@server:~# kubectl apply -f https://raw.githubusercontent.com/cilium/cilium/HEAD/examples/minikube/sw_l3_l4_l7_policy.yaml
ciliumnetworkpolicy.cilium.io/rule1 configured
```

Calling the restricted path before the new rule on L7:
```shell
root@server:~# kubectl exec tiefighter -- curl -s -XPUT deathstar.default.svc.cluster.local/v1/exhaust-port
Panic: deathstar exploded

goroutine 1 [running]:
main.HandleGarbage(0x2080c3f50, 0x2, 0x4, 0x425c0, 0x5, 0xa)
        /code/src/github.com/empire/deathstar/
        temp/main.go:9 +0x64
main.main()
        /code/src/github.com/empire/deathstar/
        temp/main.go:5 +0x85
```

After applying the manifest:
```shell
root@server:~# kubectl exec tiefighter -- curl -s -XPUT deathstar.default.svc.cluster.local/v1/exhaust-port
Access denied
```

With **Cilium L7 security policies**, we are able to restrict tiefighter's access only the required API resources on deathstar, 
thereby implementing a “least privilege” security approach for communication between microservices.

## Challenge

Ship `tiefighter` can land on `deathstar` but not `xwing` with a restriction on L3/L4:

```shell
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: rule1
spec:
  endpointSelector:
    matchLabels:
      org: empire
      class: deathstar # accessing to deathstar
  ingress:
    - fromEndpoints:
        - matchLabels:
            org: empire
            class: tiefighter # but only tiefighter
      toPorts:
        - ports:
            - port: "80"
              protocol: TCP

```

---

## Summary

**Cilium CLI first commands**
```shell
cilium install
cilium status
```

**Cilium Endpoint**
cep = Cilium endpoint = a pod in Cilium

Retrieve all endpoints:
```shell
kubectl get cep --all-namespaces
```

**Network Policies**
- Network Policies can block or allow traffic between pods
- CRD CiliumNetworkPolicy
- Restrict at L3/L4 level
- L7 Network Policies can filter on HTTP paths
- Cilium supports standard Kubernetes Network Policies


```yaml
apiVersion: cilium.io/v2
kind: CiliumNetworkPolicy
metadata:
  name: "rule1"
spec:
description: "L3-L4 policy to restrict deathstar access to empire ships only"
  endpointSelector:
    matchLabels:
      org: empire
      class: deathstar
  ingress:
    - fromEndpoints:
        - matchLabels:
            org: empire
      toPorts:
        - ports:
            - port: "80"
              protocol: TCP
```
```yaml
          rules: # L7 policy to restrict the accessible paths 
            http:
              - method: POST
                path: /v1/request-landing
```