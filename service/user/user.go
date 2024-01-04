package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/general"
	mu "github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/model/user"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/domain/utils"
	"github.com/febriandani/golang-stt-prosaAI-cloudinary/infra"
	ru "github.com/febriandani/golang-stt-prosaAI-cloudinary/repository/user"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	db     ru.DatabaseUser
	log    *logrus.Logger
	conf   general.AppService
	dbConn *infra.DatabaseList
	redis  *redis.Client
}

func newUserService(database ru.DatabaseUser, logger *logrus.Logger, dbConn *infra.DatabaseList, conf general.AppService, redis *redis.Client) UserService {
	return UserService{
		db:     database,
		log:    logger,
		conf:   conf,
		dbConn: dbConn,
		redis:  redis,
	}
}

type User interface {
	CreateRegistrationUser(ctx context.Context, data mu.RegistrationUser) (map[string]string, error)
	ForgotPassword(ctx context.Context, data mu.ForgotPasswordRequest) (map[string]string, error)
	Login(ctx context.Context, data mu.LoginRequest) (*mu.LoginResponse, map[string]string, error)
	SendOtp(ctx context.Context, data mu.OtpRequest) (*mu.ResponseCheckUser, error)
	VerifyOtp(ctx context.Context, data mu.OtpRequest) (map[string]string, error)
	GetListUsers(ctx context.Context, filter mu.FilterUser) (int64, []mu.ResponseListDataUser, map[string]string, error)
	ChangeStatusUser(ctx context.Context, userID string) (bool, map[string]string, error)
	GetDetailByUserID(ctx context.Context, userID string) (*mu.ResponseDetailDataUser, map[string]string, error)
	UpdateDataUser(ctx context.Context, data mu.UpdateUser) (map[string]string, error)
	CreateSpeechRecognizeAudio(ctx context.Context, data mu.RequestSpeechRecognizeAudio) (map[string]string, error)
	CheckUser(ctx context.Context, data mu.CheckUser) (*mu.ResponseCheckUser, map[string]string, error)
	CreateHistories(ctx context.Context, userId int64, filename, textData string) (map[string]string, error)
}

func (us UserService) CreateRegistrationUser(ctx context.Context, data mu.RegistrationUser) (map[string]string, error) {
	messages := data.Validate()
	if messages != nil {
		return messages, errors.New("data not valid")
	}

	var (
		result   bool
		resultCh = make(chan bool)
		errorCh  = make(chan error)
	)

	getDataUser, err := us.db.User.GetUserByEmail(ctx, strings.ToLower(data.Email))
	if err != nil {
		log.Println("ERROR: ", err.Error())
		return map[string]string{
			"en": "An error occurred during registration, please try again",
			"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
		}, err
	}

	if getDataUser.Email == strings.ToLower(data.Email) {
		return map[string]string{
			"en": "Email has been registered, please try with another email",
			"id": "Email telah terdaftar, silakan coba dengan email lain",
		}, errors.New("error")
	} else if getDataUser.Email == "" {
		go func(data mu.RegistrationUser) {

			tx, err := us.dbConn.Backend.Write.Begin()
			if err != nil {
				us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Registration | fail to begin transaction")
				errorCh <- err
				return
			}

			newPassword := utils.CreatePassword(5)
			password, err := utils.GeneratePassword(fmt.Sprintf("OLIN%s", newPassword))
			if err != nil {
				tx.Rollback()
				us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Registration | fail to generate password")
				errorCh <- err
				return
			}

			NewUserIDGenerate := utils.GenerateRandIntegerSixthLength()

			_, err = us.db.User.Registration(ctx, tx, mu.RegistrationUser{
				Name:      data.Name,
				Email:     strings.ToLower(strings.ToLower(data.Email)),
				Username:  data.Username,
				CreatedAt: time.Now().UTC(),
				CreatedBy: "system",
				Password:  password,
				UserID:    NewUserIDGenerate,
			})
			if err != nil {
				tx.Rollback()
				us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Registration | fail to create registration")
				errorCh <- err
				return
			}

			err = tx.Commit()
			if err != nil {
				us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Registration | fail to commit transaction")
				tx.Rollback()
				errorCh <- err
				return
			}

			resultCh <- true
		}(data)

		select {
		case result = <-resultCh:
			if result == true {
				return map[string]string{
					"en": "User is successfully registered",
					"id": "Pengguna berhasil terdaftar",
				}, nil
			} else {
				log.Println("ERROR: ", err.Error())
				return map[string]string{
					"en": "An error occurred during registration, please try again",
					"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
				}, <-errorCh
			}

		case err := <-errorCh:
			log.Println("ERROR: ", err.Error())
			return map[string]string{
				"en": "An error occurred during registration, please try again",
				"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
			}, err
		}
	} else {
		return nil, nil
	}
}

