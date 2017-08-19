// +build go1.6

// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Gogs is a painless self-hosted Git Service.
package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/gityflow/githorse/cmd"
	"github.com/gityflow/githorse/pkg/setting"
)

const APP_VER = "0.0.1"

func init() {
	setting.AppVer = APP_VER
}

func main() {
	app := cli.NewApp()
	app.Name = "GitHorse"
	app.Usage = "A painless self-hosted Git service"
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.Web,
		cmd.Serv,
		cmd.Hook,
		cmd.Cert,
		cmd.Admin,
		cmd.Import,
		cmd.Backup,
		cmd.Restore,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}