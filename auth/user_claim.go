package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/labstack/echo"
)

type UserClaim struct {
	Id            int64
	Iss           string
	ColleagueNo   string
	ColleagueName string
	TenantId      int64
	TenantCode    string
	ChannelId     int64
	ChannelCode   string
	UserName      string
}

func UserClaimMiddelware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		userClaim := UserClaim{}

		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return next(c)
			}

			si := strings.Index(token, ".")
			li := strings.LastIndex(token, ".")
			if si == -1 || li == -1 || si == li {
				return next(c)
			}

			payload := token[si+1 : li]
			if payload == "" {
				return next(c)
			}

			payloadBytes, err := decodeSegment(payload)
			if err != nil {
				return next(c)
			}

			err = json.Unmarshal(payloadBytes, &userClaim)
			if err != nil {
				return next(c)
			}

			req := c.Request()

			req = req.WithContext(context.WithValue(req.Context(), "userClaim", userClaim))
			req = req.WithContext(context.WithValue(req.Context(), "token", token))

			c.SetRequest(req)

			return next(c)
		}
	}
}

func decodeSegment(seg string) ([]byte, error) {
	if l := len(seg) % 4; l > 0 {
		seg += strings.Repeat("=", 4-l)
	}

	return base64.URLEncoding.DecodeString(seg)
}

func (UserClaim) FromCtx(ctx context.Context) UserClaim {
	v := ctx.Value("userClaim")
	if v == nil {
		return UserClaim{}
	}
	userClaim, ok := v.(UserClaim)
	if !ok {
		return UserClaim{}
	}
	return userClaim
}