func (us UserService) ForgotPassword(ctx context.Context, data mu.ForgotPasswordRequest) (map[string]string, error) {
	messages := data.Validate()
	if messages != nil {
		return messages, errors.New("data not valid")
	}

	isExist, err := us.db.User.IsExistUserByEmail(ctx, strings.ToLower(data.Email))
	if err != nil {
		us.log.WithField("request", utils.StructToString(nil)).WithError(err).Errorf("ForgotPassword | fail to get exist company user")
		return map[string]string{
			"en": "There was an error in checking user data",
			"id": "Ada kesalahan dalam memeriksa data pengguna",
		}, err
	}

	if !isExist {
		us.log.WithField("request", utils.StructToString(nil)).Errorf("ForgotPassword | user not found")
		return map[string]string{
			"en": "Email not registered",
			"id": "Email tidak terdaftar",
		}, errors.New("user not exist")
	}

	userData, err := us.db.User.GetUserByEmail(ctx, strings.ToLower(data.Email))
	if err != nil {
		us.log.WithField("request", utils.StructToString(nil)).WithError(err).Errorf("ForgotPassword | fail to get company user data")
		return map[string]string{
			"en": "There was an error in checking user data",
			"id": "Ada kesalahan dalam memeriksa data pengguna",
		}, err
	}

	isValid, err := utils.ComparePassword(userData.Password, data.NewPassword)
	if err != nil {
		us.log.WithField("request", utils.StructToString(nil)).WithError(err).Errorf("ForgotPassword | fail to compare password")
		return map[string]string{
			"en": "There was an error changing the password",
			"id": "Ada kesalahan dalam mengubah kata sandi",
		}, err
	}

	if isValid {
		us.log.WithField("request", utils.StructToString(nil)).Errorf("ForgotPassword | new password is same with current")
		return map[string]string{
			"en": "The current password cannot be the same as the password that was once created",
			"id": "Kata sandi saat ini tidak boleh sama dengan kata sandi yang pernah dibuat",
		}, errors.New("PS-011")
	}

	password, err := utils.GeneratePassword(data.NewPassword)
	if err != nil {
		us.log.WithField("request", utils.StructToString(nil)).WithError(err).Errorf("ForgotPassword | fail to generate password")
		return map[string]string{
			"en": "There was an error changing the password",
			"id": "Ada kesalahan dalam mengubah kata sandi",
		}, err
	}

	err = us.db.User.UpdatePassword(ctx, strings.ToLower(data.Email), password)
	if err != nil {
		us.log.WithField("request", utils.StructToString(nil)).WithError(err).Errorf("ForgotPassword | fail to change password from repo")
		return map[string]string{
			"en": "There was an error changing the password",
			"id": "Ada kesalahan dalam mengubah kata sandi",
		}, err
	}

	return map[string]string{
		"en": "success",
		"id": "password berhasil diganti",
	}, nil
}

