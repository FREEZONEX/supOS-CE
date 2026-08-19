package gateway

import (
	"context"
	"errors"

	"backend/internal/repo"
)

type Service struct {
	repo *repo.GatewayRepo
}

func New(ctx context.Context) *Service {
	return &Service{repo: repo.NewGatewayRepo(ctx)}
}

func (s *Service) Routes(ctx context.Context) ([]repo.GatewayRoute, error) {
	return s.repo.ListGatewayRoutes(ctx)
}

func (s *Service) CreateRoute(ctx context.Context, route repo.GatewayRoute) error {
	existing, err := s.repo.GetGatewayRoute(ctx, route.RouteKey)
	if err == nil {
		if existing.SystemBuiltin {
			return ErrSystemRouteReadonly
		}
		return ErrRouteExists
	}
	if !errors.Is(err, repo.ErrNotFound) {
		return err
	}
	route.SystemBuiltin = false
	return s.saveRoute(ctx, route)
}

func (s *Service) UpdateRoute(ctx context.Context, route repo.GatewayRoute) error {
	existing, err := s.repo.GetGatewayRoute(ctx, route.RouteKey)
	if errors.Is(err, repo.ErrNotFound) {
		return ErrRouteNotFound
	}
	if err != nil {
		return err
	}
	if existing.SystemBuiltin {
		return ErrSystemRouteReadonly
	}
	route.SystemBuiltin = false
	return s.saveRoute(ctx, route)
}

func (s *Service) saveRoute(ctx context.Context, route repo.GatewayRoute) error {
	if route.MatchType == "" {
		route.MatchType = "prefix"
	}
	if route.TimeoutMs == 0 {
		route.TimeoutMs = 10000
	}
	return s.repo.UpsertGatewayRoute(ctx, route)
}

func (s *Service) DeleteRoute(ctx context.Context, routeKey string, userID int64) error {
	existing, err := s.repo.GetGatewayRoute(ctx, routeKey)
	if errors.Is(err, repo.ErrNotFound) {
		return ErrRouteNotFound
	}
	if err != nil {
		return err
	}
	if existing.SystemBuiltin {
		return ErrSystemRouteReadonly
	}
	return s.repo.DeleteGatewayRoute(ctx, routeKey, userID)
}

var (
	ErrRouteNotFound       = errors.New("gateway route not found")
	ErrRouteExists         = errors.New("gateway route already exists")
	ErrSystemRouteReadonly = errors.New("system gateway route is read-only")
)
