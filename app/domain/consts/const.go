package consts

const (
	SERVICE_BASE_RDB_KEY  = "user.services:"
	SEND_SMS_BASE_RDB_KER = SERVICE_BASE_RDB_KEY + "send_sms:"
	USER_INFO_RDB_KEY     = SERVICE_BASE_RDB_KEY + "user_info:"
	JURIST_INFO_RDB_KEY   = SERVICE_BASE_RDB_KEY + "jurist_info:"
	RDB_TIMEOUT           = 600
)

const (
	// USERTYPE 1、普通用户，2、作者，3、认证律师. 100、管理员
	USERTYPE          = ""
	USERTYPE_ORDINARY = USERTYPE + "ordinary"
	USERTYPE_AUTHOR   = USERTYPE + "author"
	USERTYPE_JURIST   = USERTYPE + "jurist"
	USERTYPE_ADMIN    = USERTYPE + "admin"
)