func (us UserService) Login(ctx context.Context, data mu.LoginRequest) (*mu.LoginResponse, map[string]string, error) {
	messages := data.Validate()
	if messages != nil {
		return &mu.LoginResponse{
			Token: mu.JWTAccess{
				AccessToken:        "",
				AccessTokenExpired: "",
				RenewToken:         "",
				RenewTokenExpired:  "",
			},
			NamaLengkap: "",
		}, messages, errors.New("data not valid")
	}

	isExist, err := us.db.User.IsExistUserByEmail(ctx, strings.ToLower(strings.ToLower(data.Email)))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Login | fail to get exist company user")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
					RenewToken:         "",
					RenewTokenExpired:  "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "There was an error in checking user data",
				"id": "Ada kesalahan dalam memeriksa data pengguna",
			}, errors.New("FailedServer")
	}

	if !isExist {
		us.log.WithField("request", utils.StructToString(data)).Errorf("Login | user not found")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
					RenewToken:         "",
					RenewTokenExpired:  "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "Email not registered",
				"id": "Email tidak terdaftar",
			}, errors.New("UserNA")
	}

	userData, err := us.db.User.GetUserByEmail(ctx, strings.ToLower(strings.ToLower(data.Email)))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Login | fail to get company user data")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
					RenewToken:         "",
					RenewTokenExpired:  "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "There was an error in checking user data",
				"id": "Ada kesalahan dalam memeriksa data pengguna",
			}, errors.New("FailedServer")
	}

	isValid, err := utils.ComparePassword(userData.Password, data.Password)
	if err != nil {
		us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Login | fail to compare password")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
					RenewToken:         "",
					RenewTokenExpired:  "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "There was an error changing the password",
				"id": "Ada kesalahan dalam mengubah kata sandi",
			}, errors.New("FailedServer")
	}

	if !isValid {
		us.log.WithField("request", utils.StructToString(data)).Errorf("Login | incorrect password")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
					RenewToken:         "",
					RenewTokenExpired:  "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "You entered an incorrect password",
				"id": "Kata sandi yang anda masukkan salah",
			}, errors.New("FailedPassword")
	}

	session, err := utils.GetEncrypt([]byte(us.conf.KeyData.User), utils.StructToString(mu.CredentialData{
		ID:       userData.ID,
		Fullname: userData.Name,
		Email:    userData.Email,
	}))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data)).Errorf("Login | fail generate session jwt")
		log.Println(err.Error())
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
					RenewToken:         "",
					RenewTokenExpired:  "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "internal server error",
				"id": "terjadi kesalahan sistem, silahkan coba dilain waktu ",
			}, errors.New("FailedServer")
	}

	generateTime := time.Now().UTC()

	accessToken, renewToken, err := utils.GenerateJWT(session)
	if err != nil {
		us.log.WithField("request", utils.StructToString(data)).Errorf("Login | fail generate jwt token")
		return &mu.LoginResponse{
				Token: mu.JWTAccess{
					AccessToken:        "",
					AccessTokenExpired: "",
					RenewToken:         "",
					RenewTokenExpired:  "",
				},
				NamaLengkap: "",
			}, map[string]string{
				"en": "internal server error",
				"id": "terjadi kesalahan sistem, silahkan coba dilain waktu ",
			}, errors.New("FailedServer")
	}

	resLogin := mu.LoginResponse{
		Token: mu.JWTAccess{
			AccessToken:        accessToken,
			AccessTokenExpired: generateTime.Add(time.Duration(us.conf.Authorization.JWT.AccessTokenDuration) * time.Minute).Format(time.RFC3339),
			RenewToken:         renewToken,
			RenewTokenExpired:  generateTime.Add(time.Duration(us.conf.Authorization.JWT.RefreshTokenDuration) * time.Minute).Format(time.RFC3339),
		},
		NamaLengkap: userData.Name,
	}

	return &resLogin, nil, nil
}

func (us UserService) SendOtp(ctx context.Context, data mu.OtpRequest) (*mu.ResponseCheckUser, error) {

	otpCode := utils.GenerateRandIntegerFourthLengthString()

	isExist, err := us.db.User.IsExistUserByEmail(ctx, strings.ToLower(data.Email))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Failed to get data user")
		return nil, err
	}

	if !isExist {
		us.log.WithField("request", utils.StructToString(data)).Errorf("Error user not found")
		return nil, errors.New("UserNA")
	}

	value, err := us.db.User.GetOTP(ctx, fmt.Sprintf("OTP-%s", strings.ToLower(data.Email)))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Failed to get otp")
		return nil, err
	}

	log.Println("OTP VALUE :", value)

	if value.Email == "" {

		//save otp
		_, err := us.db.User.SaveOTP(ctx, nil, fmt.Sprintf("OTP-%s", strings.ToLower(data.Email)), otpCode)
		if err != nil {
			us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Failed to save otp")
			return nil, err
		}

		return &mu.ResponseCheckUser{
			Otp:      otpCode,
			IsActive: true,
			Email:    strings.ToLower(data.Email),
		}, nil

	} else if value.Email != "" && value.ExpiredAt < time.Now().Format("2006-01-02 15:04:05") {
		//tode update otp

		err := us.db.User.UpdateOTP(ctx, nil, fmt.Sprintf("OTP-%s", strings.ToLower(data.Email)), otpCode)
		if err != nil {
			us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Failed to update otp")
			return nil, err
		}

		return &mu.ResponseCheckUser{
			Otp:      otpCode,
			IsActive: true,
			Email:    strings.ToLower(data.Email),
		}, nil
	} else {
		layout := "2006-01-02 15:04:05"

		// Parse the string representation of time into a time.Time object in the "Asia/Jakarta" time zone
		location, err := time.LoadLocation("Asia/Jakarta")
		if err != nil {
			fmt.Println("Error loading time zone:", err)
			return nil, err
		}

		valueExpiredAt, err := time.ParseInLocation(layout, value.ExpiredAt, location)
		if err != nil {
			fmt.Println("Error parsing time:", err)
			return nil, err
		}

		// Convert both times to UTC
		valueExpiredAt = valueExpiredAt.UTC()
		now := time.Now().UTC()

		// Calculate the duration between value.ExpiredAt and the current time
		duration := valueExpiredAt.Sub(now)

		// Get the duration in seconds
		seconds := int(duration.Seconds())

		return &mu.ResponseCheckUser{
			Otp:      seconds,
			IsActive: true,
			Email:    strings.ToLower(data.Email),
		}, errors.New("waitGenerateOTP")
	}

}

