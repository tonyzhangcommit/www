单台CentOS服务器可以部署一个单节点的Kubernetes集群（通常称为“单节点集群”或“单节点K8s”），这种配置主要用于测试和开发环境。下面是基本步骤：

### 1. 安装依赖

首先，确保系统是最新的：

```bash
sudo yum update -y
```

安装所需的依赖项：

```bash
sudo yum install -y yum-utils device-mapper-persistent-data lvm2
```

### 2. 安装Docker

Kubernetes需要Docker作为容器运行时。安装Docker：

```bash
sudo yum install -y docker
sudo systemctl enable docker
sudo systemctl start docker
```

### 3. 添加Kubernetes Yum仓库

创建Kubernetes的Yum仓库配置文件：

```bash
cat <<EOF | sudo tee /etc/yum.repos.d/kubernetes.repo
[kubernetes]
name=Kubernetes
baseurl=https://packages.cloud.google.com/yum/repos/kubernetes-el7-\$basearch
enabled=1
gpgcheck=1
repo_gpgcheck=1
gpgkey=https://packages.cloud.google.com/yum/doc/yum-key.gpg https://packages.cloud.google.com/yum/doc/rpm-package-key.gpg
EOF
```

### 4. 安装Kubeadm, Kubelet和Kubectl

```bash
sudo yum install -y kubelet kubeadm kubectl
sudo systemctl enable kubelet
```

### 5. 初始化单节点集群

使用Kubeadm初始化Kubernetes集群：

```bash
sudo kubeadm init --pod-network-cidr=10.244.0.0/16
```

### 6. 配置Kubectl

为了使用kubectl管理集群，配置kubectl：

```bash
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

### 7. 部署Pod网络插件

使用Flannel作为网络插件：

```bash
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/master/Documentation/kube-flannel.yml
```

### 8. 允许调度Pod到主节点

默认情况下，Kubernetes不允许在主节点上调度Pod。运行以下命令来移除这个限制：

```bash
kubectl taint nodes --all node-role.kubernetes.io/master-
```

现在，你已经成功在单台CentOS服务器上部署了一个单节点的Kubernetes集群。你可以使用`kubectl`命令来管理你的集群。

如果你有任何进一步的问题或需要更多帮助，请随时告诉我。

