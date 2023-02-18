module github.com/artchitector/artchitect/soul

go 1.19

require (
	github.com/artchitector/artchitect/memory v0.0.0-20230206141224-ef4d2c479ec6
	github.com/artchitector/artchitect/model v0.0.0-20230218112449-15e526fcb934
	github.com/artchitector/artchitect/resizer v0.0.0-20230203133021-ba066d64422a
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-telegram/bot v0.5.1
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/joho/godotenv v1.4.0
	github.com/minio/minio-go/v7 v7.0.47
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.29.0
	golang.org/x/image v0.3.0
	gopkg.in/yaml.v3 v3.0.1
	gorm.io/driver/postgres v1.4.5
	gorm.io/gorm v1.24.5
)

require (
	github.com/artchitector/artchitect/bot v0.0.0-20230218165646-d26ddb6213b8 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.13.0 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.1 // indirect
	github.com/jackc/pgservicefile v0.0.0-20200714003250-2b9c44734f2b // indirect
	github.com/jackc/pgtype v1.13.0 // indirect
	github.com/jackc/pgx/v4 v4.17.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.9 // indirect
	github.com/klauspost/cpuid/v2 v2.1.0 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/minio/md5-simd v1.1.2 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646 // indirect
	github.com/rs/xid v1.4.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	golang.org/x/crypto v0.3.0 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.6.0 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
)

replace github.com/artchitector/artchitect/model => ../model

replace github.com/artchitector/artchitect/resizer => ../resizer

replace github.com/artchitector/artchitect/memory => ../memory
replace github.com/artchitector/artchitect/bot => ../bot
