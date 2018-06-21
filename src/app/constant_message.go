package app

// app sinup message
const (
	ErrCheckAppName   string = "err_check_app_name"
	ErrAppNameExisted        = "err_app_name_existed"
	ErrSaveApp               = "err_save_app"
	SaveAppSuccess           = "save_app_success"
)

// db associate message
const (
	ErrQueryUser       string = "err_query_user"
	ErrCheckUserName          = "err_check_user_name"
	ErrUserNotExisted         = "err_user_not_existed"
	ErrUserNotVerified        = "err_user_not_verified"
	ErrUserUpdate             = "err_user_update"
	ErrUserExisted            = "err_user_existed"
	ErrUserVerified           = "err_user_verified"
	ErrUserInsert             = "err_user_insert"
	ErrSessionUpdate          = "err_session_update"
)

// http message
const (
	ErrParseQuery               string = "err_parse_query"
	ErrCookieNoUserId                  = "err_cookie_no_user_id"
	ErrNoVerificationCode              = "err_no_verification_code"
	ErrVerificationCodeNotEqual        = "err_verification_code_not_equal"
)

// user message
const (
	UserLoginSuccess    string = "user_login_success"
	UserVerifiedSuccess        = "user_verified_success"
)

// DefaultMessageMap constant message map
type DefaultMessageMap struct {
	msgmap map[string]string
}

// NewDefaultMessageMap return constant message map with default value
func NewDefaultMessageMap() *DefaultMessageMap {
	mp := &DefaultMessageMap{msgmap: make(map[string]string)}
	mp.msgmap[ErrCheckAppName] = "err_check_app_name"
	mp.msgmap[ErrAppNameExisted] = "err_app_name_existed"
	mp.msgmap[ErrSaveApp] = "err_save_app"
	mp.msgmap[SaveAppSuccess] = "save_app_success"
	mp.msgmap[ErrQueryUser] = "err_query_user"
	mp.msgmap[ErrCheckUserName] = "err_check_user_name"
	mp.msgmap[ErrUserNotExisted] = "err_user_not_existed"
	mp.msgmap[ErrUserNotVerified] = "err_user_not_verified"
	mp.msgmap[ErrUserUpdate] = "err_user_update"
	mp.msgmap[ErrUserExisted] = "err_user_existed"
	mp.msgmap[ErrUserVerified] = "err_user_verified"
	mp.msgmap[ErrUserInsert] = "err_user_insert"
	mp.msgmap[ErrSessionUpdate] = "err_session_update"
	mp.msgmap[ErrParseQuery] = "err_parse_query"
	mp.msgmap[ErrCookieNoUserId] = "err_cookie_no_user_id"
	mp.msgmap[ErrNoVerificationCode] = "err_no_verification_code"
	mp.msgmap[ErrVerificationCodeNotEqual] = "err_verification_code_not_equal"
	mp.msgmap[UserLoginSuccess] = "user_login_success"
	mp.msgmap[UserVerifiedSuccess] = "user_verified_success"

	return mp
}

// Get return value by id
func (mp *DefaultMessageMap) Get(id string) string {
	if m, ok := mp.msgmap[id]; ok {
		return m
	}
	return id
}
