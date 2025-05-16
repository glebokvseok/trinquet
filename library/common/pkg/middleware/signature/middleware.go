package signmw

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/move-mates/trinquet/library/common/pkg/constants"
	"github.com/move-mates/trinquet/library/common/pkg/extensions"
	"go.uber.org/fx"
	"strconv"
	"strings"
)

type SignatureHandlerMiddleware struct {
	fx.Out

	Middleware echo.MiddlewareFunc `name:"signature_handler"`
}

func provideSignatureHandlerMiddleware(
	config signatureConfig,
) SignatureHandlerMiddleware {
	return SignatureHandlerMiddleware{
		Middleware: signatureHandler(config),
	}
}

func signatureHandler(config signatureConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			signature := ctx.Request().Header.Get(constants.SignatureHeader)
			if extensions.IsEmpty(signature) {
				return NewMalformedOrExpiredSignatureError()
			}

			timestamp, err := strconv.ParseInt(ctx.Request().Header.Get(constants.TimestampHeader), 10, 64)
			if err != nil {
				return NewMalformedOrExpiredSignatureError()
			}

			expectedSignature := generateSignature(config.RequestSigningKey, ctx.Request().Method, getFullUrl(ctx), timestamp)
			if signature != expectedSignature {
				return NewMalformedOrExpiredSignatureError()
			}

			return next(ctx)
		}
	}
}

func generateSignature(signingKey string, method string, url string, timestamp int64) string {
	dataToSign := fmt.Sprintf("%s|%s|%d", method, url, timestamp)

	h := hmac.New(sha256.New, []byte(signingKey))
	h.Write([]byte(dataToSign))

	return hex.EncodeToString(h.Sum(nil))
}

func getFullUrl(ctx echo.Context) string {
	req := ctx.Request()

	scheme := req.Header.Get(constants.ForwardedProtoHeader)
	if extensions.IsEmpty(scheme) {
		if req.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	host := req.Header.Get(constants.ForwardedHostHeader)
	if extensions.IsEmpty(host) {
		host = req.Host
	}

	var builder strings.Builder
	builder.WriteString(scheme)
	builder.WriteString("://")
	builder.WriteString(host)
	builder.WriteString(req.URL.Path)

	if extensions.IsNotEmpty(req.URL.RawQuery) {
		builder.WriteString("?")
		builder.WriteString(req.URL.RawQuery)
	}

	return builder.String()
}
