package cache

const (
	RedisKeyAuthBlacklistPrefix = "auth:blacklist:"
)

func BlacklistAccessTokenKey(jti string) string {
	return RedisKeyAuthBlacklistPrefix + jti
}
