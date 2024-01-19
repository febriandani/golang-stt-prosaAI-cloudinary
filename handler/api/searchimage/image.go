package searchimage

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	cg "github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/constants/general"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	mi "github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/image"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/utils"
	si "github.com/febriandani/golang-stt-prosaAI-cloudinary/service/searchimage"
	"github.com/sirupsen/logrus"
)

type ImageHandler struct {
	user si.ServiceImage
	conf general.AppService
	log  *logrus.Logger
}

func newImageHandler(user si.ServiceImage, conf general.AppService, logger *logrus.Logger) ImageHandler {
	return ImageHandler{
		user: user,
		conf: conf,
		log:  logger,
	}
}

func (ih ImageHandler) GetImages(res http.ResponseWriter, req *http.Request) {
	respData := &utils.ResponseDataWithTotal{
		Status: cg.Fail,
	}
	var param mi.RequestSearchImage

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

	totalData, data, message, err := ih.user.Image.GetImages(req.Context(), param)
	if err != nil {

		if err.Error() == "404Data" {
			respData.Message = message
			respData.Status = "404Data"
			utils.WriteResponse(res, respData, http.StatusNotFound)
			return
		} else {

			respData.Message = message
			utils.WriteResponse(res, respData, http.StatusInternalServerError)
			return
		}
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
