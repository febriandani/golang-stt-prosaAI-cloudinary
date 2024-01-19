package user

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	cg "github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/constants/general"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	mu "github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/user"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/utils"
	su "github.com/febriandani/golang-stt-prosaAI-cloudinary/service/user"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	user su.ServiceUser
	conf general.AppService
	log  *logrus.Logger
}

func newUserHandler(user su.ServiceUser, conf general.AppService, logger *logrus.Logger) UserHandler {
	return UserHandler{
		user: user,
		conf: conf,
		log:  logger,
	}
}

func (uh UserHandler) RegistrationUser(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mu.RegistrationUser

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	message, err := uh.user.User.CreateRegistrationUser(req.Context(), param)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) UpdatePassword(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mu.ForgotPasswordRequest

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	message, err := uh.user.User.ForgotPassword(req.Context(), param)
	if err != nil {
		if err.Error() == "PS-011" {
			respData.Message = message
			respData.Status = "PS-011"
			utils.WriteResponse(res, respData, http.StatusNotAcceptable)
			return
		} else {
			respData.Message = message
			utils.WriteResponse(res, respData, http.StatusInternalServerError)
			return
		}

	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) LoginUser(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mu.LoginRequest

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {

		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	data, message, err := uh.user.User.Login(req.Context(), param)
	if err != nil {
		if err.Error() == "FailedServer" {
			respData.Message = message
			utils.WriteResponse(res, respData, http.StatusInternalServerError)
			return
		} else if err.Error() == "FailedPassword" {
			respData.Message = message
			respData.Status = "FailedPassword"
			utils.WriteResponse(res, respData, http.StatusForbidden)
			return
		} else if err.Error() == "UserNA" {
			respData.Message = message
			respData.Status = "FailedUserNA"
			utils.WriteResponse(res, respData, http.StatusNotFound)
			return
		} else {
			respData.Message = message
			utils.WriteResponse(res, respData, http.StatusInternalServerError)
			return
		}

	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
		Detail:  data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) SendOtp(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mu.OtpRequest

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	data, err := uh.user.User.SendOtp(req.Context(), param)
	if err != nil {
		if err.Error() == "waitGenerateOTP" {
			respData.Message = map[string]string{
				"en": fmt.Sprintf("Please wait %d second to request new otp code", data.Otp),
				"id": fmt.Sprintf("Mohon tunggu selama %d detik untuk request kode otp", data.Otp),
			}
			respData.Status = "waitGenerateOTP"
			utils.WriteResponse(res, respData, http.StatusTooManyRequests)
			return
		} else if err.Error() == "UserNA" {
			respData.Message = map[string]string{
				"en": "User not found",
				"id": "User tidak ditemukan",
			}
			respData.Status = "FailedUserNA"
			utils.WriteResponse(res, respData, http.StatusNotFound)
			return
		} else {
			respData.Message = map[string]string{
				"en": "Failed to request new otp code",
				"id": "Ada kesalahan saat request kode otp baru",
			}
			utils.WriteResponse(res, respData, http.StatusInternalServerError)
			return
		}

	}

	respData = &utils.ResponseDataV2{
		Status: cg.Success,
		Message: map[string]string{
			"en": "Successfully",
			"id": "Berhasil",
		},
		Detail: data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) VerifyOtp(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mu.OtpRequest

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	message, err := uh.user.User.VerifyOtp(req.Context(), param)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) GetListUsers(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataWithTotal{
		Status: cg.Fail,
	}
	var param mu.FilterUser

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	dataSession, err := utils.GetUserIDFromToken(fmt.Sprintf("%v", req.Context().Value(cg.SessionContextKey)), uh.conf.KeyData.User)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	var session cg.CredentialData

	err = json.Unmarshal([]byte(dataSession), &session)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	totalData, data, message, err := uh.user.User.GetListUsers(req.Context(), param)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataWithTotal{
		Status:    cg.Success,
		Message:   message,
		TotalData: totalData,
		Detail:    data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) ChangeStatusUser(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}

	userID := req.URL.Query().Get("user_id")

	dataSession, err := utils.GetUserIDFromToken(fmt.Sprintf("%v", req.Context().Value(cg.SessionContextKey)), uh.conf.KeyData.User)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	var session cg.CredentialData

	err = json.Unmarshal([]byte(dataSession), &session)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	status, message, err := uh.user.User.ChangeStatusUser(req.Context(), userID)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
		Detail:  status,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) GetDetailUser(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataWithTotal{
		Status: cg.Fail,
	}

	userID := req.URL.Query().Get("user_id")

	dataSession, err := utils.GetUserIDFromToken(fmt.Sprintf("%v", req.Context().Value(cg.SessionContextKey)), uh.conf.KeyData.User)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	var session cg.CredentialData

	err = json.Unmarshal([]byte(dataSession), &session)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	data, message, err := uh.user.User.GetDetailByUserID(req.Context(), userID)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataWithTotal{
		Status:  cg.Success,
		Message: message,
		Detail:  data,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) UpdateUser(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mu.UpdateUser

	dataSession, err := utils.GetUserIDFromToken(fmt.Sprintf("%v", req.Context().Value(cg.SessionContextKey)), uh.conf.KeyData.User)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	var session cg.CredentialData

	err = json.Unmarshal([]byte(dataSession), &session)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	param.CreatedBy = session.Email

	message, err := uh.user.User.UpdateDataUser(req.Context(), param)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: message,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) CheckUser(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}
	var param mu.CheckUser

	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataEmpty,
			"id": cg.HandlerErrorRequestDataEmptyID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBody, &param)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorRequestDataNotValid,
			"id": cg.HandlerErrorRequestDataNotValidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	response, message, err := uh.user.User.CheckUser(req.Context(), param)
	if err != nil {
		respData.Message = message
		utils.WriteResponse(res, respData, http.StatusInternalServerError)
		return
	}

	respData = &utils.ResponseDataV2{
		Detail:  response,
		Status:  cg.Success,
		Message: message,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}

func (uh UserHandler) UploadFile(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataV2{
		Status: cg.Fail,
	}

	dataSession, err := utils.GetUserIDFromToken(fmt.Sprintf("%v", req.Context().Value(cg.SessionContextKey)), uh.conf.KeyData.User)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	var session cg.CredentialData

	err = json.Unmarshal([]byte(dataSession), &session)
	if err != nil {
		respData.Message = map[string]string{
			"en": cg.HandlerErrorAuthInvalid,
			"id": cg.HandlerErrorAuthInvalidID,
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	// err := req.ParseMultipartForm(20 << 30) // 10 MB limit for the uploaded file
	// if err != nil {
	// 	log.Println("Error >> 10 mb ", err)
	// 	respData.Message = map[string]string{
	// 		"en": "Error uploading file",
	// 		"id": "Error uploading file",
	// 	}
	// 	utils.WriteResponse(res, respData, http.StatusBadRequest)
	// 	return
	// }

	file, handler, err := req.FormFile("upload")
	if err != nil {
		log.Println("Error form file ", err)
		respData.Message = map[string]string{
			"en": "Error uploading file",
			"id": "Error uploading file",
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Specify the directory where you want to save the uploaded file
	tempDir := "./temporaryUploads/"

	// Ensure the upload directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		os.Mkdir(tempDir, os.ModePerm)
	}

	// Create a new file in the specified directory
	dst, err := os.Create(tempDir + handler.Filename)
	if err != nil {
		log.Println("Error save dir ", err)
		respData.Message = map[string]string{
			"en": "Error uploading file",
			"id": "Error uploading file",
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}
	defer dst.Close()

	// Copy the contents of the uploaded file to the new file
	_, err = io.Copy(dst, file)
	if err != nil {
		respData.Message = map[string]string{
			"en": "Error uploading file",
			"id": "Error uploading file",
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	fmt.Printf("File %s uploaded successfully!", handler.Filename)

	cld, err := utils.CredentialsCloudinary()
	if err != nil {
		log.Println("error ", err)
		respData.Message = map[string]string{
			"en": "Error credentials file",
			"id": "Error credentials file",
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	resultCloudinary, err := utils.UploadCloudinary(cld, req.Context(), fmt.Sprintf("%s%s", tempDir, handler.Filename))
	if err != nil {
		log.Println("error ", err)
		respData.Message = map[string]string{
			"en": "Error uploading file",
			"id": "Error uploading file",
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	resultAudio := utils.SpeechToTextApi(req.Context(), fmt.Sprintf("%s%s", tempDir, handler.Filename))
	if err != nil {
		respData.Message = map[string]string{
			"en": "Error convert speech to text",
			"id": "Error convert speech to text",
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	if _, err := uh.user.User.CreateHistories(req.Context(), session.ID, resultCloudinary, resultAudio.TranscriptData); err != nil {
		respData.Message = map[string]string{
			"en": "Error to save history",
			"id": "Error to save history",
		}
		utils.WriteResponse(res, respData, http.StatusBadRequest)
		return
	}

	if err == nil {
		err := os.RemoveAll(tempDir)
		if err != nil {
			log.Println("error delete temp directory", err)
			respData.Message = map[string]string{
				"en": "Error uploading file",
				"id": "Error uploading file",
			}
			utils.WriteResponse(res, respData, http.StatusBadRequest)
			return
		}
	}

	respData = &utils.ResponseDataV2{
		Status:  cg.Success,
		Message: nil,
		Detail:  resultAudio.TranscriptData,
	}

	utils.WriteResponse(res, respData, http.StatusOK)
	return
}
