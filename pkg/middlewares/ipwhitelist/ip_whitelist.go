package ipallowlist

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	"github.com/traefik/traefik/v2/pkg/ip"
	"github.com/traefik/traefik/v2/pkg/log"
	"github.com/traefik/traefik/v2/pkg/middlewares"
	"github.com/traefik/traefik/v2/pkg/tracing"
)

const (
	typeName = "IPAllowLister"
)

// ipAllowLister is a middleware that provides Checks of the Requesting IP against a set of Whitelists.
type ipAllowLister struct {
	next        http.Handler
	allowLister *ip.Checker
	strategy    ip.Strategy
	name        string
}

// New builds a new IPAllowLister given a list of CIDR-Strings to allowlist.
func New(ctx context.Context, next http.Handler, config dynamic.IPAllowList, name string) (http.Handler, error) {
	logger := log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName))
	logger.Debug("Creating middleware")

	if len(config.SourceRange) == 0 {
		return nil, errors.New("sourceRange is empty, IPAllowLister not created")
	}

	checker, err := ip.NewChecker(config.SourceRange)
	if err != nil {
		return nil, fmt.Errorf("cannot parse CIDR allowlist %s: %w", config.SourceRange, err)
	}

	strategy, err := config.IPStrategy.Get()
	if err != nil {
		return nil, err
	}

	logger.Debugf("Setting up IPAllowLister with sourceRange: %s", config.SourceRange)

	return &ipAllowLister{
		strategy:    strategy,
		allowLister: checker,
		next:        next,
		name:        name,
	}, nil
}

func (wl *ipAllowLister) GetTracingInformation() (string, ext.SpanKindEnum) {
	return wl.name, tracing.SpanKindNoneEnum
}

func (wl *ipAllowLister) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := middlewares.GetLoggerCtx(req.Context(), wl.name, typeName)
	logger := log.FromContext(ctx)

	err := wl.allowLister.IsAuthorized(wl.strategy.GetIP(req))
	if err != nil {
		logMessage := fmt.Sprintf("rejecting request %+v: %v", req, err)
		logger.Debug(logMessage)
		tracing.SetErrorWithEvent(req, logMessage)
		reject(ctx, rw)
		return
	}
	logger.Debugf("Accept %s: %+v", wl.strategy.GetIP(req), req)

	wl.next.ServeHTTP(rw, req)
}

func reject(ctx context.Context, rw http.ResponseWriter) {
	statusCode := http.StatusForbidden

	rw.WriteHeader(statusCode)
	_, err := rw.Write([]byte(http.StatusText(statusCode)))
	if err != nil {
		log.FromContext(ctx).Error(err)
	}
}
