1. Redis部署

trivy依赖redis运行, 可以采用redis中间件服务或者自建, 使用下面命了自建redis, 下面的配置是基于自建redis的.

```
docker run -d --restart=always --name=redis -p 36379:6379 -e REDIS_PASSWORD=<redis-password> bitnami/redis:latest
```

2. 启动

运行下面的命令启动trivy-scanner

```
docker stop harbor-scanner-trivy \
&& docker rm harbor-scanner-trivy \
&& docker run --restart=always -d --name=harbor-scanner-trivy \
    -p 38080:8080 \
    -u 10000:10000 \
    -v /data/trivy:/home/scanner/.cache \
    -e SCANNER_LOG_LEVEL=trace \
    -e SCANNER_TRIVY_DEBUG_MODE=true \
    -e SCANNER_TRIVY_INSECURE=true \
    -e SCANNER_TRIVY_OFFLINE_SCAN=true \
    -e SCANNER_TRIVY_SKIP_UPDATE=true \
    -e SCANNER_TRIVY_SEVERITY="HIGH,CRITICAL" \
    -e SCANNER_REDIS_URL="redis://:<redis-password>@10.206.38.68:36379/1" \
    -e SCANNER_RULE_CHECKER_ADMIN_USERNAME=<admin-user> \
    -e SCANNER_RULE_CHECKER_ADMIN_PASSWORD="<admin-password>" \
    -e SCANNER_RULE_CHECKER_BASE_IMAGE_DIGESTS="c4c7334c2caba18f404262545f78ef8911e74b9334d852192ff9f225051fdb16" \
    -e SCANNER_RULE_CHECKER_IMAGE_LABELS="78b72b3a80deaae8b73474934b74bba16da5460dcb4a5c7a67f29f9a917dcfac" \
    registry.cn-beijing.aliyuncs.com/llaoj/harbor-scanner-trivy:rule-v0.30.2-0.4
```

解释:

- scanner暴露38080端口
- `SCANNER_TRIVY_SKIP_UPDATE=true`和`SCANNER_TRIVY_SKIP_UPDATE=true`是因为都是调用国外服务, 大概率会失败.改镜像打包了最新trivy-db可以暂时不用更新.
- `SCANNER_REDIS_URL`配置redis地址, 格式: `redis://:<password>@<host>:<port>/<db-index>`
- `SCANNER_RULE_CHECKER_ADMIN_USERNAME`和`SCANNER_RULE_CHECKER_ADMIN_PASSWORD`为harbor管理员账号密码, scanner需要获取image的build history所以需要管理员授权.
- `SCANNER_RULE_CHECKER_BASE_IMAGE_DIGESTS`合规基础镜像字符串, 逗号分隔.
- `SCANNER_RULE_CHECKER_IMAGE_LABELS`合规label字符串, 逗号分隔.
- `registry.cn-beijing.aliyuncs.com/llaoj/harbor-scanner-trivy:rule-v0.30.2-0.4`scanner服务镜像, 建议上传到自有镜像仓库中.