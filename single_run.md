1. redis

```
docker run -d --restart=always --name=redis -p 36379:6379 -e REDIS_PASSWORD=VQ2WLo_URY bitnami/redis:latest
```


2. 启动

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
    -e SCANNER_REDIS_URL="redis://:VQ2WLo_URY@10.206.38.68:36379/1" \
    -e SCANNER_RULE_CHECKER_ADMIN_USERNAME=admin \
    -e SCANNER_RULE_CHECKER_ADMIN_PASSWORD="7_T^2nNFRT" \
    -e SCANNER_RULE_CHECKER_BASE_IMAGE_DIGESTS="c4c7334c2caba18f404262545f78ef8911e74b9334d852192ff9f225051fdb16" \
    -e SCANNER_RULE_CHECKER_IMAGE_LABELS="78b72b3a80deaae8b73474934b74bba16da5460dcb4a5c7a67f29f9a917dcfac" \
    registry.cn-beijing.aliyuncs.com/llaoj/harbor-scanner-trivy:rule-v0.30.2-0.3
```

3. 定期下载db

