package client

import (
	"encoding/json"
	"errors"
	"github.com/black40x/tunl-core/commands"
	"github.com/black40x/tunl-core/tunl"
	"net"
	"os"
	"sync"
	"time"
	"tunl-cli/cmd/options"
	"tunl-cli/cmd/tui"
)

type HttpRequestReceiver func(r *commands.HttpRequest, body []byte, d time.Duration)

type Client struct {
	opt            *options.Options
	conn           *tunl.TunlConn
	receiveRequest HttpRequestReceiver
	bodySync       sync.Map
	mu             sync.Mutex
	version        string
	serverVersion  string
	serverPrivate  bool
}

func NewTunlClient(opt *options.Options, version string) *Client {
	return &Client{
		opt:      opt,
		bodySync: sync.Map{},
		version:  version,
	}
}

func (c *Client) SetHttpRequestReceiver(receiver HttpRequestReceiver) {
	c.receiveRequest = receiver
}

func (c *Client) sendJsonMessage(uuid string, j map[string]interface{}, status int32) {
	data, err := json.Marshal(j)
	if err == nil {
		c.conn.Send(&commands.HttpResponse{
			Uuid:          uuid,
			Header:        []*commands.Header{{Key: "Content-Type", Value: []string{"application/json"}}},
			ContentLength: int64(len(data)),
			Status:        status,
		})
		c.conn.Send(&commands.BodyChunk{
			Uuid: uuid,
			Body: data,
			Eof:  true,
		})
	}
}

func (c *Client) handleCommand(cmd *commands.Transfer) {
	switch cmd.GetCommand().(type) {
	case *commands.Transfer_ServerHeader:
		c.serverVersion = cmd.GetServerHeader().Version
		c.serverPrivate = cmd.GetServerHeader().Private

		if c.serverPrivate && c.opt.ServerPassword == "" {
			tui.PrintError(errors.New("need password for connect to private server"))
			os.Exit(1)
		} else {
			c.conn.Send(&commands.ClientConnect{
				Version:  Version,
				Password: c.opt.ServerPassword,
			})
		}
	case *commands.Transfer_ServerConnect:
		tui.PrintConnectionScreen(
			*c.opt,
			cmd.GetServerConnect().PublicUrl,
			c.version,
			c.serverVersion,
			c.serverPrivate,
			cmd.GetServerConnect().Expire,
		)
	case *commands.Transfer_Error:
		if cmd.GetError().Code == tunl.ErrorSessionExpired {
			tui.PrintInfo(cmd.GetError().GetMessage())
			os.Exit(1)
		} else if cmd.GetError().Code == tunl.ErrorServerFull {
			tui.PrintWarning(cmd.GetError().GetMessage())
			os.Exit(1)
		} else if cmd.GetError().Code == tunl.ErrorUnauthorized {
			tui.PrintError(errors.New(cmd.GetError().GetMessage()))
			os.Exit(1)
		} else {
			tui.PrintError(errors.New(cmd.GetError().GetMessage()))
		}
	case *commands.Transfer_BodyChunk:
		if v, ok := c.bodySync.Load(cmd.GetBodyChunk().Uuid); ok {
			v.(chan *commands.BodyChunk) <- cmd.GetBodyChunk()
		}
	case *commands.Transfer_HttpRequest:
		c.bodySync.Store(cmd.GetHttpRequest().Uuid, make(chan *commands.BodyChunk, 100))
		go c.processRequestCommand(cmd.GetHttpRequest())
	}
}

func (c *Client) Connect() error {
	cn, err := net.Dial("tcp", c.opt.ServerAddr)
	if err != nil {
		return err
	}

	c.conn = tunl.NewTunlConn(cn)

	c.conn.SetOnCommand(func(cmd *commands.Transfer) {
		c.handleCommand(cmd)
	})
	c.conn.SetOnDisconnected(func() {
		tui.PrintError(errors.New("disconnected from server"))
	})
	c.conn.SetOnError(func(err error) {
		tui.PrintError(err)
	})

	go c.conn.HandleConnection()

	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
