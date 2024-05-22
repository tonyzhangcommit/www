#### K8s相关文档
1. k8s 是什么？有什么作用？
Kubernetes（简称K8s）是一个用于自动化部署、扩展和管理容器化应用程序的开源平台。它最初由Google设计，并于2014年开源，现在由云原生计算基金会（CNCF）进行维护。

### Kubernetes的主要作用包括：

1. **自动化部署和回滚**：Kubernetes可以自动化部署应用程序，并在需要时进行回滚。
2. **弹性伸缩**：根据负载情况，自动调整应用实例的数量，以满足应用的需求。
3. **服务发现和负载均衡**：Kubernetes提供内置的服务发现和负载均衡功能，使得应用可以通过DNS名称或IP地址访问。
4. **存储编排**：Kubernetes可以自动挂载存储系统，如本地存储、云提供商的存储（如AWS、GCP等），以及网络存储（如NFS）。
5. **自我修复**：自动重新启动失败的容器、替换被杀死的容器、杀死不响应用户定义的健康检查的容器等。
6. **配置管理和机密管理**：管理应用程序的配置和机密信息（如密码、OAuth令牌等），并在容器中使用它们而不需要重建容器镜像。
7. **容器编排**：管理和调度容器的运行、监控容器状态，并通过API控制和管理容器的生命周期。

### 重点概念 Node pod container,pool 之间的关系
在Kubernetes（K8s）中，节点（Node）、Pod、容器（Container）和资源池（Pool）之间有明确的层次结构和关系。以下是它们之间的关系和作用：

### 节点（Node）
- **定义**：节点是Kubernetes集群中的一台虚拟机或物理机，它是Kubernetes运行时的工作单元。
- **作用**：每个节点上运行着容器运行时（如Docker）、Kubelet（负责与Kubernetes主控平面通信）和Kube-proxy（处理网络路由）。

### Pod
- **定义**：Pod是Kubernetes中的最小可调度单元，一个Pod可以包含一个或多个紧密耦合的容器。
- **作用**：Pod中的容器共享同一个网络命名空间和存储卷，这使得它们能够相互通信和共享资源。Pod定义了应用程序的运行环境和逻辑单位。

### 容器（Container）
- **定义**：容器是应用程序及其依赖项的打包单元，它是通过容器化技术（如Docker）创建的轻量级、可移植的运行环境。
- **作用**：容器运行在Pod中，提供隔离的运行环境。每个容器包含应用程序代码、运行时、库和配置文件。

### 资源池（Resource Pool）
- **定义**：资源池在Kubernetes中并不是一个具体的概念，而是指集群中可用资源的集合。这些资源包括CPU、内存、存储等。
- **作用**：资源池是供Pod调度和运行时使用的资源总和。Kubernetes通过调度器将Pod分配到合适的节点上，以最优方式利用集群资源。

### 关系
1. **容器在Pod中**：一个Pod可以包含一个或多个容器，这些容器共享同一个网络和存储卷。
2. **Pod在节点上**：多个Pod可以部署在同一个节点上，节点提供Pod所需的计算资源。
3. **节点在集群中**：多个节点组成一个Kubernetes集群，节点通过Kubernetes主控平面进行管理和协调。
4. **资源池管理**：资源池由集群中的所有节点资源构成，Kubernetes调度器从资源池中分配资源给各个Pod。

通过上述关系，Kubernetes实现了容器化应用的自动化部署、伸缩和管理，使得应用能够高效、可靠地运行在分布式环境中。

#### k8s 中， ingress  VS  ingress Controller

在Kubernetes中，Ingress Controller和Ingress资源是用来管理外部访问集群内服务的流量的。它们的功能和作用如下：

### Ingress Controller
- **定义**：Ingress Controller是一个专门的负载均衡器，负责根据Ingress资源的配置，动态地配置HTTP和HTTPS路由规则，从而管理外部流量的进入。
- **作用**：Ingress Controller通常运行在Kubernetes集群中，它监控集群内的Ingress资源，根据这些资源定义的规则（如路径匹配、主机名匹配等），配置负载均衡器（如NGINX、HAProxy、Traefik等）来处理外部流量的路由。

