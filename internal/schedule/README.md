# sophie schedule 

调度系统分为两个部分:manager和worker。manager 接受gateway 传递过来的job请求，负责对job进行同一的管理、worker节点状态的监控以及负载均衡
（可配置不同的负载均衡策略）；worker节点负责定时任务的具体执行。
