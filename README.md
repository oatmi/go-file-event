# gofileevent

## 期望的使用模式
1. 初始化 client，需要传入文件路径或者目录路径
2. 开始订阅 Event，可以传入感兴趣的event 类型；
   1. 得到一个Subscription，可用于 err 判断或者终止订阅；
   2. 从 chan 中得到 Event 数据，进行后续的处理；
3. 支持多个订阅者