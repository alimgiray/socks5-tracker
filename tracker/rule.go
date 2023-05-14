package tracker

import (
	"context"
	"log"

	"github.com/things-go/go-socks5"
)

type usageLimitRule struct {
	tracker *usageTracker
}

func (r *usageLimitRule) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	if r.tracker.HasGlobalLimitExceeded() {
		log.Println("Global usage limit exceeded")
		return ctx, false
	}

	user := req.AuthContext.Payload[usernameKey.String()]
	if r.tracker.HasUserLimitExceeded(user) {
		log.Printf("User %s usage limit exceeded\n", user)
		return ctx, false
	}

	return ctx, true
}
