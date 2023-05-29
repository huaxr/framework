// Author: huaxr
// Time: 2022-10-26 09:59
// Git: huaxr

package define

type ArchTag string

const (
	// normal log visible in debug console
	ArchError = "arch_error"
	ArchFatal = "arch_fatal"
	// system log which user should not take care of.
	HostEmpty          ArchTag = "arch_host_empty"
	ServiceContaminate ArchTag = "arch_sf_contaminate"
	ServiceDiff        ArchTag = "arch_sf_diff"
	EtcdPut            ArchTag = "arch_etcd_put"
	EtcdGet            ArchTag = "arch_etcd_get"
	EtcdDel            ArchTag = "arch_etcd_del"
	EtcdLeaseExpire    ArchTag = "arch_etcd_lease_expire"
	EtcdWatchPut       ArchTag = "arch_etcd_watch_put"
	EtcdWatchDel       ArchTag = "arch_etcd_watch_del"
	GinxQps            ArchTag = "arch_ginx_qps"
	GinxLimit          ArchTag = "arch_ginx_limit"
	GrpcQps            ArchTag = "arch_grpc_qps"
	GrpcClosed         ArchTag = "arch_grpc_closed"
	GrpcCircuitOpen    ArchTag = "arch_grpc_circuit"
	TcmUpdate          ArchTag = "arch_tcm_update"
)
