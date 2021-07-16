# Kata Containers v2.1 源码分析
[代码仓库](https://github.com/kata-containers/kata-containers)

kata 启动容器命令
```sh
/usr/local/bin/containerd-shim-kata-v2 -namespace k8s.io -address /run/containerd/containerd.sock -publish-binary /usr/bin/containerd -id 41153ea62c0b0528e58eae5a345ce39e2e29c48b22a5664d2440313fde10f45c -debug
```