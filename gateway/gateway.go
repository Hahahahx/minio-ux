package gateway

// Package gateway imports ipfs gateway only.
// To enable all minio gateways please also import "github.com/minio/mino/cmd/gateway"

import (
	_ "github.com/Hahahahx/minio-ux/gateway/ipfs"
)
