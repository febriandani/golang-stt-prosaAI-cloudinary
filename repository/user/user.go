package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pharmaniaga/auth-user/domain/model/user"
	"github.com/pharmaniaga/auth-user/infra"
	"github.com/sirupsen/logrus"
)

type UserConfig struct {
	db  *infra.DatabaseList
	log *logrus.Logger
}

func newDatabaseUser(db *infra.DatabaseList, logger *logrus.Logger) UserConfig {
	return UserConfig{
		db:  db,
		log: logger,
	}
}

type User interface {
	Registration(ctx context.Context, tx *sql.Tx, data user.RegistrationUser) (int64, error)
	IsExistUserByEmail(ctx context.Context, email string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*user.ResponseDataUser, error)
	UpdatePassword(ctx context.Context, email, password string) error
	GetListUser(ctx context.Context, filter user.FilterUser) ([]user.ResponseListDataUser, error)
	GetTotalListUser(ctx context.Context, filter user.FilterUser) (int64, error)
	ChangeStatusUser(ctx context.Context, status bool, userID string) error
	GetUserByUserID(ctx context.Context, userID string) (*user.ResponseDataUser, error)
	GetDetailByUserID(ctx context.Context, userID string) (*user.ResponseDetailDataUser, error)
	UpdateDataUser(ctx context.Context, tx *sql.Tx, data user.UpdateUser) error
	SaveOTP(ctx context.Context, tx *sql.Tx, email string, otpCode int) (int64, error)
	GetOTP(ctx context.Context, email string) (*user.OtpResponse, error)
	UpdateOTP(ctx context.Context, tx *sql.Tx, email string, otpCode int) error
	SaveAudio(ctx context.Context, tx *sql.Tx, userID int64, filename, textData string) (int64, error)
}

const (
	uQCreateUser = `INSERT INTO public.users
	(  
	name, 
	username,
	email, 
	"password", 
	created_at, 
	created_by, 
	"userId")
	VALUES(?,?,?,?,?,?,?) returning id;
	`

	uqUpdateUser = `UPDATE public.users
	SET
	employee_name = ?,
	"position" = ?,
	email = ?,
	partner_id = ?,
	user_role = ?,
	mpi = ?,
	partner = ?,
	product = ?,
	discount = ?,
	setting_product = ?,
	"transaction" = ?,
	billing = ?,
	homepage = ?,
	imin = ?,
	purchase = ?,
	news = ?,
	terms_condition = ?,
	"subscription" = ?,
	refferal_fee = ?,
	platform_fee = ?,
	marketing = ?,
	transaction_qris = ?,
	tax_invoice = ?,
	user_management = ?,
	reporting_transaction = ?,
	reporting_finance = ?,
	subscription_marketing = ?,
	updated_at = ?,
	updated_by = ?,
	phone = ?,
	type_partner = ?
	WHERE
		user_id = ? returning id; 
	`

	uqGetDataUserByEmail = `SELECT id, 
	email, password, name
FROM public.users te where te.email = ? or te.username = ?;
	`

	uqGetDataUserByUserID = `SELECT id, 
	employee_name, 
	"position", 
	email, 
	"password", 
	partner_id, 
	is_product, 
	is_product_child, 
	is_list_user, 
	is_list_user_child, 
	is_inventory, 
	is_inventory_child, 
	is_transaction, 
	is_transaction_child,
	is_active
	FROM public.users te where te.user_id = ?;
	`

	uQCreateOTP = `INSERT INTO public.otp
	(  
	email, 
	otp_code,
	expired_at, 
	created_at)
	VALUES(?,?,?,?) returning id;
	`

	uQCreateHistories = `INSERT INTO public.histories
	(  
	"userId", 
	filename,
	text_data,
	created_at, 
	created_by)
	VALUES(?,?,?,?,?) returning id;
	`

	uQUpdateOTP = `UPDATE public.otp SET otp_code = ?, expired_at = ?, updated_at = ? where email = ?`
)

