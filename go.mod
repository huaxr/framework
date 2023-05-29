module github.com/huaxr/framework

go 1.16

// export GOPROXY=goproxy.io

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4

replace github.com/Shopify/sarama => github.com/Shopify/sarama v1.25.0

require (
	cloud.google.com/go/bigquery v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/storage/azblob v0.3.0 // indirect
	github.com/Shopify/sarama v1.25.0
	github.com/VictoriaMetrics/fastcache v1.6.0 // indirect
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/agiledragon/gomonkey v2.0.2+incompatible
	github.com/alecthomas/units v0.0.0-20190924025748-f65c72e2690d // indirect
	github.com/alibaba/sentinel-golang v1.0.2
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5 // indirect
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a // indirect
	github.com/aws/aws-sdk-go-v2/config v1.1.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.1.1 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/bwmarrin/snowflake v0.3.0
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cloudflare/cloudflare-go v0.14.0 // indirect
	github.com/consensys/gnark-crypto v0.4.1-0.20210426202927-39ac3d4b3f1f // indirect
	github.com/coreos/bbolt v0.0.0-00010101000000-000000000000 // indirect
	github.com/coreos/etcd v3.3.27+incompatible
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/deepmap/oapi-codegen v1.8.2 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/docker/docker v1.6.2 // indirect
	github.com/docopt/docopt-go v0.0.0-20180111231733-ee0de3bc6815 // indirect
	github.com/dop251/goja v0.0.0-20220405120441-9037c2b61cbf // indirect
	github.com/edwingeng/doublejump v1.0.1
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/ethereum/go-ethereum v1.9.20
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/fjl/gencodec v0.0.0-20220412091415-8bb9e558978c // indirect
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/getkin/kin-openapi v0.76.0 // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gin-contrib/pprof v1.2.1
	github.com/gin-gonic/gin v1.5.0
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20200222043503-6f7a984d4dc4 // indirect
	github.com/go-kit/log v0.2.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-openapi/jsonreference v0.19.3 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/go-resty/resty/v2 v2.7.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/go-task/slim-sprig v0.0.0-20210107165309-348f09dbbbc0 // indirect
	github.com/goccy/go-json v0.9.7 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.3.0 // indirect
	github.com/golang/mock v1.4.4 // indirect
	github.com/golang/protobuf v1.3.2
	github.com/google/gnostic v0.3.0 // indirect
	github.com/google/gofuzz v1.1.1-0.20200604201612-c04b05f3adfa // indirect
	github.com/google/martian/v3 v3.0.0 // indirect
	github.com/google/pprof v0.0.0-20210407192527-94a9f03dee38 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/graph-gophers/graphql-go v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0 // indirect
	github.com/hashicorp/go-bexpr v0.1.10 // indirect
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/holiman/uint256 v1.2.0 // indirect
	github.com/huin/goupnp v1.0.3 // indirect
	github.com/influxdata/influxdb v1.8.3 // indirect
	github.com/influxdata/influxdb-client-go/v2 v2.4.0 // indirect
	github.com/influxdata/line-protocol v0.0.0-20210311194329-9aa0e372d097 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/jedisct1/go-minisign v0.0.0-20190909160543-45766022959e // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/karalabe/usb v0.0.2 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/nacos-group/nacos-sdk-go v1.0.8
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/nsqio/go-nsq v1.1.0
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/onsi/ginkgo v1.12.1 // indirect
	github.com/onsi/gomega v1.10.0 // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/prometheus/client_golang v1.5.1
	github.com/prometheus/common v0.10.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/shirou/gopsutil/v3 v3.21.6 // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/cast v1.5.0
	github.com/stoewer/go-strcase v1.2.0 // indirect
	github.com/stretchr/testify v1.8.0
	github.com/supranational/blst v0.3.8-0.20220526154634-513d2456b344 // indirect
	github.com/ugorji/go v1.2.7 // indirect
	github.com/urfave/cli/v2 v2.10.2 // indirect
	github.com/valyala/fastrand v1.1.0
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opencensus.io v0.22.4 // indirect
	go.uber.org/automaxprocs v1.5.1
	go.uber.org/zap v1.23.0
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/time v0.1.0
	golang.org/x/tools v0.1.12 // indirect
	golang.org/x/xerrors v0.0.0-20220517211312-f3a8303e98df // indirect
	google.golang.org/api v0.22.0 // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/grpc v1.31.0
	google.golang.org/protobuf v1.26.0-rc.1 // indirect
	k8s.io/apimachinery v0.19.0-alpha.1
	k8s.io/gengo v0.0.0-20210813121822-485abfe95c7c // indirect
	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	rsc.io/quote/v3 v3.1.0 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	xorm.io/xorm v1.3.2
)
