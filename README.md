# minio gateway

首先在我们的网关开发中，需要先实现一个 init()函数，该函数会被 minio 主进程调用
函数主要实现 `minio.RegisterGatewayCommand(cmd cli.Command)`

如下，是 nas 网关中的实现：

```Go

func init() {
	const nasGatewayTemplate = `NAME` // 忽略不重要的信息

	minio.RegisterGatewayCommand(cli.Command{
		Name:               minio.NASBackendGateway,    // minio定义的全局字符串 "nas"
		Usage:              "Network-attached storage (NAS)",
		Action:             nasGatewayMain,         // 主要执行函数
		CustomHelpTemplate: nasGatewayTemplate,
		HideHelpCommand:    true,
	})
}

```

他是提供给命令行工具的重要参数，比如 Name 则是网关的名字，Action 是调用该命令时执行的主要进程，也是最重要的地方。

在该函数中，我们会获得一个`ctx`，最重要的是需要调用`minio.StartGateway(ctx *cli.Context, gw minio.Gateway)`

```Go
func nasGatewayMain(ctx *cli.Context) {
	// 此处验证官网的参数，nas网关还要求后面跟一个path，也就是地址参数，比如：
    // minio gateway nas /mnt/minio/
	if !ctx.Args().Present() || ctx.Args().First() == "help" {
		cli.ShowCommandHelpAndExit(ctx, minio.NASBackendGateway, 1)
	}

    // 那么，path : "/mnt/minio"
	minio.StartGateway(ctx, &NAS{ctx.Args().First()})
}

```

在 nas 的代码中最终会给`StartGateway`传入一个结构体，该结构体是接口`minio.Gateway`的实现。

nas 中的结构体，以及 minio 所定义的相关接口：

```Go
// 该结构体需要实现minio.Gateway
type NAS struct {
	path string
}

type Gateway interface {
    // 只要返回前面的网关名字符串就好
	Name() string

	NewGatewayLayer(creds madmin.Credentials) (ObjectLayer, error)
}

// 用来存放权限信息
type Credentials struct {
	AccessKey    string    `xml:"AccessKeyId" json:"accessKey,omitempty"`
	SecretKey    string    `xml:"SecretAccessKey" json:"secretKey,omitempty"`
	SessionToken string    `xml:"SessionToken" json:"sessionToken,omitempty"`
	Expiration   time.Time `xml:"Expiration" json:"expiration,omitempty"`
}
```

在前面的`StartGateway`中，我们实现了`minio.Gateway`接口，`NewGatewayLayer`和前面的流程差不多，也需要实现 minio 所提供的另一个接口`minio.ObjectLayer`并返回它，同时也提供了一个参数`minio.Credentials`，该参数存放着 minio 启动时的操作权限，包括账号和密钥以及有效期等

```Go

func (g *NAS) NewGatewayLayer(creds madmin.Credentials) (minio.ObjectLayer, error) {
	var err error
    // ObjectLayer实例
	newObject, err := minio.NewFSObjectLayer(g.path)
	if err != nil {
		return nil, err
	}
	return &nasObjects{newObject}, nil
}

```

在 nas 中，ObjectLayer 实例由`minio.NewFSObjectLayer`创建，其用法与`minio server /mnt/minio`并无区别。而在我们自定义的网关中，并非这么实现，比如我们需要对接 ipfs，其中的诸多参数需要我们自己考量。

```Go
// FSObjectLayer实例
	fs := &FSObjects{
		fsPath:       fsPath,
		metaJSONFile: fsMetaJSONFile,
		fsUUID:       fsUUID,
		rwPool: &fsIOPool{
			readersMap: make(map[string]*lock.RLockedFile),
		},
		nsMutex:       newNSLock(false),
		listPool:      NewTreeWalkPool(globalLookupTimeout),
		appendFileMap: make(map[string]*fsAppendFile),
		diskMount:     mountinfo.IsLikelyMountPoint(fsPath),
	}

```

还记得我们需要实现`minio.ObjectLayer`接口吗。

