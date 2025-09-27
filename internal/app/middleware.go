package app

import "github.com/Mafit1/notes-app/internal/api/common/middleware"

func (app *App) AuthMW() *middleware.AuthMW {
	if app.authMW == nil {
		app.authMW = middleware.NewAuthMW(app.JwtService())
	}
	return app.authMW
}

func (app *App) MetricsMW() *middleware.MetricsMW {
	if app.metricsMW == nil {
		app.metricsMW = middleware.NewMetricsMW()
	}
	return app.metricsMW
}
