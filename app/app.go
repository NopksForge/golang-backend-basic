package app

type ctxKey string

const (
	refIDKey      ctxKey = "ref-id"
	forwardCtxKey ctxKey = "forwarding"

	RedisProductKey string = "product"
	RedisUserKey    string = "user"
)
