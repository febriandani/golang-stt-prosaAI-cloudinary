package general

const (
	APIHeaderContentType   string = "Content-Type"
	APIHeaderBorzoToken    string = "X-DV-Auth-Token"
	APIHeaderJetClientKey  string = "clientkey"
	APIHeaderAuthorization string = "Authorization"
)

const (
	APIHeaderContentTypeJSon           string = "application/json"
	APIHeaderContentTypeFormURLEncoded string = "application/x-www-form-urlencoded"
)

const (
	ImageFileNotFound string = "http: no such file"
)

const (
	HandlerErrorAuthInvalid                string = "authorization invalid"
	HandlerErrorAuthInvalidID              string = "authorization tidak valid"
	HandlerErrorResponseKeyIDEmpty         string = "key id cannot be empty"
	HandlerErrorRequestDataNotValid        string = "request data not valid"
	HandlerErrorRequestDataNotValidID      string = "data request tidak valid"
	HandlerErrorRequestDataEmpty           string = "request data empty"
	HandlerErrorRequestDataEmptyID         string = "data request kosong"
	HandlerErrorRequestDataFormatInvalid   string = "request data format invalid"
	HandlerErrorRequestDataFormatInvalidID string = "format data request salah"
	HandlerErrorCookiesEmpty               string = "key data cannot be empty"
	HandlerErrorCookiesInvalid             string = "key data invalid"
	HandlerErrorKeyIDInvalid               string = "key id invalid"
	HandlerErrorImageSizeTooLarge          string = "image too large, max size 1 Mb"
	HandlerErrorImageSizeTooLargeID        string = "file terlalu besar, maks. 1 Mb"
	HandlerErrorImageDataInvalid           string = "image data invalid"
	HandlerErrorImageDataInvalidID         string = "data gambar tidak sesuai"
	HandlerErrorImageDataEmpty             string = "image data cannot be empty"
	HandlerErrorImageDataEmptyID           string = "data gambar tidak boleh kosong"
	HandlerErrorFileSizeTooLarge           string = "file too large, max size 1 Mb"
	HandlerErrorFileDataInvalid            string = "file data invalid"
	HandlerErrorFileDataEmpty              string = "file data cannot be empty"
)

type CredentialData struct {
	ID        int64  `json:"id"`
	CompanyID int64  `json:"company_id"`
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
}
