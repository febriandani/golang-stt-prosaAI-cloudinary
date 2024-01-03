package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pharmaniaga/auth-user/domain/model/general"
	"github.com/pharmaniaga/auth-user/handler/api"
)

func getV1(freeRoute, router, routerJWT, wsRoute *mux.Router, conf *general.AppService, handler api.Handler) {

	// router.HandleFunc("/v1/upload/image", handler.Media.Spaces.UploadImage).Methods(http.MethodPost)
	// router.HandleFunc("/v1/upload/file", handler.Media.Spaces.UploadFile).Methods(http.MethodPost)
	// router.HandleFunc("/v1/file/remove", handler.Media.Spaces.RemoveFile).Methods(http.MethodDelete)
	freeRoute.HandleFunc("/v1/registration-user", handler.User.User.RegistrationUser).Methods(http.MethodPost)
	freeRoute.HandleFunc("/v1/login-user", handler.User.User.LoginUser).Methods(http.MethodPost)
	freeRoute.HandleFunc("/v1/checkuserv2", handler.User.User.CheckUser).Methods(http.MethodPost)
	freeRoute.HandleFunc("/v1/change-password-user", handler.User.User.UpdatePassword).Methods(http.MethodPost)
	freeRoute.HandleFunc("/v1/request-otp", handler.User.User.SendOtp).Methods(http.MethodPost)
	router.HandleFunc("/v1/verify-otp", handler.User.User.VerifyOtp).Methods(http.MethodPost)
	routerJWT.HandleFunc("/v1/users", handler.User.User.GetListUsers).Methods(http.MethodPost)
	routerJWT.HandleFunc("/v1/change-status-user", handler.User.User.ChangeStatusUser).Methods(http.MethodPost)
	routerJWT.HandleFunc("/v1/user-detail", handler.User.User.GetDetailUser).Methods(http.MethodGet)
	routerJWT.HandleFunc("/v1/update-user", handler.User.User.UpdateUser).Methods(http.MethodPost)
	// freeRoute.Handle("/", http.FileServer(http.Dir("static"))) // Serve a simple HTML/JS client for demonstration
	routerJWT.HandleFunc("/v1/upload", handler.User.User.UploadFile).Methods(http.MethodPost)
}
