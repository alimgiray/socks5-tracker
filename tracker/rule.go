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

	// Check global usage
	if uint64(r.tracker.globalLimit) < r.tracker.globalUsage {
		log.Println("Global usage limit exceeded")
		return ctx, false
	}

	// Check user usage
	user := req.AuthContext.Payload["username"]
	if uint64(r.tracker.perUserLimit) < r.tracker.perUserUsage[user] {
		log.Printf("User %s usage limit exceeded\n", user)
		return ctx, false
	}

	return ctx, true
}
