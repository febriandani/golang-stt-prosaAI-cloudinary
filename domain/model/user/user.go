package user

import (
	"time"

	"net/mail"

	"github.com/lib/pq"
	"gopkg.in/guregu/null.v4"
)

type RegistrationUser struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt time.Time
	CreatedBy string
	UpdatedBy string
	UserID    int
}

type ForgotPasswordRequest struct {
	Email         string `json:"email"`
	NewPassword   string `json:"new_password"`
	ReNewPassword string `json:"repeat_new_password"`
}

type ResponseDataUser struct {
	ID       int64  `json:"-" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"-" db:"password"`
	IsActive bool   `json:"-" db:"is_active"`
}

type ResponseDetailDataUser struct {
	ID                          int64         `json:"user_id" db:"user_id"`
	Name                        string        `json:"name" db:"employee_name"`
	Position                    string        `json:"position" db:"position"`
	Email                       string        `json:"email" db:"email"`
	Partner                     string        `json:"partner_id" db:"partner_id"`
	PartnerType                 int64         `json:"partner_type" db:"type_partner"`
	PartnerName                 string        `json:"partner_name" db:"partner_name"`
	Phone                       string        `json:"phone" db:"phone"`
	IsUserChild                 pq.Int64Array `json:"isUserChild" db:"user_role"`
	IsMpiChild                  pq.Int64Array `json:"isMpiChild" db:"mpi"`
	IsPartnerChild              pq.Int64Array `json:"isPartnerChild" db:"partner"`
	IsProductChild              pq.Int64Array `json:"isProductChild" db:"product"`
	IsDiscountChild             pq.Int64Array `json:"isDiscountChild" db:"discount"`
	IsSettingProductChild       pq.Int64Array `json:"isSettingProductChild" db:"setting_product"`
	IsTransactionChild          pq.Int64Array `json:"isTransactionChild" db:"transaction"`
	IsBillingChild              pq.Int64Array `json:"isBillingChild" db:"billing"`
	IsHomepageChild             pq.Int64Array `json:"isHomepageChild" db:"homepage"`
	IsIminChild                 pq.Int64Array `json:"isIminChild" db:"imin"`
	IsPurchaseChild             pq.Int64Array `json:"isPurchaseChild" db:"purchase"`
	IsNewsChild                 pq.Int64Array `json:"isNewsChild" db:"news"`
	IsTermsConditionChild       pq.Int64Array `json:"isTermsConditionChild" db:"terms_condition"`
	IsSubscriptionChild         pq.Int64Array `json:"isSubscriptionChild" db:"subscription"`
	IsRefferalFeeChild          pq.Int64Array `json:"isRefferalFeeChild" db:"refferal_fee"`
	IsPlatformFeeChild          pq.Int64Array `json:"isPlatformFeeChild" db:"platform_fee"`
	IsMarketingFeeChild         pq.Int64Array `json:"isMarketingFeeChild" db:"marketing"`
	IsTransactionQrisChild      pq.Int64Array `json:"isTransactionQrisChild" db:"transaction_qris"`
	IsTaxInvoiceChild           pq.Int64Array `json:"isTaxInvoiceChild" db:"tax_invoice"`
	IsUserManagementChild       pq.Int64Array `json:"isUserManagementChild" db:"user_management"`
	IsReportingTransactionChild pq.Int64Array `json:"isReportingTransactionChild" db:"reporting_transaction"`
	IsReportingFinanceChild     pq.Int64Array `json:"isReportingFinanceChild" db:"reporting_finance"`
	IsSubcriptionMarketingChild pq.Int64Array `json:"isSubcriptionMarketingChild" db:"subscription_marketing"`
}

type OtpRequest struct {
	Email string      `json:"email" db:"email"`
	OTP   null.String `json:"otp" db:"otp"`
}

type ResponseOtpFromRedis struct {
	Email       string `json:"email"`
	Otp         string `json:"otp"`
	ExpiredDate string `json:"expired_date"`
}

type ResponseListDataUser struct {
	UserID       string `json:"user_id" db:"user_id"`
	EmployeeName string `json:"user_name" db:"employee_name"`
	Email        string `json:"user_email" db:"user_email"`
	Position     string `json:"position" db:"position"`
	CreatedAt    string `json:"date_created" db:"date_created"`
	Status       string `json:"status" db:"status_name"`
	StatusID     int    `json:"status_id" db:"is_active"`
}

type FilterUser struct {
	Username   null.String `json:"username"`
	Email      null.String `json:"email"`
	Status     null.Int    `json:"status"`
	Offset     null.Int    `json:"offset"`
	Limit      null.Int    `json:"limit"`
	Statusinit string
}

func (r *ForgotPasswordRequest) Validate() map[string]string {
	if r.Email == "" {
		return map[string]string{
			"en": "Email cannot be empty",
			"id": "Email tidak boleh kosong",
		}
	}

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		return map[string]string{
			"en": "Incorrect email format",
			"id": "Format email salah",
		}
	}

	if r.NewPassword == "" {
		return map[string]string{
			"en": "New Password cannot be empty",
			"id": "ulangi kata sandi baru tidak boleh kosong",
		}
	}

	if r.ReNewPassword == "" {
		return map[string]string{
			"en": "Repeat New Password cannot be empty",
			"id": "Ulangi kata sandi baru tidak boleh kosong",
		}
	}
	return nil
}

