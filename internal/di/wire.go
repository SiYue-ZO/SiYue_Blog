//go:build wireinject
// +build wireinject

package di

import (
	"net/http"

	"github.com/google/wire"
)

func InitHTTPServer() (*http.Server, error) {
	wire.Build(AppSet)
	return nil, nil
}
