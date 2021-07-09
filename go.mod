module github.com/Hahahahx/minio-ux

go 1.16

replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

require (
	github.com/fatih/color v1.12.0
	github.com/ipfs/go-ipfs-api v0.2.0
	github.com/minio/cli v1.22.0
	github.com/minio/madmin-go v1.0.15
	github.com/minio/minio v0.0.0-20210709055349-b6dd9b55a709
// github.com/fatih/color v1.12.0
// github.com/ipfs/go-ipfs-api v0.2.0
// github.com/minio/cli v1.22.0
// github.com/minio/madmin-go v1.0.13 // 只能是13版，目前升级成15会有编译问题，部分结构体发生变化
// github.com/minio/minio v0.0.0-20210708194325-84a64a7e479e
)