func (r *RegistrationUser) Validate() map[string]string {
	if r.Email == "" {
		return map[string]string{
			"en": "Email cannot be empty",
			"id": "Email tidak boleh kosong",
		}
	}

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		return map[string]string{
			"en": "Incorrect email format",
			"id": "Format email salah",
		}
	}

	if r.Name == "" {
		return map[string]string{
			"en": "Name cannot be empty",
			"id": "Nama tidak boleh kosong",
		}
	}

	if r.Password == "" {
		return map[string]string{
			"en": "Password cannot be empty",
			"id": "Password tidak boleh kosong",
		}
	}

	return nil
}

func (r *UpdateUser) Validate() map[string]string {
	if r.Email == "" {
		return map[string]string{
			"en": "Email cannot be empty",
			"id": "Email tidak boleh kosong",
		}
	}

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		return map[string]string{
			"en": "Incorrect email format",
			"id": "Format email salah",
		}
	}

	if r.Name == "" {
		return map[string]string{
			"en": "Name cannot be empty",
			"id": "Nama tidak boleh kosong",
		}
	}

	if r.Position == "" {
		return map[string]string{
			"en": "Position cannot be empty",
			"id": "Posisi tidak boleh kosong",
		}
	}

	return nil
}

func (r *OtpRequest) Validate() map[string]string {
	if r.Email == "" {
		return map[string]string{
			"en": "Email cannot be empty",
			"id": "Email tidak boleh kosong",
		}
	}

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		return map[string]string{
			"en": "Incorrect email format",
			"id": "Format email salah",
		}
	}
	return nil
}

type UpdateUser struct {
	Name                         string        `json:"name"`
	Position                     string        `json:"position"`
	Email                        string        `json:"email"`
	Partner                      string        `json:"partner"`
	IsUserChild                  pq.Int64Array `json:"isUserChild"`
	IsMpiChild                   pq.Int64Array `json:"isMpiChild"`
	IsPartnerChild               pq.Int64Array `json:"isPartnerChild"`
	IsProductChild               pq.Int64Array `json:"isProductChild"`
	IsDiscountChild              pq.Int64Array `json:"isDiscountChild"`
	IsSettingProductChild        pq.Int64Array `json:"isSettingProductChild"`
	IsTransactionChild           pq.Int64Array `json:"isTransactionChild"`
	IsBillingChild               pq.Int64Array `json:"isBillingChild"`
	IsHomepageChild              pq.Int64Array `json:"isHomepageChild"`
	IsIminChild                  pq.Int64Array `json:"isIminChild"`
	IsPurchaseChild              pq.Int64Array `json:"isPurchaseChild"`
	IsNewsChild                  pq.Int64Array `json:"isNewsChild"`
	IsTermsConditionChild        pq.Int64Array `json:"isTermsConditionChild"`
	IsSubscriptionChild          pq.Int64Array `json:"isSubscriptionChild"`
	IsRefferalFeeChild           pq.Int64Array `json:"isRefferalFeeChild"`
	IsPlatformFeeChild           pq.Int64Array `json:"isPlatformFeeChild"`
	IsMarketingFeeChild          pq.Int64Array `json:"isMarketingFeeChild"`
	IsTransactionQrisChild       pq.Int64Array `json:"isTransactionQrisChild"`
	IsTaxInvoiceChild            pq.Int64Array `json:"isTaxInvoiceChild"`
	IsUserManagementChild        pq.Int64Array `json:"isUserManagementChild"`
	IsReportingTransactionChild  pq.Int64Array `json:"isReportingTransactionChild"`
	IsReportingFinanceChild      pq.Int64Array `json:"isReportingFinanceChild"`
	IsSubscriptionMarketingChild pq.Int64Array `json:"isSubscriptionMarketingChild"`
	CreatedAt                    time.Time
	CreatedBy                    string
	UpdatedBy                    string
	Password                     string
	UserID                       int64  `json:"user_id"`
	Phone                        string `json:"phone"`
	TypePartner                  int64  `json:"type_partner"`
}

type RequestSpeechRecognizeAudio struct {
	Username string `json:"username"`
	Audio1   string `json:"audio1"`
	Audio2   string `json:"audio2"`
	Audio3   string `json:"audio3"`
}

type CheckUser struct {
	Email string `json:"email"`
}

type ForgotPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
}

type ResponseCheckUser struct {
	Otp      int    `json:"otp"`
	IsActive bool   `json:"isActive"`
	Email    string `json:"email"`
}

type OtpResponse struct {
	OtpCode   int64  `json:"otp" db:"otp_code"`
	ExpiredAt string `json:"expiredAt" db:"expired_at"`
	Email     string `json:"email" db:"email"`
}
