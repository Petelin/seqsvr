# seqsvr

## 全局唯一序列号生成服务

- 基于微信亿万级序列号生成器理论
- 为每一个entityID生成单调递增的seqID
- 千万qps(理论上可以无限扩展)
- 节点故障自动切流量
- 支持手动一键切流量，动态增加减少机器
- 可选客户端/服务端路由表，整个架构可以放在微服务下面


## 性能，单机

全部部署台同一台机器上，5秒统计值。 
```
>metrics: 21:35:24.940022 histogram alloc:FetchNextSeqNum
>metrics: 21:35:24.940033   count:          253593
>metrics: 21:35:24.940036   min:               147
>metrics: 21:35:24.940039   max:               306
>metrics: 21:35:24.940067   mean:              202.80
>metrics: 21:35:24.940074   stddev:             59.49
>metrics: 21:35:24.940079   median:            166.00
>metrics: 21:35:24.940082   75%:               269.50
>metrics: 21:35:24.940098   95%:               306.00
>metrics: 21:35:24.940102   99%:               306.00
>metrics: 21:35:24.940104   99.9%:             306.00
```