### Ingress
- **定义**：Ingress是一种Kubernetes资源对象，用来定义如何将外部HTTP和HTTPS流量路由到集群内的服务。它可以配置基于域名和路径的路由规则，还可以配置SSL/TLS终止。
- **作用**：Ingress资源提供了一个灵活的方式来管理七层（应用层）的网络流量路由。它可以配置不同的路由规则来将请求转发到不同的服务。

### 四层代理 vs 七层代理
- **四层代理（L4 Proxy）**：指的是传输层的代理，处理TCP/UDP协议的流量转发。Kubernetes中的Service类型，如LoadBalancer和NodePort，能够实现四层代理。这种代理可以基于IP地址和端口进行流量转发，但无法基于HTTP请求的路径或主机名进行更细粒度的控制。
- **七层代理（L7 Proxy）**：指的是应用层的代理，处理HTTP/HTTPS协议的流量转发。Ingress资源可以实现七层代理。通过Ingress资源，用户可以定义基于HTTP请求的路径、主机名等信息的路由规则，从而实现更加细粒度的流量控制和管理。

### 总结
- **Ingress Controller**：负责监控和实现Ingress资源定义的七层代理规则。它能够动态配置集群内的负载均衡器，实现外部流量的智能路由。
- **Ingress**：是一种资源对象，用来定义七层代理规则，包括基于域名和路径的路由，以及SSL/TLS终止等功能。

通过使用Ingress和Ingress Controller，Kubernetes能够实现更高级的流量管理和控制，为集群内的服务提供灵活的外部访问方式。

在Kubernetes中，Service对象和Ingress对象是两种主要用于网络流量管理的资源，它们各自有不同的功能和作用。以下是对这两种对象的详细介绍：

### Service 对象
**定义**：Service是一种抽象的方式，用于定义一组Pod的逻辑集合以及如何访问这些Pod（通常是通过DNS名称）。

**功能和作用**：
1. **负载均衡**：Service对象可以在一组Pod之间实现负载均衡。当多个Pod提供相同的服务时，Service会将请求分发给这些Pod，确保流量均匀分布。
2. **服务发现**：Kubernetes为每个Service分配一个固定的IP地址和DNS名称。集群内的其他Pod可以通过这个IP地址或DNS名称来访问该Service。
3. **网络代理**：Service可以充当网络代理，将请求从一个稳定的端点转发到动态变化的Pod集合。Kubernetes中有多种类型的Service：
   - **ClusterIP**（默认）：服务只能在集群内部访问。
   - **NodePort**：服务可以通过集群每个节点的静态端口进行访问。
   - **LoadBalancer**：创建一个外部负载均衡器（如果在云环境中），将流量转发到Service。
   - **ExternalName**：将Service映射到一个外部的DNS名称。

**示例**：
```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: MyApp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
```
上面的示例定义了一个名为`my-service`的Service，它选择标签为`app: MyApp`的Pod，将流量从Service的端口80转发到Pod的端口8080。

### Ingress 对象
**定义**：Ingress是一种Kubernetes资源，用于管理外部HTTP和HTTPS流量的路由。它提供了七层（应用层）路由功能，可以基于HTTP请求的路径和主机名进行流量转发。

**功能和作用**：
1. **路径和主机名路由**：Ingress允许用户定义基于URL路径和主机名的路由规则，将外部请求转发到不同的Service。这样可以在一个IP地址和端口下管理多个服务。
2. **SSL/TLS终止**：Ingress可以配置SSL/TLS终止，使得HTTPS流量在Ingress处被解密，然后以HTTP的形式传递到集群内部的服务。
3. **增强的流量控制**：Ingress提供了更多高级的流量控制功能，如重定向、Rewrite、负载均衡策略等。