```Go

// ObjectLayer implements primitives for object API layer.
type ObjectLayer interface {
	// Locking operations on object.
	NewNSLock(bucket string, objects ...string) RWLocker

	// Storage operations.
	Shutdown(context.Context) error
	NSScanner(ctx context.Context, bf *bloomFilter, updates chan<- madmin.DataUsageInfo) error

	BackendInfo() madmin.BackendInfo
	StorageInfo(ctx context.Context) (StorageInfo, []error)
	LocalStorageInfo(ctx context.Context) (StorageInfo, []error)

	// Bucket operations.
	MakeBucketWithLocation(ctx context.Context, bucket string, opts BucketOptions) error
	GetBucketInfo(ctx context.Context, bucket string) (bucketInfo BucketInfo, err error)
	ListBuckets(ctx context.Context) (buckets []BucketInfo, err error)
	DeleteBucket(ctx context.Context, bucket string, forceDelete bool) error
	ListObjects(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result ListObjectsInfo, err error)
	ListObjectsV2(ctx context.Context, bucket, prefix, continuationToken, delimiter string, maxKeys int, fetchOwner bool, startAfter string) (result ListObjectsV2Info, err error)
	ListObjectVersions(ctx context.Context, bucket, prefix, marker, versionMarker, delimiter string, maxKeys int) (result ListObjectVersionsInfo, err error)
	// Walk lists all objects including versions, delete markers.
	Walk(ctx context.Context, bucket, prefix string, results chan<- ObjectInfo, opts ObjectOptions) error

	// Object operations.

	// GetObjectNInfo returns a GetObjectReader that satisfies the
	// ReadCloser interface. The Close method unlocks the object
	// after reading, so it must always be called after usage.
	//
	// IMPORTANTLY, when implementations return err != nil, this
	// function MUST NOT return a non-nil ReadCloser.
	GetObjectNInfo(ctx context.Context, bucket, object string, rs *HTTPRangeSpec, h http.Header, lockType LockType, opts ObjectOptions) (reader *GetObjectReader, err error)
	GetObjectInfo(ctx context.Context, bucket, object string, opts ObjectOptions) (objInfo ObjectInfo, err error)
	PutObject(ctx context.Context, bucket, object string, data *PutObjReader, opts ObjectOptions) (objInfo ObjectInfo, err error)
	CopyObject(ctx context.Context, srcBucket, srcObject, destBucket, destObject string, srcInfo ObjectInfo, srcOpts, dstOpts ObjectOptions) (objInfo ObjectInfo, err error)
	DeleteObject(ctx context.Context, bucket, object string, opts ObjectOptions) (ObjectInfo, error)
	DeleteObjects(ctx context.Context, bucket string, objects []ObjectToDelete, opts ObjectOptions) ([]DeletedObject, []error)
	TransitionObject(ctx context.Context, bucket, object string, opts ObjectOptions) error
	RestoreTransitionedObject(ctx context.Context, bucket, object string, opts ObjectOptions) error

	// Multipart operations.
	ListMultipartUploads(ctx context.Context, bucket, prefix, keyMarker, uploadIDMarker, delimiter string, maxUploads int) (result ListMultipartsInfo, err error)
	NewMultipartUpload(ctx context.Context, bucket, object string, opts ObjectOptions) (uploadID string, err error)
	CopyObjectPart(ctx context.Context, srcBucket, srcObject, destBucket, destObject string, uploadID string, partID int,
		startOffset int64, length int64, srcInfo ObjectInfo, srcOpts, dstOpts ObjectOptions) (info PartInfo, err error)
	PutObjectPart(ctx context.Context, bucket, object, uploadID string, partID int, data *PutObjReader, opts ObjectOptions) (info PartInfo, err error)
	GetMultipartInfo(ctx context.Context, bucket, object, uploadID string, opts ObjectOptions) (info MultipartInfo, err error)
	ListObjectParts(ctx context.Context, bucket, object, uploadID string, partNumberMarker int, maxParts int, opts ObjectOptions) (result ListPartsInfo, err error)
	AbortMultipartUpload(ctx context.Context, bucket, object, uploadID string, opts ObjectOptions) error
	CompleteMultipartUpload(ctx context.Context, bucket, object, uploadID string, uploadedParts []CompletePart, opts ObjectOptions) (objInfo ObjectInfo, err error)

	// Policy operations
	SetBucketPolicy(context.Context, string, *policy.Policy) error
	GetBucketPolicy(context.Context, string) (*policy.Policy, error)
	DeleteBucketPolicy(context.Context, string) error

	// Supported operations check
	IsNotificationSupported() bool
	IsListenSupported() bool
	IsEncryptionSupported() bool
	IsTaggingSupported() bool
	IsCompressionSupported() bool

	SetDriveCounts() []int // list of erasure stripe size for each pool in order.

	// Healing operations.
	HealFormat(ctx context.Context, dryRun bool) (madmin.HealResultItem, error)
	HealBucket(ctx context.Context, bucket string, opts madmin.HealOpts) (madmin.HealResultItem, error)
	HealObject(ctx context.Context, bucket, object, versionID string, opts madmin.HealOpts) (madmin.HealResultItem, error)
	HealObjects(ctx context.Context, bucket, prefix string, opts madmin.HealOpts, fn HealObjectFn) error

	// Backend related metrics
	GetMetrics(ctx context.Context) (*BackendMetrics, error)

	// Returns health of the backend
	Health(ctx context.Context, opts HealthOptions) HealthResult
	ReadHealth(ctx context.Context) bool

	// Metadata operations
	PutObjectMetadata(context.Context, string, string, ObjectOptions) (ObjectInfo, error)

	// ObjectTagging operations
	PutObjectTags(context.Context, string, string, string, ObjectOptions) (ObjectInfo, error)
	GetObjectTags(context.Context, string, string, ObjectOptions) (*tags.Tags, error)
	DeleteObjectTags(context.Context, string, string, ObjectOptions) (ObjectInfo, error)
}
```

