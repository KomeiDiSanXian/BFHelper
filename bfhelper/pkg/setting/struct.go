package setting

// AccountSettingS 账号相关设置
type AccountSettingS struct {
	Username string
	Password string
	Session  string // Session is X-Gatewaysession
	Token    string // Token is bearerAccessToken
	SID      string // SID is cookie sid
	Remid    string // Remid is cookie remid
}

// SessionAPISettingS 获取Session的API设置
type SessionAPISettingS struct {
	SakuraID    string
	SakuraToken string
	MFASecret   string
}

// BFEACSettingS BFEAC 设置
type BFEACSettingS struct {
	APIKey string
}

// TraceSettingS 追踪设置
type TraceSettingS struct {
	Enabled  bool
	UseHTTPS bool
	URL      string
}