**示例**：
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
spec:
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /foo
        pathType: Prefix
        backend:
          service:
            name: foo-service
            port:
              number: 80
      - path: /bar
        pathType: Prefix
        backend:
          service:
            name: bar-service
            port:
              number: 80
```
上面的示例定义了一个名为`my-ingress`的Ingress，包含两个路径规则：访问`myapp.example.com/foo`时流量转发到`foo-service`，访问`myapp.example.com/bar`时流量转发到`bar-service`。

### 总结
- **Service对象**：主要用于集群内部的服务发现和负载均衡，能够通过ClusterIP、NodePort和LoadBalancer等类型实现不同的访问方式。
- **Ingress对象**：主要用于管理集群外部到内部的HTTP和HTTPS流量，提供基于路径和主机名的路由规则，以及SSL/TLS终止等高级功能。

通过结合使用Service和Ingress，Kubernetes能够提供灵活且强大的网络流量管理方案，满足各种应用场景的需求。

#### k8s中常用的对象
在Kubernetes中，有许多不同类型的对象，每种对象都有特定的功能和用途。以下是一些常用的Kubernetes对象及其功能和作用：

### Pod
- **定义**：Pod是Kubernetes中最小的可部署单元。一个Pod可以包含一个或多个紧密耦合的容器，通常共享同一个网络命名空间和存储卷。
- **作用**：Pod为容器提供了运行环境和资源隔离。

### ReplicaSet (RS)
- **定义**：ReplicaSet是用来确保某个Pod的指定数量的副本始终在运行。
- **作用**：通过维护指定数量的Pod副本来保证应用的高可用性和负载均衡。

### Deployment
- **定义**：Deployment是用来管理Pod和ReplicaSet的控制器，提供声明式的更新能力。
- **作用**：支持滚动更新和回滚，能够方便地进行应用版本的升级和降级。

### StatefulSet
- **定义**：StatefulSet是用来管理有状态应用的工作负载对象，确保Pod具有唯一性和顺序性。
- **作用**：适用于需要持久存储和有序部署、更新的应用，如数据库、分布式文件系统等。

### DaemonSet
- **定义**：DaemonSet确保所有（或某些）节点上运行一个Pod的副本。
- **作用**：常用于集群中每个节点上都需要运行的任务，如日志收集器、监控代理等。

### Job
- **定义**：Job是用来创建一次性任务的控制器，确保任务成功完成。
- **作用**：用于执行短期的、一次性的任务，如批处理任务。

### CronJob
- **定义**：CronJob是用来创建周期性任务的控制器，按照计划的时间点运行Job。
- **作用**：适用于需要定期运行的任务，如定时备份、定期报告等。

### Service
- **定义**：Service是一种抽象，定义一组Pod的逻辑集合以及如何访问这些Pod。
- **作用**：提供服务发现和负载均衡，能够通过ClusterIP、NodePort、LoadBalancer等方式暴露服务。

### Ingress
- **定义**：Ingress是用来管理外部访问集群内服务的HTTP和HTTPS流量的资源对象。
- **作用**：提供七层（应用层）的流量管理和控制，如基于路径和主机名的路由、SSL/TLS终止等。

### ConfigMap
- **定义**：ConfigMap是用来存储非机密数据的键值对。
- **作用**：用于配置数据的管理，能够在Pod中作为环境变量、命令行参数或配置文件使用。

### Secret
- **定义**：Secret是用来存储和管理敏感信息的对象，如密码、OAuth令牌、SSH密钥等。
- **作用**：安全地分发和使用敏感数据，能够在Pod中作为环境变量、命令行参数或文件使用。

### PersistentVolume (PV)
- **定义**：PersistentVolume是集群中由管理员提供的存储资源。
- **作用**：提供独立于Pod生命周期的持久存储。

### PersistentVolumeClaim (PVC)
- **定义**：PersistentVolumeClaim是用户对PersistentVolume的存储请求。
- **作用**：用户通过PVC请求PV中的存储资源，实现持久存储的绑定和使用。

### Namespace
- **定义**：Namespace是用来将集群资源划分为多个虚拟集群的机制。
- **作用**：提供资源的逻辑隔离和管理，适用于多租户环境。

### Role 和 RoleBinding
- **定义**：Role定义了在某个命名空间内的权限，RoleBinding将Role绑定到用户、组或服务账号上。
- **作用**：提供命名空间级别的权限控制。

### ClusterRole 和 ClusterRoleBinding
- **定义**：ClusterRole定义了集群范围内的权限，ClusterRoleBinding将ClusterRole绑定到用户、组或服务账号上。
- **作用**：提供集群级别的权限控制。

这些对象共同构成了Kubernetes的核心功能，帮助用户实现应用的部署、管理、扩展和维护。

#### k8s 对象和组件的区别
上述提到的对象是Kubernetes中的资源类型（也称为API对象或资源对象）。它们并不是Kubernetes的组件，而是用户在Kubernetes集群中定义和管理的配置实体。这些对象是以声明性配置的形式（通常是YAML或JSON文件）来描述的，Kubernetes根据这些配置来管理集群中的资源和工作负载。

这些对象在Kubernetes中以API资源的形式存在，通过Kubernetes API服务器进行管理和操作。用户可以使用`kubectl`命令行工具、Kubernetes Dashboard或直接通过Kubernetes API与这些资源进行交互。

### Kubernetes中的组件与对象的关系
- **API服务器（kube-apiserver）**：Kubernetes的核心组件，负责处理REST API请求，验证和配置数据对象，如Pod、Service等。
- **控制器管理器（kube-controller-manager）**：负责运行控制循环来管理和维护集群的状态，包括ReplicaSet、Deployment、Job等控制器。
- **调度器（kube-scheduler）**：负责将未指定节点的Pod分配到合适的节点上。
- **etcd**：分布式键值存储，用于存储所有Kubernetes集群的状态数据。
- **Kubelet**：运行在每个节点上，负责管理节点上的Pod和容器，确保它们按照定义的状态运行。
- **Kube-proxy**：负责维护节点网络规则，处理Service对象的负载均衡和网络代理功能。

### Kubernetes对象的声明性配置
这些Kubernetes对象以声明性配置的形式定义，通常使用YAML或JSON文件。以下是一些示例配置文件：

#### Pod 示例
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: my-pod
spec:
  containers:
  - name: my-container
    image: nginx
```

