package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sunshineplan/metadata"
	"github.com/sunshineplan/password"
	"github.com/sunshineplan/service"
	"github.com/sunshineplan/stock"
	_ "github.com/sunshineplan/stock/eastmoney"
	"github.com/sunshineplan/utils"
	"github.com/sunshineplan/utils/flags"
	"github.com/sunshineplan/utils/httpsvr"
)

var (
	self string
	priv *rsa.PrivateKey

	server = httpsvr.New()
	svc    = service.New()
	meta   metadata.Server

	joinPath = filepath.Join
	dir      = filepath.Dir
)

func init() {
	var err error
	self, err = os.Executable()
	if err != nil {
		svc.Fatalln("Failed to get self path:", err)
	}
	svc.Name = "MyStocks"
	svc.Desc = "Instance to serve My Stocks"
	svc.Exec = run
	svc.TestExec = test
	svc.Options = service.Options{
		Dependencies:       []string{"After=network.target"},
		Environment:        map[string]string{"GIN_MODE": "release"},
		RemoveBeforeUpdate: []string{"dist/assets"},
		ExcludeFiles:       []string{"scripts/mystocks.conf"},
	}
	svc.RegisterCommand("add", "add user", func(arg ...string) error {
		return addUser(arg[0])
	}, 1)
	svc.RegisterCommand("delete", "delete user", func(arg ...string) error {
		if utils.Confirm("Do you want to delete this user?", 3) {
			return deleteUser(arg[0])
		}
		return nil
	}, 1)

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(flag.CommandLine.Output(), svc.Usage())
	}
}

var (
	universal = flag.Bool("universal", false, "Use Universal account id or not")
	maxRetry  = flag.Int("retry", 5, "Max number of retries on wrong password")
	refresh   = flag.Int("refresh", 3, "Refresh Interval")
	pemPath   = flag.String("pem", "", "PEM file Path")
	logPath   = flag.String("log", "", "Log file path")
)

func main() {
	flag.StringVar(&meta.Addr, "server", "", "Metadata Server Address")
	flag.StringVar(&meta.Header, "header", "", "Verify Header Header Name")
	flag.StringVar(&meta.Value, "value", "", "Verify Header Value")
	flag.StringVar(&server.Unix, "unix", "", "UNIX-domain Socket")
	flag.StringVar(&server.Host, "host", "0.0.0.0", "Server Host")
	flag.StringVar(&server.Port, "port", "12345", "Server Port")
	flag.StringVar(&svc.Options.UpdateURL, "update", "", "Update URL")
	flag.StringVar(&svc.Options.PIDFile, "pid", "/var/run/mystocks.pid", "PID file path")
	flags.SetConfigFile(joinPath(dir(self), "config.ini"))
	flags.Parse()

	password.SetMaxAttempts(*maxRetry)
	if *pemPath != "" {
		b, err := os.ReadFile(*pemPath)
		if err != nil {
			svc.Fatal(err)
		}
		block, _ := pem.Decode(b)
		if block == nil {
			svc.Fatal("no PEM data is found")
		}
		priv, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			svc.Fatal(err)
		}
	}
	stock.SetTimeout(*refresh)

	if err := svc.ParseAndRun(flag.Args()); err != nil {
		svc.Fatal(err)
	}
}
