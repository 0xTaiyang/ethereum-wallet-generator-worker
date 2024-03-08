package main

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/seth-shi/ethereum-wallet-generator-nodes/internal"
	"github.com/urfave/cli/v2"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var (
	masterCommand = &cli.Command{
		Name:  "master",
		Usage: "启动 HTTP 服务器, 收集信息",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "port",
				Value: 8000,
				Usage: "HTTP 服务端口",
			},

			&cli.StringFlag{
				Name:  "prefix",
				Value: "",
				Usage: "钱包地址前缀",
			},
			&cli.StringFlag{
				Name:  "suffix",
				Value: "",
				Usage: "钱包地址后缀",
			},
			&cli.StringFlag{
				Name:  "key",
				Value: "",
				Usage: "通讯 key,如果不指定,那么自动生成",
			},
		},
		Before: func(cCtx *cli.Context) (err error) {

			var (
				prefix = cCtx.String("prefix")
				suffix = cCtx.String("suffix")
				key    = strings.TrimSpace(cCtx.String("key"))
				port   = cCtx.Int("port")
			)
			if cCtx.String("prefix") == "" && cCtx.String("suffix") == "" {
				return errors.New("钱包前缀和后缀不能同时为空")
			}

			if Master, err = internal.NewMaster(port, prefix, suffix); err != nil {
				return
			}

			if key != "" {
				Master.Config.Key = key
			}

			return nil
		},
		Action: func(cCtx *cli.Context) error {

			defer func() {
				Master.FilePoint.Close()
			}()

			go Master.StartWebServer()
			Master.Run()

			return nil
		},
	}
	nodeCommand = &cli.Command{
		Name:  "node",
		Usage: "生成节点",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "server",
				Value: "http://localhost:8000",
				Usage: "服务器地址",
			},
			&cli.UintFlag{
				Name:  "c",
				Value: 0,
				Usage: "并发线程, 默认 CPU 个数",
			},
		},
		Before: func(cCtx *cli.Context) error {
			var (
				c          = cCtx.Uint("c")
				serverHost = cCtx.String("server")
			)
			if c == 0 {
				c = uint(runtime.NumCPU())
			}

			// 从服务端获取配置
			var apiRes internal.GetConfigRequest
			resp, err := resty.New().SetTimeout(time.Second * 3).R().SetResult(&apiRes).Get(serverHost)
			if err != nil {
				return err
			}

			if resp.StatusCode() != http.StatusOK {
				return errors.New(fmt.Sprintf("http get status %s", resp.Status()))
			}

			if Node, err = internal.NewNode(serverHost, apiRes, c); err != nil {
				return err
			}

			return nil
		},
		Action: func(cCtx *cli.Context) error {

			Node.Run()
			return nil
		},
	}
)
