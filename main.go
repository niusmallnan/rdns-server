package main

import (
	"net/http"
	"os"

	"github.com/niusmallnan/rdns-server/backend"
	"github.com/niusmallnan/rdns-server/backend/etcd"
	"github.com/niusmallnan/rdns-server/service"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var VERSION = "v0.0.0-dev"

func main() {
	app := cli.NewApp()
	app.Name = "rdns-server"
	app.Version = VERSION
	app.Usage = "You need help!"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug, d",
			EnvVar: "RANCHER_DEBUG",
		},
		cli.StringFlag{
			Name:   "listen",
			Value:  ":9333",
			EnvVar: "RANCHER_SERVICE_LISTEN_PORT",
		},
		cli.StringFlag{
			Name:   "backend",
			Value:  "etcd",
			EnvVar: "RANCHER_SERVICE_BACKEND",
		},
		cli.StringFlag{
			Name:   "etcd-endpoints",
			Value:  "http://127.0.0.1:2379",
			EnvVar: "RANCHER_ETCD_ENDPOINTS",
		},
		cli.StringFlag{
			Name:   "etcd-prepath",
			Value:  "/rdns",
			EnvVar: "RANCHER_ETCD_PREPATH",
		},
	}
	app.Action = func(ctx *cli.Context) {
		if err := appMain(ctx); err != nil {
			logrus.Errorf("error: %v", err)
			os.Exit(1)
		}
	}

	app.Run(os.Args)
}

func appMain(ctx *cli.Context) error {
	if ctx.Bool("debug") {
		logrus.SetLevelString("debug")
	}

	var (
		b     *backend.Backend
		error err
	)
	switch name {
	case etcd.ETCD_BACKEND:
		b, err = etcd.NewEtcdBackend(ctx.String("etcd-endpoints"), ctx.String("etcd-prepath"))
	}
	if err != nil {
		return errors.Wrapf(err, "Failed to init backend %s", ctx.String("backend"))
	}
	backend.SetBackend(b)

	done := make(chan error)

	go func() {
		router := service.NewRouter()
		done <- http.ListenAndServe(ctx.String(Listen), router)
	}()

	return <-done
}
