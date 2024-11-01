>  yaml文件常用管理方式
>
> * 开源工具，如：kustomize
> * git管理
> * 使用nginx代理


```shell
#利用nignx代理
#创建统一个文件夹，例如/data/k8s-yaml/nginx，后续的yaml文件分文件夹分类放置
#然后可以使用
#kubectl apply -f http://k8s-yaml.alnk.com/nginx/dp.yaml

#nginx配置
[root@ops nginx]# cat /etc/nginx/conf.d/k8s-yaml.od.com.conf 
server {
    listen       80;
    server_name  k8s-yaml.alnk.com;

    location / {
        autoindex on;
        default_type text/plain;
        root /data/k8s-yaml;
    }
}

kubectl apply -f http://k8s-yaml.alnk.com/nginx/dp.yaml

```