很多对吧，这很夸张，我知道，但是别急。nas中实现了三个函数`IsListenSupported`，`StorageInfo`，`IsTaggingSupported`


```Go


// IsListenSupported returns whether listen bucket notification is applicable for this gateway.
func (n *nasObjects) IsListenSupported() bool {
	return false
}

func (n *nasObjects) StorageInfo(ctx context.Context) (si minio.StorageInfo, _ []error) {
	si, errs := n.ObjectLayer.StorageInfo(ctx)
	si.Backend.GatewayOnline = si.Backend.Type == madmin.FS
	si.Backend.Type = madmin.Gateway
	return si, errs
}

func (n *nasObjects) IsTaggingSupported() bool {
	return true
}

```


据我了解nas网关提供的功能其实和server是一样的，在其他网关开发中和他还是有所区别的，比如s3、hdfs等。在上面的`ObjectLayer`诸多方法中，有许多的功能是无需实现的，在这里minio给我们提供了一个接口`minio.GatewayUnsupported`，它实现了一些在网关中不太重要的方法，好让我们在注册一个ObjectLayer时不需要自己手动实现这些方法。

以下是经过简化后的接口列表，也就是说这些方法也还是需要我们实现的，在s3、hdfs中也都完成了对他们的实现

```Go

// ObjectLayer implements primitives for object API layer.
type ObjectLayer interface {

	// Storage operations.
	// issue：关闭？
	Shutdown(context.Context) error
	// 存储信息
	StorageInfo(ctx context.Context) (StorageInfo, []error)

	// Bucket operations.
	// 创建一个Bucket
	MakeBucketWithLocation(ctx context.Context, bucket string, opts BucketOptions) error
	// 获取Bucket信息
	GetBucketInfo(ctx context.Context, bucket string) (bucketInfo BucketInfo, err error)
	ListBuckets(ctx context.Context) (buckets []BucketInfo, err error)
	DeleteBucket(ctx context.Context, bucket string, forceDelete bool) error
	ListObjects(ctx context.Context, bucket, prefix, marker, delimiter string, maxKeys int) (result ListObjectsInfo, err error)

	// Object operations.
	GetObjectNInfo(ctx context.Context, bucket, object string, rs *HTTPRangeSpec, h http.Header, lockType LockType, opts ObjectOptions) (reader *GetObjectReader, err error)
	GetObjectInfo(ctx context.Context, bucket, object string, opts ObjectOptions) (objInfo ObjectInfo, err error)
	PutObject(ctx context.Context, bucket, object string, data *PutObjReader, opts ObjectOptions) (objInfo ObjectInfo, err error)
	DeleteObject(ctx context.Context, bucket, object string, opts ObjectOptions) (ObjectInfo, error)
	DeleteObjects(ctx context.Context, bucket string, objects []ObjectToDelete, opts ObjectOptions) ([]DeletedObject, []error)
}
```

```