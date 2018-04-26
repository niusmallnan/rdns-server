package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/niusmallnan/rdns-server/backend"
	"github.com/niusmallnan/rdns-server/backend/etcd"
	"github.com/niusmallnan/rdns-server/service"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
			Name:   "root-domain",
			Value:  "lb.rancher.cloud",
			EnvVar: "RANCHER_ROOT_DOMAIN",
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
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.Infof("Using backend: %s", ctx.String("backend"))
	var (
		b   backend.Backend
		err error
	)
	switch ctx.String("backend") {
	case etcd.BackendName:
		etcdEndpoints := strings.Split(ctx.String("etcd-endpoints"), ",")
		b, err = etcd.NewEtcdBackend(etcdEndpoints, ctx.String("etcd-prepath"), ctx.String("root-domain"))
	default:
		err = errors.Errorf("Failed to find backend %s", ctx.String("backend"))
	}
	if err != nil {
		return errors.Wrapf(err, "Failed to init backend %s", ctx.String("backend"))
	}
	backend.SetBackend(b)

	done := make(chan error)

	go func() {
		router := service.NewRouter()
		done <- http.ListenAndServe(ctx.String("listen"), router)
	}()

	return <-done
}
