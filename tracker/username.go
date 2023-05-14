package tracker

import "context"

type trackerCtxKey string

const usernameKey trackerCtxKey = "username"

func (c trackerCtxKey) String() string {
	return string(c)
}

func PutUsername(ctx context.Context, username string) context.Context {
	ctx = context.WithValue(ctx, usernameKey, username)
	return ctx
}

func GetUsername(ctx context.Context) string {
	return ctx.Value(usernameKey).(string)
}
