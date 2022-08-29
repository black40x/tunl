package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/black40x/tunl-core/commands"
	"github.com/black40x/tunl-core/tunl"
	"github.com/urfave/cli/v2"
	"os"
	"os/signal"
	"strings"
	"time"
	"tunl-cli/cmd/client"
	"tunl-cli/cmd/monitor"
	"tunl-cli/cmd/options"
	"tunl-cli/cmd/tui"
)

func StartTunlClient(opt *options.Options) error {
	tunlClient := client.NewTunlClient(opt, client.Version)
	err := tunlClient.Connect()
	if err != nil {
		return err
	}

	tunlMonitor := monitor.NewTunlMonitor()

	tunlClient.SetHttpRequestReceiver(func(r *commands.HttpRequest, b []byte, d time.Duration) {
		mes := map[string]interface{}{
			"duration":       d.Milliseconds(),
			"uuid":           r.GetUuid(),
			"method":         r.GetMethod(),
			"proto":          r.GetProto(),
			"uri":            r.GetUri(),
			"remote_address": r.GetRemoteAddr(),
			"header":         r.Header,
			"cookies":        r.Cookies,
			"body_type":      false,
			"body":           false,
		}

		if r.IsJson() {
			mes["body_type"] = "json"
			mes["body"] = string(b)
		} else if _, ok := r.IsFormData(); ok {
			buf := bytes.NewBuffer(b)
			files := map[string]string{}
			form, err := r.ParseFormData(buf)
			if err == nil {
				mes["body_type"] = "form-data"

				for f, _ := range form.File {
					files[f] = "[binary]"
				}

				mes["body"] = map[string]interface{}{
					"values": form.Value,
					"files":  files,
				}
			}
		}

		tunlMonitor.SendJsonMessage(mes)
	})

	if opt.Monitor {
		err = tunlMonitor.Start(opt.MonitorAddr.ToString())
		if err != nil {
			return errors.New(fmt.Sprintf("can't start tunl monitor at %s", opt.MonitorAddr.ToString()))
		}
	}

	// Wait for exit
	var wait time.Duration
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	tui.PrintInfo("Shutting down tunl cli")

	tunlMonitor.Shutdown(ctx)
	tunlClient.Close()

	return nil
}

func app() {
	app := &cli.App{
		Name:      "tunl",
		Usage:     "make your localhost visible outside and inspect traffic",
		UsageText: "tunl [command] [flag]",
		Version:   client.Version,
		Commands: []*cli.Command{
			{
				Name:        "http",
				Usage:       "start HTTP tunnel",
				Description: "Start tunnel from localhost to the http://<subdomain>.tunl.online,\nand run Traffic Monitor.",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "tunl-host",
						Usage:       "set tunl connection host",
						DefaultText: "tunl.online",
						Value:       "tunl.online",
					},
					&cli.StringFlag{
						Name:        "tunl-port",
						Usage:       "set tunl connection port",
						DefaultText: "5000",
						Value:       "5000",
					},
					&cli.BoolFlag{
						Name:        "monitor",
						Usage:       "enable/disable Traffic monitor",
						DefaultText: "true",
						Value:       true,
					},
					&cli.StringFlag{
						Name:        "monitor-host",
						Usage:       "set Traffic monitor host",
						DefaultText: "127.0.0.1",
						Value:       "127.0.0.1",
					},
					&cli.StringFlag{
						Name:        "monitor-port",
						Usage:       "set Traffic monitor port",
						DefaultText: "6060",
						Value:       "6060",
					},
					&cli.StringFlag{
						Name:  "server-pass",
						Usage: "set server authenticate password",
					},
					&cli.StringFlag{
						Name:  "basic-auth",
						Usage: "use Basic Auth 'user:password'",
					},
					&cli.StringSliceFlag{
						Name:  "req-header",
						Usage: "add 'key:value' to request header",
					},
					&cli.StringSliceFlag{
						Name:  "resp-header",
						Usage: "add 'key:value' to response header",
					},
				},
				Action: func(cCtx *cli.Context) error {
					localAddr, err := tunl.NewAddress(cCtx.Args().First())
					if err != nil {
						return err
					}

					monitorAddr, err := tunl.NewAddress(cCtx.String("monitor-host") + ":" + cCtx.String("monitor-port"))
					if err != nil {
						return err
					}

					opt := &options.Options{
						LocalAddr:       localAddr,
						ServerAddr:      cCtx.String("tunl-host") + ":" + cCtx.String("tunl-port"),
						HttpTimeout:     time.Second * 15,
						Monitor:         cCtx.Bool("monitor"),
						MonitorAddr:     monitorAddr,
						RequestHeaders:  client.ArrToHeaders(cCtx.StringSlice("req-header"), "="),
						ResponseHeaders: client.ArrToHeaders(cCtx.StringSlice("resp-header"), "="),
					}
					if cCtx.String("basic-auth") != "" {
						up := strings.Split(cCtx.String("basic-auth"), ":")
						if len(up) >= 2 {
							opt.BasicAuth = &options.BasicAuth{
								Login: up[0],
								Pass:  up[1],
							}
						} else {
							return errors.New("basic auth invalid login:password")
						}
					}

					return StartTunlClient(opt)
				},
			},
			{
				Name:  "version",
				Usage: "print tunl cli version",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		tui.PrintError(err)
		os.Exit(1)
	}
}

func main() {
	app()
}