#### Service 示例
```yaml
apiVersion: v1
kind: Service
metadata:
  name: my-service
spec:
  selector:
    app: MyApp
  ports:
    - protocol: TCP
      port: 80
      targetPort: 8080
  type: ClusterIP
```

#### Deployment 示例
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: MyApp
  template:
    metadata:
      labels:
        app: MyApp
    spec:
      containers:
      - name: my-container
        image: nginx
```

#### Ingress 示例
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-ingress
spec:
  rules:
  - host: myapp.example.com
    http:
      paths:
      - path: /foo
        pathType: Prefix
        backend:
          service:
            name: foo-service
            port:
              number: 80
```

用户将这些配置文件应用到Kubernetes集群中，Kubernetes API服务器会解析这些文件并创建相应的对象。控制器管理器和其他组件会根据这些对象的定义来管理集群中的资源和工作负载，确保它们处于期望的状态。

### 总结
Kubernetes中的对象（如Pod、Service、Deployment等）是声明性配置的实体，用于定义和管理集群中的资源和工作负载。它们通过Kubernetes API进行管理，而Kubernetes组件（如API服务器、控制器管理器、调度器等）负责确保这些对象按预期运行和维护集群的状态。




#### 启动minnikube 命令     minikube start --vm-driver=docker --base-image="anjone/kicbase" --registry-mirror=https://lbhsgfub.mirror.aliyuncs.com