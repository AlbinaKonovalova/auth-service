package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	errCookieMissing   = errors.New("cookie missing")
	errCookieMalformed = errors.New("cookie malformed")
)

func (c *AuthServiceController) setRefreshCookie(ctx context.Context, raw string) error {
	cookie := &http.Cookie{
		Name:     c.cookie.Name,
		Value:    raw,
		Path:     c.cookie.Path,
		Domain:   c.cookie.Domain,
		HttpOnly: true,
		Secure:   c.cookie.IsSecure(),
		SameSite: parseSameSite(c.cookie.SameSite),
		MaxAge:   int(c.cookie.TTL.Seconds()),
	}
	return grpc.SetHeader(ctx, metadata.Pairs("set-cookie", cookie.String()))
}

func (c *AuthServiceController) clearRefreshCookie(ctx context.Context) error {
	cookie := &http.Cookie{
		Name:     c.cookie.Name,
		Value:    "",
		Path:     c.cookie.Path,
		Domain:   c.cookie.Domain,
		HttpOnly: true,
		Secure:   c.cookie.IsSecure(),
		SameSite: parseSameSite(c.cookie.SameSite),
		MaxAge:   -1,
	}
	return grpc.SetHeader(ctx, metadata.Pairs("set-cookie", cookie.String()))
}

func (c *AuthServiceController) readRefreshCookie(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("%w: no incoming metadata", errCookieMalformed)
	}

	cookies := md.Get("cookie")
	if len(cookies) == 0 {
		return "", errCookieMissing
	}

	header := http.Header{}
	for _, v := range cookies {
		header.Add("Cookie", v)
	}
	r := &http.Request{Header: header}

	cookie, err := r.Cookie(c.cookie.Name)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return "", errCookieMissing
		}
		return "", fmt.Errorf("%w: %w", errCookieMalformed, err)
	}

	return cookie.Value, nil
}

func parseSameSite(s string) http.SameSite {
	switch strings.ToLower(s) {
	case "lax":
		return http.SameSiteLaxMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteStrictMode
	}
}