func (us UserService) VerifyOtp(ctx context.Context, data mu.OtpRequest) (map[string]string, error) {
	value, err := us.redis.Get(context.Background(), fmt.Sprintf("BACKSTAGE-%s", strings.ToLower(data.Email))).Result()
	if err != nil {
		if err != redis.Nil {
			log.Println("ERROR GET REDIS", err.Error())
			return getErrorMessageVerifyOTP(), err
		}

		return map[string]string{
			"en": "Otp code has expired, please request otp code again.",
			"id": "Kode otp telah kedaluwarsa, silakan minta kode otp lagi.",
		}, errors.New("ERROR")
	}

	var user mu.ResponseOtpFromRedis
	err = json.Unmarshal([]byte(value), &user)
	if err != nil {
		log.Println("ERROR UNMARSHAL", err.Error())
		return getErrorMessageVerifyOTP(), err
	}

	if user.Otp == data.OTP.String {

		err := us.redis.Del(ctx, fmt.Sprintf("BACKSTAGE-%s", strings.ToLower(data.Email))).Err()
		if err != nil {
			log.Println("ERROR DEL REDIS", err.Error())
			return getErrorMessageVerifyOTP(), err
		}

		return map[string]string{
			"en": "Otp code successfully verified",
			"id": "Kode otp berhasil diverifikasi",
		}, nil
	}

	return map[string]string{
		"en": "Invalid otp code",
		"id": "Kode otp tidak valid",
	}, errors.New("ERROR")
}

func (us UserService) GetListUsers(ctx context.Context, filter mu.FilterUser) (int64, []mu.ResponseListDataUser, map[string]string, error) {

	totalData, err := us.db.User.GetTotalListUser(ctx, filter)
	if err != nil {
		log.Println("ERROR GET TOTAL DATA USERS", err.Error())
		return totalData, nil, map[string]string{
			"en": "There was an error during get data users",
			"id": "Ada kesalahan saat menampilkan data staff",
		}, err
	}

	data, err := us.db.User.GetListUser(ctx, filter)
	if err != nil {
		log.Println("ERROR GET DATA USERS", err.Error())
		return totalData, nil, map[string]string{
			"en": "There was an error during get data users",
			"id": "Ada kesalahan saat menampilkan data staff",
		}, err
	}

	return totalData, data, map[string]string{
		"en": "Successfully retrieved data users",
		"id": "Berhasil menampilkan data staff",
	}, nil
}

func (us UserService) ChangeStatusUser(ctx context.Context, userID string) (bool, map[string]string, error) {

	var status bool

	dataUser, err := us.db.User.GetUserByUserID(ctx, userID)
	if err != nil {
		return status, map[string]string{
			"en": "There was an error during change status user",
			"id": "Ada kesalahan saat mengganti status user",
		}, err
	}

	if dataUser.IsActive == true {
		status = false
	} else {
		status = true
	}

	err = us.db.User.ChangeStatusUser(ctx, status, userID)
	if err != nil {
		return status, map[string]string{
			"en": "There was an error during change status user",
			"id": "Ada kesalahan saat mengganti status user",
		}, err
	}

	return status, map[string]string{
		"en": "Successfully changed status",
		"id": "Berhasil mengganti status",
	}, nil
}

func (us UserService) GetDetailByUserID(ctx context.Context, userID string) (*mu.ResponseDetailDataUser, map[string]string, error) {
	detail, err := us.db.User.GetDetailByUserID(ctx, userID)
	if err != nil {
		log.Println("ERROR GET DATA USERS", err.Error())
		return nil, map[string]string{
			"en": "There was an error during get data",
			"id": "Ada kesalahan saat menampilkan data",
		}, err
	}

	return detail, map[string]string{
		"en": "Successfully",
		"id": "Berhasil",
	}, nil
}