func (uc UserConfig) Registration(ctx context.Context, tx *sql.Tx, data user.RegistrationUser) (int64, error) {
	param := make([]interface{}, 0)

	param = append(param, data.Name)
	param = append(param, data.Username)
	param = append(param, data.Email)
	param = append(param, data.Password)
	param = append(param, data.CreatedAt)
	param = append(param, data.CreatedBy)
	param = append(param, data.UserID)

	query, args, err := uc.db.Backend.Write.In(uQCreateUser, param...)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	query = uc.db.Backend.Write.Rebind(query)

	var res *sql.Row
	if tx == nil {
		res = uc.db.Backend.Write.QueryRow(ctx, query, args...)
	} else {
		res = tx.QueryRowContext(ctx, query, args...)
	}

	if err != nil {
		return 0, err
	}

	err = res.Err()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	var id int64
	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (uc UserConfig) IsExistUserByEmail(ctx context.Context, email string) (bool, error) {
	var result bool

	uQIsExistUserByEmailActive := ` select exists(select te.id from users te where ( te.email = ? or te.username = ? ) and is_active = true)`

	query, args, err := uc.db.Backend.Read.In(uQIsExistUserByEmailActive, email, email)
	if err != nil {
		return result, err
	}

	query = uc.db.Backend.Read.Rebind(query)
	err = uc.db.Backend.Read.GetContext(ctx, &result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return result, err
	}

	return result, nil
}

func (uc UserConfig) GetUserByEmail(ctx context.Context, email string) (*user.ResponseDataUser, error) {
	var result user.ResponseDataUser

	query, args, err := uc.db.Backend.Read.In(uqGetDataUserByEmail, email, email)
	if err != nil {
		return nil, err
	}

	query = uc.db.Backend.Read.Rebind(query)
	err = uc.db.Backend.Read.GetContext(ctx, &result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &result, nil
}

func (uc UserConfig) UpdatePassword(ctx context.Context, email, password string) error {
	q := `	UPDATE public.users
	SET "password"=?, updated_at=?, updated_by=?
	WHERE email = ?;`

	query, args, err := uc.db.Backend.Read.In(q, password, time.Now().UTC(), "system", email)
	if err != nil {
		return err
	}

	query = uc.db.Backend.Read.Rebind(query)
	res, err := uc.db.Backend.Write.ExecContext(ctx, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		err = fmt.Errorf("no rows inserted")
		return err
	}

	return nil
}

func (uc UserConfig) GetListUser(ctx context.Context, filter user.FilterUser) ([]user.ResponseListDataUser, error) {
	var result []user.ResponseListDataUser

	if filter.Status.Int64 == 2 {
		filter.Statusinit = `te.is_active = true`
	} else if filter.Status.Int64 == 3 {
		filter.Statusinit = `te.is_active = false`
	}

	q := `SELECT te.user_id, te.employee_name, te.email as user_email, te."position", TO_CHAR(te.created_at, 'yyyy-mm-dd hh24:mi:ss') AS date_created,
		CASE
			WHEN te.is_active = true THEN 'Aktif'
			WHEN te.is_active = false THEN 'Tidak Aktif'
		END AS status_name, CASE
			WHEN te.is_active = true THEN 2
			WHEN te.is_active = false THEN 3
		END AS is_active
		FROM users te`

	queryStatement, args2 := buildQueryStatementGetListUser(q, filter)

	query, args, err := uc.db.Backend.Read.In(queryStatement, args2...)
	if err != nil {
		return nil, err
	}

	query = uc.db.Backend.Read.Rebind(query)
	err = uc.db.Backend.Read.Select(&result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return result, nil
}

func (uc UserConfig) GetTotalListUser(ctx context.Context, filter user.FilterUser) (int64, error) {
	var result int64

	q := `select count(te.id) from users te`

	queryStatement, args2 := buildQueryStatementGetTotalListUser(q, filter)

	query, args, err := uc.db.Backend.Read.In(queryStatement, args2...)
	if err != nil {
		return result, err
	}

	query = uc.db.Backend.Read.Rebind(query)
	err = uc.db.Backend.Read.Get(&result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return result, err
	}

	return result, nil
}

func (uc UserConfig) ChangeStatusUser(ctx context.Context, status bool, userID string) error {
	q := `UPDATE public.users
	SET is_active=?
	WHERE user_id = ?;`

	query, args, err := uc.db.Backend.Read.In(q, status, userID)
	if err != nil {
		return err
	}

	query = uc.db.Backend.Read.Rebind(query)
	res, err := uc.db.Backend.Write.ExecContext(ctx, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		err = fmt.Errorf("no rows inserted")
		return err
	}

	return nil
}

func (uc UserConfig) GetUserByUserID(ctx context.Context, userID string) (*user.ResponseDataUser, error) {
	var result user.ResponseDataUser

	query, args, err := uc.db.Backend.Read.In(uqGetDataUserByUserID, userID)
	if err != nil {
		return nil, err
	}

	query = uc.db.Backend.Read.Rebind(query)
	err = uc.db.Backend.Read.GetContext(ctx, &result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &result, nil
}

func (uc UserConfig) GetDetailByUserID(ctx context.Context, userID string) (*user.ResponseDetailDataUser, error) {
	var result user.ResponseDetailDataUser

	q := `select 
	tb.user_id, 
	tb.employee_name, 
	tb."position", 
	tb.email, 
	tb.partner_id, 
	tp.company_name as partner_name,
	tb.user_role, 
	tb.mpi, 
	tb.partner, 
	tb.product, 
	tb.discount, 
	tb.setting_product, 
	tb."transaction", 
	tb.billing, 
	tb.homepage, 
	tb.imin, 
	tb.purchase, 
	tb.news, 
	tb.terms_condition, 
	tb."subscription", 
	tb.refferal_fee, 
	tb.platform_fee, 
	tb.marketing, 
	tb.transaction_qris, 
	tb.tax_invoice, 
	tb.user_management, 
	tb.reporting_transaction, 
	tb.reporting_finance,
	tb.phone,
	tb.type_partner,
	subscription_marketing
	from public.users tb
	left join tbmstr_partner tp on tp.partner_id = tb.partner_id
	where tb.user_id = ?`

	query, args, err := uc.db.Backend.Read.In(q, userID)
	if err != nil {
		return nil, err
	}

	query = uc.db.Backend.Read.Rebind(query)
	err = uc.db.Backend.Read.GetContext(ctx, &result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &result, nil
}

func buildQueryStatementGetListUser(baseQuery string, filter user.FilterUser) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.Email.Valid && filter.Email.String != "" {
		conditions = append(conditions, "te.email ILIKE '%' || ? || '%'")
		args = append(args, filter.Email.String)
	}

	if filter.Status.Valid && (filter.Status.Int64 == 2 || filter.Status.Int64 == 3) {
		conditions = append(conditions, filter.Statusinit)
	}

	if filter.Username.Valid && filter.Username.String != "" {
		conditions = append(conditions, "te.employee_name ILIKE '%' || ? || '%'")
		args = append(args, filter.Username.String)
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
	}

	baseQuery += " ORDER BY te.employee_name ASC"

	if filter.Offset.Valid && filter.Limit.Valid && filter.Offset.Int64 != 0 && filter.Limit.Int64 != 0 {
		baseQuery += " OFFSET ((? - 1) * ?) ROWS FETCH NEXT ? ROWS ONLY"
		args = append(args, filter.Offset.Int64, filter.Limit.Int64, filter.Limit.Int64)
	}

	return baseQuery, args
}

func buildQueryStatementGetTotalListUser(baseQuery string, filter user.FilterUser) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.Email.Valid && filter.Email.String != "" {
		conditions = append(conditions, "te.email ILIKE '%' || ? || '%'")
		args = append(args, filter.Email.String)
	}

	if filter.Status.Valid && (filter.Status.Int64 == 2 || filter.Status.Int64 == 3) {
		conditions = append(conditions, filter.Statusinit)
	}

	if filter.Username.Valid && filter.Username.String != "" {
		conditions = append(conditions, "te.employee_name ILIKE '%' || ? || '%'")
		args = append(args, filter.Username.String)
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
	}

	if filter.Offset.Valid && filter.Limit.Valid && filter.Offset.Int64 != 0 && filter.Limit.Int64 != 0 {
		baseQuery += " OFFSET ((? - 1) * ?) ROWS FETCH NEXT ? ROWS ONLY"
		args = append(args, filter.Offset.Int64, filter.Limit.Int64, filter.Limit.Int64)
	}

	return baseQuery, args
}

func (uc UserConfig) UpdateDataUser(ctx context.Context, tx *sql.Tx, data user.UpdateUser) error {
	param := make([]interface{}, 0)

	param = append(param, data.Name)
	param = append(param, data.Position)
	param = append(param, data.Email)
	param = append(param, data.Partner)
	param = append(param, data.IsUserChild)
	param = append(param, data.IsMpiChild)
	param = append(param, data.IsPartnerChild)
	param = append(param, data.IsProductChild)
	param = append(param, data.IsDiscountChild)
	param = append(param, data.IsSettingProductChild)
	param = append(param, data.IsTransactionChild)
	param = append(param, data.IsBillingChild)
	param = append(param, data.IsHomepageChild)
	param = append(param, data.IsIminChild)
	param = append(param, data.IsPurchaseChild)
	param = append(param, data.IsNewsChild)
	param = append(param, data.IsTermsConditionChild)
	param = append(param, data.IsSubscriptionChild)
	param = append(param, data.IsRefferalFeeChild)
	param = append(param, data.IsPlatformFeeChild)
	param = append(param, data.IsMarketingFeeChild)
	param = append(param, data.IsTransactionQrisChild)
	param = append(param, data.IsTaxInvoiceChild)
	param = append(param, data.IsUserManagementChild)
	param = append(param, data.IsReportingTransactionChild)
	param = append(param, data.IsReportingFinanceChild)
	param = append(param, data.IsSubscriptionMarketingChild)
	param = append(param, data.CreatedAt)
	param = append(param, data.CreatedBy)
	param = append(param, data.Phone)
	param = append(param, data.TypePartner)
	param = append(param, data.UserID)

	query, args, err := uc.db.Backend.Read.In(uqUpdateUser, param...)
	if err != nil {
		return err
	}

	query = uc.db.Backend.Read.Rebind(query)
	res, err := uc.db.Backend.Write.ExecContext(ctx, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		err = fmt.Errorf("no rows inserted")
		return err
	}

	return nil
}

func (uc UserConfig) SaveOTP(ctx context.Context, tx *sql.Tx, email string, otpCode int) (int64, error) {
	param := make([]interface{}, 0)

	param = append(param, email)
	param = append(param, otpCode)
	param = append(param, time.Now().Add(time.Second*60).Format("2006-01-02 15:04:05"))
	param = append(param, time.Now())

	query, args, err := uc.db.Backend.Write.In(uQCreateOTP, param...)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	query = uc.db.Backend.Write.Rebind(query)

	var res *sql.Row
	if tx == nil {
		res = uc.db.Backend.Write.QueryRow(ctx, query, args...)
	} else {
		res = tx.QueryRowContext(ctx, query, args...)
	}

	if err != nil {
		return 0, err
	}

	err = res.Err()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	var id int64
	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (uc UserConfig) GetOTP(ctx context.Context, email string) (*user.OtpResponse, error) {
	var result user.OtpResponse

	q := `SELECT email, otp_code, expired_at FROM otp where email = ?`

	query, args, err := uc.db.Backend.Read.In(q, email)
	if err != nil {
		return nil, err
	}

	query = uc.db.Backend.Read.Rebind(query)
	err = uc.db.Backend.Read.GetContext(ctx, &result, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &result, nil
}

func (uc UserConfig) UpdateOTP(ctx context.Context, tx *sql.Tx, email string, otpCode int) error {
	param := make([]interface{}, 0)

	param = append(param, otpCode)
	param = append(param, time.Now().Add(time.Second*60).Format("2006-01-02 15:04:05"))
	param = append(param, time.Now())
	param = append(param, email)

	query, args, err := uc.db.Backend.Read.In(uQUpdateOTP, param...)
	if err != nil {
		return err
	}

	query = uc.db.Backend.Read.Rebind(query)
	res, err := uc.db.Backend.Write.ExecContext(ctx, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		err = fmt.Errorf("no rows inserted")
		return err
	}

	return nil
}

func (uc UserConfig) SaveAudio(ctx context.Context, tx *sql.Tx, userID int64, filename, textData string) (int64, error) {
	param := make([]interface{}, 0)

	param = append(param, userID)
	param = append(param, filename)
	param = append(param, textData)
	param = append(param, time.Now().Add(time.Second*60).Format("2006-01-02 15:04:05"))
	param = append(param, "system")

	query, args, err := uc.db.Backend.Write.In(uQCreateHistories, param...)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	query = uc.db.Backend.Write.Rebind(query)

	var res *sql.Row
	if tx == nil {
		res = uc.db.Backend.Write.QueryRow(ctx, query, args...)
	} else {
		res = tx.QueryRowContext(ctx, query, args...)
	}

	if err != nil {
		return 0, err
	}

	err = res.Err()
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	var id int64
	err = res.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
