package auth

// Claims 表示 JWT 令牌的声明信息
type Claims struct {
	UserID uint     `json:"user_id"` // 用户ID
	Email  string   `json:"email"`   // 用户邮箱
	Name   string   `json:"name"`    // 用户姓名
	Roles  []string `json:"roles"`   // 用户角色列表
}

// TokenResponse 表示令牌响应（已废弃：请使用 TokenPairResponse）
type TokenResponse struct {
	Token string `json:"token"` // JWT 令牌
}

// TokenPairResponse 表示访问令牌和刷新令牌对的响应
type TokenPairResponse struct {
	AccessToken  string `json:"access_token"`  // 访问令牌
	RefreshToken string `json:"refresh_token"` // 刷新令牌
	TokenType    string `json:"token_type"`    // 令牌类型（通常为 "Bearer"）
	ExpiresIn    int64  `json:"expires_in"`    // 访问令牌过期时间（秒）
}

// RefreshTokenRequest 表示刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // 刷新令牌
}
