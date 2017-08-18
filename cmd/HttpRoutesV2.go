package cmd

import (
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
	"github.com/gityflow/githorse/routes.v2"
)

func routesV2(c *cli.Context,m *macaron.Macaron, reqSignIn, ignSignIn, ignSignInAndCsrf, reqSignOut macaron.Handler){
	// FIXME: not all routes need go through same middlewares.
	// Especially some AJAX requests, we can reduce middleware number to improve performance.
	// Routers.
	m.Get("/v2", ignSignIn, routes_v2.Home)
}
