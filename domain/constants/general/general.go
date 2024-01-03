package general

const (
	SourceFromDB    string = "db"
	SourceFromCache string = "cache"
)

// List database name
const (
	DatabaseRead  = "read"
	DatabaseWrite = "write"
)

const (
	ImageMaxSize  int64 = 2097152
	FileMaxSize   int64 = 2097152
	VideoMaxSize  int64 = 100000000
	MultiPartSize int64 = 100000000

	MimeTypeImage string = "image"
	MimeTypeVideo string = "video"

	ImageTypeJPEG string = "image/jpeg"
	ImageTypePNG  string = "image/png"
	ImageTypeWebp string = "image/webp"
	ImageTypeAll  string = "image/*"

	VideoTypeFLV         string = "video/x-flv"
	VideoTypeMP4         string = "video/mp4"
	VideoTypeMPEGURL     string = "application/x-mpegURL"
	VideoTypeMP2T        string = "video/MP2T"
	VideoType3gpp        string = "video/3gpp"
	VideoTypeQuicktime   string = "video/quicktime"
	VideoTypeMSVideo     string = "video/x-msvideo"
	VideoTypeWMV         string = "video/x-ms-wmv"
	VideoTypeOctetStream string = "application/octet-stream"
)
