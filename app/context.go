package app

import (
	"github.com/go-macaron/cache"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

type Context struct {
	*macaron.Context
	Cache    cache.Cache
	Session  session.Store
	UserName string
	IsSigned bool
}

func SignedIn(session session.Store) (bool, string) {
	seUserName := session.Get(SessionKey)
	if seUserName != nil {
		return true, seUserName.(string)
	}
	return false, ""
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context, cache cache.Cache, sess session.Store) {
		ctx := &Context{
			Context: c,
			Cache:   cache,
			Session: sess,
		}

		has, username := SignedIn(ctx.Session)
		if has {
			ctx.IsSigned = true
			ctx.Data["IsSigned"] = ctx.IsSigned
			ctx.Data["SignedUser"] = username
		}
		c.Map(ctx)
	}
}
