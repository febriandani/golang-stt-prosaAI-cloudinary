package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	constants "github.com/pharmaniaga/auth-user/domain/constants/general"
	"github.com/pharmaniaga/auth-user/domain/model/general"
)

type ResponseHTTP struct {
	StatusCode int
	Response   ResponseData
}

type ResponseData struct {
	Status  string      `json:"status"`
	Source  string      `json:"source,omitempty"`
	Message string      `json:"message,omitempty"`
	Detail  interface{} `json:"detail"`
}

type ResponseDataV2 struct {
	Status  string            `json:"status"`
	Message map[string]string `json:"message,omitempty"`
	Detail  interface{}       `json:"detail,omitempty"`
}

type ResponseDataWithTotal struct {
	Status      string            `json:"status"`
	Message     map[string]string `json:"message,omitempty"`
	TotalData   int64             `json:"total_data,omitempty"`
	TotalAmount int64             `json:"total_amount,omitempty"`
	Detail      interface{}       `json:"detail,omitempty"`
}

type ResponseInvoice struct {
	Status  string                `json:"status,omitempty"`
	Message map[string]string     `json:"message,omitempty"`
	Detail  DetailResponseInvoice `json:"detail,omitempty"`
}

type DetailResponseInvoice struct {
	ListMedicine  interface{} `json:"list_medicine,omitempty"`
	Subtotal      interface{} `json:"subtotal,omitempty"`
	DetailPayment interface{} `json:"detail_payment,omitempty"`
	DetailInvoice interface{} `json:"detail_invoice,omitempty"`
}

// Response is the new type for define all of the response from service
type Response interface{}

var (
	ErrRespServiceMaintance = ResponseHTTP{
		StatusCode: http.StatusServiceUnavailable,
		Response:   ResponseData{Status: constants.Fail}}
	ErrRespUnauthorize = ResponseHTTP{
		StatusCode: http.StatusUnauthorized,
		Response:   ResponseData{Status: constants.Fail}}
	ErrRespAuthInvalid = ResponseHTTP{
		StatusCode: http.StatusUnauthorized,
		Response:   ResponseData{Status: constants.Fail}}
	ErrRespBadRequest = ResponseHTTP{
		StatusCode: http.StatusBadRequest,
		Response:   ResponseData{Status: constants.Fail}}
	ErrRespInternalServer = ResponseHTTP{
		StatusCode: http.StatusServiceUnavailable,
		Response:   ResponseData{Status: constants.Fail}}
)

func WriteResponse(res http.ResponseWriter, resp Response, code int) {
	res.Header().Set("Content-Type", "application/json")
	r, _ := json.Marshal(resp)

	res.WriteHeader(code)
	res.Write(r)
	return
}

func WriteResponseWs(ws *websocket.Conn, res http.ResponseWriter, resp interface{}, code int) {
	err := ws.WriteJSON(resp)
	if err != nil {
		log.Println("ERROR write ws response", err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	ws.Close()

	return
}

type Error struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	Title  string `json:"title"`
}

func NewError(id string, status string, title string) *Error {
	return &Error{
		Id:     id,
		Status: status,
		Title:  title,
	}
}

func (rd *ResponseData) GenerateErrorResponse(data *general.ResponseData, errorMsg string) {
	data.Error = errorMsg
	rd.Detail = data
}

func GenerateUUID() string {
	/* generating random no invoice */
	rand.Seed(time.Now().Unix())
	currentTime := time.Now().UTC()
	rangeLower := 12345
	rangeUpper := 98765
	randomNum := rangeLower + rand.Intn(rangeUpper-rangeLower+1)

	length := 6 // change the length as needed

	date := currentTime.Format("02")
	year := currentTime.Format("2006")
	randomString := fmt.Sprintf("%v%d%v", date, int64(randomNum), year)
	for len(randomString) < length {
		randomString += strconv.Itoa(rand.Intn(10))
	}
	return randomString[:length]

}

func ConvertEpochToDateTime(epoch int64) string {
	// Convert epoch timestamp to time.Time in UTC
	t := time.Unix(epoch/1000, 0).UTC()

	// Convert UTC time to the local time zone
	localTime := t.Local()

	// Format the local time in the desired layout
	formattedTime := localTime.Format("2006-01-02 15:04:05")

	return formattedTime
}
