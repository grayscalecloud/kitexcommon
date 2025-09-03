package tools

import (
	"context"
	"net/http"

	pt "aidanwoods.dev/go-paseto"
	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/grayscalecloud/kitexcommon/consts"
	"github.com/grayscalecloud/kitexcommon/consts/errno"
	"github.com/grayscalecloud/kitexcommon/model"
	"github.com/hertz-contrib/gzip"
	"github.com/hertz-contrib/paseto"
)

func CommonMW() []app.HandlerFunc {
	return []app.HandlerFunc{
		// use cors mw
		//middleware.Cors(),
		// use recovery mw
		//middleware.Recovery(),
		// use gzip mw
		gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{".jpg", ".mp4", ".png"})),
	}
}

func PasetoAuth(audience string, pi model.PasetoConfig) app.HandlerFunc {

	pf, err := paseto.NewV4PublicParseFunc(pi.PubKey, []byte(pi.Implicit), paseto.WithAudience(audience), paseto.WithNotBefore())
	if err != nil {
		hlog.Fatal(err)
	}
	sh := func(ctx context.Context, c *app.RequestContext, token *pt.Token) {
		aid, err := token.GetString("id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, BuildBaseResp(errno.BadRequest.WithMessage("missing accountID in token")))
			c.Abort()
			return
		}
		// 将 AccountID 存储在 context 中
		ctx = context.WithValue(ctx, consts.AccountID, aid)
		ctx = metainfo.WithValue(ctx, consts.AccountID, aid)
		c.Next(ctx)
	}

	eh := func(ctx context.Context, c *app.RequestContext) {
		c.JSON(http.StatusUnauthorized, BuildBaseResp(errno.BadRequest.WithMessage("invalid token")))
		c.Abort()
	}
	return paseto.New(paseto.WithParseFunc(pf), paseto.WithSuccessHandler(sh), paseto.WithErrorFunc(eh))
}

func PasetoAuthNotForce(audience string, pi model.PasetoConfig) app.HandlerFunc {
	pf, err := paseto.NewV4PublicParseFunc(pi.PubKey, []byte(pi.Implicit), paseto.WithAudience(audience), paseto.WithNotBefore())
	if err != nil {
		hlog.Fatal(err)
	}
	sh := func(ctx context.Context, c *app.RequestContext, token *pt.Token) {
		aid, err := token.GetString("id")
		if err != nil {
			c.JSON(http.StatusUnauthorized, BuildBaseResp(errno.BadRequest.WithMessage("missing accountID in token")))
			c.Abort()
			return
		}
		// 将 AccountID 存储在 context 中
		ctx = context.WithValue(ctx, consts.AccountID, aid)
		ctx = metainfo.WithValue(ctx, consts.AccountID, aid)
		c.Next(ctx)
	}

	eh := func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
	}
	return paseto.New(paseto.WithParseFunc(pf), paseto.WithSuccessHandler(sh), paseto.WithErrorFunc(eh))
}
