package main

import (
	"kiramishima/m-backend/bootstrap"

	"go.uber.org/fx"
)

func main() {
	fx.New(bootstrap.Module).Run()
}
