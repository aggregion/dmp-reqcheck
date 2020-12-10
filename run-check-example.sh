agg_proxy_host="5.53.126.108"
agg_proxy_http="$agg_proxy_host:80"

http_proxy="http://$agg_proxy_http" \
https_proxy="http://$agg_proxy_http" \
DRC_CONCURRENCY=6 \
DRC_COMMON_PROXY="$agg_proxy_http" \
DRC_COMMON_CASTEST="$agg_proxy_host:28765" \
DRC_COMMON_CASPROD="$agg_proxy_host:18765" \
./bin/dmp-reqcheck check --roles $1 --hosts $2
