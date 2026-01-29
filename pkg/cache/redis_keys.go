package cache

const (
	RedisKeyAuthBlacklistPrefix      = "auth:blacklist:"
	RedisKeyPasswordResetEmailPrefix = "auth:password_reset:email:"
	RedisKeyPasswordResetTokenPrefix = "auth:password_reset:token:"
)

func BlacklistAccessTokenKey(jti string) string {
	return RedisKeyAuthBlacklistPrefix + jti
}

func PasswordResetEmailRateLimitKey(email string) string {
	return RedisKeyPasswordResetEmailPrefix + email
}

func PasswordResetTokenKey(uuid string) string {
	return RedisKeyPasswordResetTokenPrefix + uuid
}
