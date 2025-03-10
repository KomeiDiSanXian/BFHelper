module github.com/KomeiDiSanXian/BFHelper

go 1.22
toolchain go1.23.1

require (
	github.com/FloatTech/zbpctrl v1.7.0
	github.com/FloatTech/zbputils v1.7.1
	github.com/jinzhu/gorm v1.9.16
	github.com/sirupsen/logrus v1.9.3
	github.com/tidwall/gjson v1.18.0
	github.com/wdvxdr1123/ZeroBot v1.8.1
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.34.0
	go.opentelemetry.io/otel/sdk v1.34.0
	go.opentelemetry.io/otel/trace v1.34.0
)

require (
	github.com/RomiChan/websocket v1.4.3-0.20220227141055-9b2c6168c9c5 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/fumiama/terasu v0.0.0-20240507144117-547a591149c0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.25.1 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/sagikazarmark/locafero v0.6.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/proto/otlp v1.5.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240613232115-7f521ea00fb8 // indirect
	golang.org/x/net v0.34.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250115164207-1a7da9e5054f // indirect
	google.golang.org/grpc v1.69.4 // indirect
	google.golang.org/protobuf v1.36.3 // indirect
)

require (
	github.com/fsnotify/fsnotify v1.8.0
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/google/uuid v1.6.0
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	golang.org/x/sys v0.31.0
)

require (
	github.com/FloatTech/floatbox v0.0.0-20240505082030-226ec6713e14
	github.com/FloatTech/gg v1.1.3 // indirect
	github.com/FloatTech/imgfactory v0.2.2-0.20230315152233-49741fc994f9 // indirect
	github.com/FloatTech/rendercard v0.1.2 // indirect
	github.com/FloatTech/sqlite v1.7.0 // indirect
	github.com/FloatTech/ttl v0.0.0-20240716161252-965925764562 // indirect
	github.com/RomiChan/syncx v0.0.0-20240418144900-b7402ffdebc7 // indirect
	github.com/disintegration/imaging v1.6.2 // indirect
	github.com/ericpauley/go-quantize v0.0.0-20200331213906-ae555eb2afa4 // indirect
	github.com/fumiama/cron v1.3.0 // indirect
	github.com/fumiama/go-base16384 v1.7.0 // indirect
	github.com/fumiama/go-registry v0.2.7 // indirect
	github.com/fumiama/go-simple-protobuf v0.2.0 // indirect
	github.com/fumiama/gofastTEA v0.0.10 // indirect
	github.com/fumiama/imgsz v0.0.4 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pquerna/otp v1.4.0
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/spf13/viper v1.19.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.34.0
	golang.org/x/image v0.18.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	modernc.org/libc v1.61.0 // indirect
	modernc.org/mathutil v1.6.0 // indirect
	modernc.org/memory v1.8.0 // indirect
	modernc.org/sqlite v1.33.1 // indirect

)

replace modernc.org/sqlite => github.com/fumiama/sqlite3 v1.20.0-with-win386

replace github.com/FloatTech/zbputils => github.com/KomeiDiSanXian/zbputils v0.0.0-20230923095115-55ba2c51620d

replace github.com/remyoudompheng/bigfft => github.com/fumiama/bigfft v0.0.0-20211011143303-6e0bfa3c836b
