package general

type SectionService struct {
	// AppAccount
	AppName         string `json:"APP_NAME"`
	AppEnvirontment string `json:"APP_ENV"`
	AppURL          string `json:"APP_URL"`
	AppPort         string `json:"APP_PORT"`
	AppSecretKey    string `json:"APP_KEY"`

	// RouteAccount
	RouteMethods string `json:"ROUTES_METHODS"`
	RouteHeaders string `json:"ROUTES_HEADERS"`
	RouteOrigins string `json:"ROUTES_ORIGIN"`

	// DatabaseAccount
	// Read
	DatabaseReadUsername     string `json:"DATABASE_READ_USERNAME"`
	DatabaseReadPassword     string `json:"DATABASE_READ_PASSWORD"`
	DatabaseReadURL          string `json:"DATABASE_READ_URL"`
	DatabaseReadPort         string `json:"DATABASE_READ_PORT"`
	DatabaseReadDBName       string `json:"DATABASE_READ_NAME"`
	DatabaseReadMaxIdleConns string `json:"DATABASE_READ_MAXIDLECONNS"`
	DatabaseReadMaxOpenConns string `json:"DATABASE_READ_MAXOPENCONNS"`
	DatabaseReadMaxLifeTime  string `json:"DATABASE_READ_MAXLIFETIME"`
	DatabaseReadTimeout      string `json:"DATABASE_READ_TIMEOUT"`
	DatabaseReadSSLMode      string `json:"DATABASE_READ_SSL_MODE"`

	// Write
	DatabaseWriteUsername     string `json:"DATABASE_WRITE_USERNAME"`
	DatabaseWritePassword     string `json:"DATABASE_WRITE_PASSWORD"`
	DatabaseWriteURL          string `json:"DATABASE_WRITE_URL"`
	DatabaseWritePort         string `json:"DATABASE_WRITE_PORT"`
	DatabaseWriteDBName       string `json:"DATABASE_WRITE_NAME"`
	DatabaseWriteMaxIdleConns string `json:"DATABASE_WRITE_MAXIDLECONNS"`
	DatabaseWriteMaxOpenConns string `json:"DATABASE_WRITE_MAXOPENCONNS"`
	DatabaseWriteMaxLifeTime  string `json:"DATABASE_WRITE_MAXLIFETIME"`
	DatabaseWriteTimeout      string `json:"DATABASE_WRITE_TIMEOUT"`
	DatabaseWriteSSLMode      string `json:"DATABASE_WRITE_SSL_MODE"`

	// Redis
	RedisUsername     string `json:"REDIS_USERNAME"`
	RedisPassword     string `json:"REDIS_PASSWORD"`
	RedisURL          string `json:"REDIS_URL"`
	RedisPort         string `json:"REDIS_PORT"`
	RedisMinIdleConns string `json:"REDIS_MINIDLECONNS"`
	RedisTimeout      string `json:"REDIS_TIMEOUT"`

	// Authorization
	// JWT
	AuthorizationJWTIsActive              string `json:"AUTHORIZATION_JWT_IS_ACTIVE"`
	AuthorizationJWTAccessTokenSecretKey  string `json:"AUTHORIZATION_JWT_ACCESS_TOKEN_SECRET_KEY"`
	AuthorizationJWTAccessTokenDuration   string `json:"AUTHORIZATION_JWT_ACCESS_TOKEN_DURATION"`
	AuthorizationJWTRefreshTokenSecretKey string `json:"AUTHORIZATION_JWT_REFRESH_TOKEN_SECRET_KEY"`
	AuthorizationJWTRefreshTokenDuration  string `json:"AUTHORIZATION_JWT_REFRESH_TOKEN_DURATION"`

	// Public
	AuthorizationPublicSecretKey string `json:"AUTHORIZATION_PUBLIC_SECRET_KEY"`

	// Key Account
	KeyAccountUser string `json:"KEY_USER"`

	// Minio
	MinioBucketName  string `json:"MINIO_BUCKET_NAME"`
	MinioEndpoint    string `json:"MINIO_ENDPOINT"`
	MinioKey         string `json:"MINIO_ACCESS_KEY_ID"`
	MinioSecret      string `json:"MINIO_SECRET_ACCESS_KEY"`
	MinioRegion      string `json:"MINIO_REGION"`
	MinioTempFolder  string `json:"MINIO_TEMPFOLDER"`
	MinioBaseURL     string `json:"MINIO_BASE_URL"`
	MinioURLDuration string `json:"MINIO_URL_DURATION"`
}

type AppService struct {
	App           AppAccount   `json:",omitempty"`
	Route         RouteAccount `json:",omitempty"`
	DatabaseUser  DatabaseUser `json:",omitempty"`
	Redis         RedisAccount `json:",omitempty"`
	Authorization AuthAccount  `json:",omitempty"`
	KeyData       KeyAccount   `json:",omitempty"`
	Minio         MinioSecret  `json:",omitempty"`
}

type AppAccount struct {
	Name         string `json:",omitempty"`
	Environtment string `json:",omitempty"`
	URL          string `json:",omitempty"`
	Port         string `json:",omitempty"`
	SecretKey    string `json:",omitempty"`
}

type RouteAccount struct {
	Methods []string `json:",omitempty"`
	Headers []string `json:",omitempty"`
	Origins []string `json:",omitempty"`
}

type DatabaseUser struct {
	Read  DBDetailUser `json:",omitempty"`
	Write DBDetailUser `json:",omitempty"`
}

type DBDetailUser struct {
	Username     string `json:",omitempty"`
	Password     string `json:",omitempty"`
	URL          string `json:",omitempty"`
	Port         string `json:",omitempty"`
	DBName       string `json:",omitempty"`
	MaxIdleConns int    `json:",omitempty"`
	MaxOpenConns int    `json:",omitempty"`
	MaxLifeTime  int    `json:",omitempty"`
	Timeout      string `json:",omitempty"`
	SSLMode      string `json:",omitempty"`
}

type RedisAccount struct {
	Username     string `json:",omitempty"`
	Password     string `json:",omitempty"`
	URL          string `json:",omitempty"`
	Port         int    `json:",omitempty"`
	MinIdleConns int    `json:",omitempty"`
	Timeout      string `json:",omitempty"`
}

type AuthAccount struct {
	JWT    JWTCredential    `json:",omitempty"`
	Public PublicCredential `json:",omitempty"`
}

type JWTCredential struct {
	IsActive              bool   `json:",omitempty"`
	AccessTokenSecretKey  string `json:",omitempty"`
	AccessTokenDuration   int    `json:",omitempty"`
	RefreshTokenSecretKey string `json:",omitempty"`
	RefreshTokenDuration  int    `json:",omitempty"`
}

type PublicCredential struct {
	SecretKey string `json:",omitempty"`
}

type KeyAccount struct {
	User string `json:",omitempty"`
}

type MinioSecret struct {
	BucketName string `json:",omitempty"`
	Endpoint   string `json:",omitempty"`
	Key        string `json:",omitempty"`
	Secret     string `json:",omitempty"`
	Region     string `json:",omitempty"`
	TempFolder string `json:",omitempty"`
	BaseURL    string `json:",omitempty"`
}
