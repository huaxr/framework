// Author: XinRui Hua
// Time:   2023/02/02 12:46
// Git:    huaxr

package degradeAgent

// 自动降级、手动降级agent
//存储放火：包括Redis、MongoDB、ES、MySQL等错误量和耗时的放火。
//
//依赖服务放火：包括下游调用超时、抛异常等
//
//消息放火：主要是针对MQ的放火，包括减少Consumer、增加Producer发送耗时等。
//
//服务器放火：主要包括CPU、内存等服务器资源放火。
//
//网络放火：主要包括网络带宽、网络节点的放火，例如前面提到的盲测。
