# Sophie部署指南
设置环境变量

## https服务创建
1. 安装cfssl工具集
```
cfssl: 证书签发工具
cfssljson: 将cfssl生成的证书（json格式）变为文件承载式证书
```

执行命令:
```
cd ${SOPHIE_ROOT_DIR}/configs
mkdir cert
cfssl gencert -initca cfssl/ca-csr.json | cfssljson -bare ca 
ls ca*
mv ca* cert
```
生成ca-key.pem(私钥) 和 ca.pem(公钥)

2. 修改hosts
```
sudo tee -a /etc/hosts << EOF
127.0.0.1 sophie.gateway.com
127.0.0.1 sophie.system.com
EOF
```

3. 配置gateway
```
cd $SOPHIE_ROOT/configs
sudo mkdir /var/run/sophie
cfssl gencert -ca=cert/ca.pem -ca-key=cert/ca-key.pem -config=cfssl/ca-config.json \
-profile=sophie cfssl/sophie-gateway-csr.json | cfssljson -bare sophie-gateway
sudo cp sophie-gateway*pem /var/run/sophie/
mv sophie-gateway*pem cert
mv sophie-gateway*csr cert
```