func (us UserService) UpdateDataUser(ctx context.Context, data mu.UpdateUser) (map[string]string, error) {
	messages := data.Validate()
	if messages != nil {
		return messages, errors.New("data not valid")
	}

	var (
		result   bool
		resultCh = make(chan bool)
		errorCh  = make(chan error)
	)

	go func(data mu.UpdateUser) {

		tx, err := us.dbConn.Backend.Write.Begin()
		if err != nil {
			us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Registration | fail to begin transaction")
			errorCh <- err
			return
		}

		err = us.db.User.UpdateDataUser(ctx, tx, mu.UpdateUser{
			Name:                         data.Name,
			Position:                     data.Position,
			Email:                        strings.ToLower(data.Email),
			Partner:                      data.Partner,
			IsUserChild:                  data.IsUserChild,
			IsMpiChild:                   data.IsMpiChild,
			IsPartnerChild:               data.IsPartnerChild,
			IsProductChild:               data.IsProductChild,
			IsDiscountChild:              data.IsDiscountChild,
			IsSettingProductChild:        data.IsSettingProductChild,
			IsTransactionChild:           data.IsTransactionChild,
			IsBillingChild:               data.IsBillingChild,
			IsHomepageChild:              data.IsHomepageChild,
			IsIminChild:                  data.IsIminChild,
			IsPurchaseChild:              data.IsPurchaseChild,
			IsNewsChild:                  data.IsNewsChild,
			IsTermsConditionChild:        data.IsTermsConditionChild,
			IsSubscriptionChild:          data.IsSubscriptionChild,
			IsRefferalFeeChild:           data.IsRefferalFeeChild,
			IsPlatformFeeChild:           data.IsPlatformFeeChild,
			IsMarketingFeeChild:          data.IsMarketingFeeChild,
			IsTransactionQrisChild:       data.IsTransactionQrisChild,
			IsTaxInvoiceChild:            data.IsTaxInvoiceChild,
			IsUserManagementChild:        data.IsUserManagementChild,
			IsReportingTransactionChild:  data.IsReportingTransactionChild,
			IsReportingFinanceChild:      data.IsReportingFinanceChild,
			IsSubscriptionMarketingChild: data.IsSubscriptionMarketingChild,
			CreatedAt:                    time.Now().UTC(),
			CreatedBy:                    data.CreatedBy,
			Phone:                        data.Phone,
			TypePartner:                  data.TypePartner,
			UserID:                       data.UserID,
		})
		if err != nil {
			tx.Rollback()
			us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("User Management | fail to Update User Management")
			errorCh <- err
			return
		}

		err = tx.Commit()
		if err != nil {
			us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("User Management | fail to commit transaction")
			tx.Rollback()
			errorCh <- err
			return
		}

		resultCh <- true

	}(data)

	select {
	case result = <-resultCh:
		if result == true {

			return map[string]string{
				"en": "User is successfully registered",
				"id": "Pengguna berhasil terdaftar",
			}, nil
		} else {
			return map[string]string{
				"en": "An error occurred during update data user, please try again",
				"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
			}, <-errorCh
		}

	case err := <-errorCh:
		log.Println("ERROR: ", err.Error())
		return map[string]string{
			"en": "An error occurred during update data user, please try again",
			"id": "Terjadi kesalahan saat registrasi, silakan coba lagi",
		}, err
	}
}

func (us UserService) CreateSpeechRecognizeAudio(ctx context.Context, data mu.RequestSpeechRecognizeAudio) (map[string]string, error) {

	return nil, nil
}

func getErrorMessageVerifyOTP() map[string]string {
	return map[string]string{
		"en": "There was an error during otp verification",
		"id": "Ada kesalahan saat verifikasi otp",
	}
}

func getErrorMessage() map[string]string {
	return map[string]string{
		"en": "There was an error during otp request",
		"id": "Ada kesalahan saat permintaan otp",
	}
}

func (us UserService) CheckUser(ctx context.Context, data mu.CheckUser) (*mu.ResponseCheckUser, map[string]string, error) {

	otpCode := utils.GenerateRandIntegerFourthLengthString()

	//checkIsExistUser
	isExist, err := us.db.User.IsExistUserByEmail(ctx, strings.ToLower(data.Email))
	if err != nil {
		us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Failed to check data user ")
		return nil, nil, err
	}

	if !isExist {
		us.log.WithField("request", utils.StructToString(data)).WithError(err).Errorf("Uset not found ")
		return nil, nil, errors.New("UserNA")
	}

	return &mu.ResponseCheckUser{
		Otp:      otpCode,
		IsActive: true,
		Email:    strings.ToLower(data.Email),
	}, nil, nil
}

func (us UserService) CreateHistories(ctx context.Context, userId int64, filename, textData string) (map[string]string, error) {

	_, err := us.db.User.SaveAudio(ctx, nil, userId, filename, textData)
	if err != nil {
		us.log.WithField("request ", utils.StructToString(userId)).WithError(err).Errorf("Failed to save audio")
		return nil, err
	}

	return nil, nil
}
