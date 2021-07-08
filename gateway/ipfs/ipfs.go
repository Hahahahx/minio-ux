package ipfs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/minio/cli"
	"github.com/minio/madmin-go"
	minio "github.com/minio/minio/cmd"

	shell "github.com/ipfs/go-ipfs-api"
)

// be used in cli.Command
const (
	ipfsBackend = "ipfs"
)

// IPFS implements Gateway{Name,NewGatewayLayer}.
type IPFS struct {
	// need your IPFS host address and port.
	host string
}

// implements gateway for MinIO and S3 compatible object storage servers.
type ipfsObjects struct {
	minio.GatewayUnsupported
	ipfs *shell.Shell
}

func init() {

	const ipfsGatewayTemplate = `NAME:
	{{.HelpName}} - {{.Usage}}
  
  USAGE:
	{{.HelpName}} {{if .VisibleFlags}}[FLAGS]{{end}} PATH
  {{if .VisibleFlags}}
  FLAGS:
	{{range .VisibleFlags}}{{.}}
	{{end}}{{end}}
  PATH:
	path to NAS mount point
  
  EXAMPLES:
	1. Start minio gateway server for IPFS backend
	   {{.Prompt}} {{.EnvVarSetCommand}} MINIO_ROOT_USER{{.AssignmentOperator}}accesskey
	   {{.Prompt}} {{.EnvVarSetCommand}} MINIO_ROOT_PASSWORD{{.AssignmentOperator}}secretkey
	   {{.Prompt}} {{.HelpName}} /shared/nasvol
  
	2. Start minio gateway server for IPFS backend with edge caching enabled
	   {{.Prompt}} {{.EnvVarSetCommand}} MINIO_ROOT_USER{{.AssignmentOperator}}accesskey
	   {{.Prompt}} {{.EnvVarSetCommand}} MINIO_ROOT_PASSWORD{{.AssignmentOperator}}secretkey
	   {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_DRIVES{{.AssignmentOperator}}"/mnt/drive1,/mnt/drive2,/mnt/drive3,/mnt/drive4"
	   {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_EXCLUDE{{.AssignmentOperator}}"bucket1/*,*.png"
	   {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_QUOTA{{.AssignmentOperator}}90
	   {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_AFTER{{.AssignmentOperator}}3
	   {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_WATERMARK_LOW{{.AssignmentOperator}}75
	   {{.Prompt}} {{.EnvVarSetCommand}} MINIO_CACHE_WATERMARK_HIGH{{.AssignmentOperator}}85
	   {{.Prompt}} {{.HelpName}} /shared/nasvol
  `

	if err := minio.RegisterGatewayCommand(cli.Command{
		Name:               ipfsBackend,
		Usage:              "InterPlanetary File System(IPFS)",
		Action:             ipfsGatewayMain,
		CustomHelpTemplate: ipfsGatewayTemplate,
		HideHelpCommand:    true,
	}); err != nil {
		panic(err)
	}
}

func ipfsGatewayMain(ctx *cli.Context) {

	if ctx.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(ctx, ipfsBackend, 1)
	}

	minio.StartGateway(ctx, &IPFS{host: ctx.Args().First()})
}

// Name implements Gateway interface.
func (g *IPFS) Name() string {
	return ipfsBackend
}

// NewGatewayLayer returns ipfs gatewaylayer (instance of ipfsObject)
// func (g *IPFS) NewGatewayLayer(creds madmin.Credentials) (*ipfsObjects, error) {
func (g *IPFS) NewGatewayLayer(creds madmin.Credentials) (minio.ObjectLayer, error) {

	// Where your local node is running on localhost:5001
	host := If(g.host != "", g.host, "localhost:5001").(string)
	sh := shell.NewShell(host)
	out, err := sh.ID()
	if err != nil {
		fmt.Println("IPFS connected failure " + host)
		return nil, err
	}

	fmt.Printf("\nIPFS is run in: " + host)
	fmt.Printf("\nIpfs.ID: %s\n", out.ID)

	return &ipfsObjects{
		ipfs: sh,
	}, nil
}

func (i *ipfsObjects) Shutdown(ctx context.Context) error {
	return nil
}

type Repo struct {
	NumObjects uint
	RepoSize   uint64
	StorageMax uint64
	RepoPath   string
	Version    string
}

func (i *ipfsObjects) StorageInfo(ctx context.Context) (si minio.StorageInfo, err []error) {

	var out Repo

	if error := i.ipfs.Request("stats/repo").Exec(context.Background(), &out); error != nil {
		return
	}

	si.Disks = []madmin.Disk{{
		UsedSpace:  out.RepoSize,
		TotalSpace: out.StorageMax,
		DrivePath:  out.RepoPath,
	}}

	si.Backend.Type = madmin.Gateway
	si.Backend.GatewayOnline = true

	return si, nil
}

func (i *ipfsObjects) LocalStorageInfo(ctx context.Context) (si minio.StorageInfo, err []error) {
	return i.StorageInfo(ctx)
}

func (i *ipfsObjects) MakeBucketWithLocation(ctx context.Context, bucket string, opts minio.BucketOptions) (err error) {
	i.ipfs.AddDir(bucket)
	return
}

func (i *ipfsObjects) GetBucketInfo(ctx context.Context, bucket string) (bucketInfo minio.BucketInfo, err error) {
	return
}
func (i *ipfsObjects) ListBuckets(ctx context.Context) (buckets []minio.BucketInfo, err error) {
	return
}

func (i *ipfsObjects) DeleteBucket(ctx context.Context, bucket string, forceDelete bool) (err error) {
	return
}

func (i *ipfsObjects) ListObjects(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result minio.ListObjectsInfo, err error) {
	return
}

func (i *ipfsObjects) GetObjectNInfo(ctx context.Context, bucket, object string, rs *minio.HTTPRangeSpec, h http.Header, lockType minio.LockType, opts minio.ObjectOptions) (reader *minio.GetObjectReader, err error) {
	return
}
func (i *ipfsObjects) GetObjectInfo(ctx context.Context, bucket, object string, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
	return
}
func (i *ipfsObjects) PutObject(ctx context.Context, bucket, object string, data *minio.PutObjReader, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
	return
}
func (i *ipfsObjects) DeleteObject(ctx context.Context, bucket, object string, opts minio.ObjectOptions) (objInfo minio.ObjectInfo, err error) {
	return
}
func (i *ipfsObjects) DeleteObjects(ctx context.Context, bucket string, objects []minio.ObjectToDelete, opts minio.ObjectOptions) (delObject []minio.DeletedObject, err []error) {
	return
}
