package client

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/black40x/tunl-core/commands"
	"github.com/black40x/tunl-core/tunl"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"tunl-cli/cmd/tui"
)

func (c *Client) httpRequest(r *commands.HttpRequest) (res *http.Response, err error) {
	ts := time.Now()
	client := &http.Client{}
	client.Timeout = c.opt.HttpTimeout
	buf := bytes.NewBuffer([]byte{})

	ch, ok := c.bodySync.Load(r.Uuid)
	if r.ContentLength > 0 && ok {
		var bodySize int64 = 0

	L:
		for {
			select {
			case body := <-ch.(chan *commands.BodyChunk):
				bodySize += int64(len(body.Body))
				buf.Write(body.Body)
				if bodySize >= r.ContentLength {
					break L
				}
			case <-time.After(time.Second * 60):
				return nil, errors.New("request timeout")
			}
		}
	}

	req, err := http.NewRequest(r.Method, c.opt.LocalAddr.ToProtoString()+r.Uri, buf)

	if err != nil {
		return nil, err
	}

	for _, ck := range r.Cookies {
		req.AddCookie(&http.Cookie{
			Name:     ck.Name,
			Value:    ck.Value,
			Path:     ck.Path,
			Domain:   ck.Domain,
			Expires:  time.UnixMicro(ck.Expires),
			HttpOnly: ck.HttpOnly,
			Secure:   ck.Secure,
		})
	}

	for _, h := range r.Header {
		for _, v := range h.GetValue() {
			req.Header.Add(h.GetKey(), v)
		}
	}

	for k, v := range c.opt.RequestHeaders {
		req.Header.Add(k, v)
	}

	duration := time.Since(ts)
	tui.PrintURL(r.Method, r.Uri, duration)
	if c.receiveRequest != nil {
		c.receiveRequest(r, buf.Bytes(), duration)
	}

	return client.Do(req)
}

func (c *Client) processWeb(r *commands.HttpRequest) {
	res, err := c.httpRequest(r)
	if err != nil {
		c.conn.Send(&commands.HttpResponse{
			Uuid:          r.Uuid,
			ContentLength: 0,
			ErrorCode:     int64(tunl.ErrorClientResponse),
		})
	} else {
		mes := commands.HttpResponse{
			Uuid:          r.Uuid,
			ContentLength: res.ContentLength,
			Proto:         res.Proto,
			Status:        int32(res.StatusCode),
			ErrorCode:     -1,
		}

		for k, v := range res.Header {
			mes.Header = append(mes.Header, &commands.Header{Key: k, Value: v})
		}

		for k, v := range c.opt.ResponseHeaders {
			mes.Header = append(mes.Header, &commands.Header{Key: k, Value: []string{v}})
		}

		_, err := c.conn.Send(&mes)
		if err != nil {
			c.conn.Send(&commands.HttpResponse{
				Uuid:          r.Uuid,
				ContentLength: 0,
				ErrorCode:     int64(tunl.ErrorClientResponse),
			})
		} else {
			re := bufio.NewReader(res.Body)
			buf := make([]byte, 0, tunl.ReaderSize)
			for {
				n, _ := re.Read(buf[:cap(buf)])
				buf = buf[:n]
				if n == 0 {
					if res.ContentLength == -1 {
						c.conn.Send(&commands.BodyChunk{
							Uuid: r.Uuid,
							Body: nil,
							Eof:  true,
						})
					}
					break
				} else {
					c.conn.Send(&commands.BodyChunk{
						Uuid: r.Uuid,
						Body: buf,
						Eof:  false,
					})
				}
			}
		}
	}
}

func (c *Client) getFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}
	contentType := http.DetectContentType(buffer)
	return contentType, nil
}

func (c *Client) processDir(r *commands.HttpRequest) {
	tui.PrintURL(r.Method, r.Uri, 0)
	if c.receiveRequest != nil {
		c.receiveRequest(r, nil, 0)
	}

	uri := r.Uri
	if uri == "/" {
		uri = uri + "index.html"
	}

	p := path.Join(c.opt.LocalAddr.ToString(), uri)
	st, err := os.Stat(p)
	if err != nil || st.IsDir() {
		c.sendJsonMessage(
			r.Uuid,
			map[string]interface{}{
				"error":   true,
				"message": "File not exists",
			},
			http.StatusNotFound,
		)

		return
	}

	f, err := os.Open(p)
	if err != nil {
		c.sendJsonMessage(
			r.Uuid,
			map[string]interface{}{
				"error":   true,
				"message": "File not exists",
			},
			http.StatusNotFound,
		)

		return
	}
	defer f.Close()
	ct, _ := c.getFileContentType(f)
	f.Seek(0, 0)

	if strings.ToLower(filepath.Ext(p)) == ".html" || strings.ToLower(filepath.Ext(p)) == ".htm" {
		ct = strings.ReplaceAll(ct, "text/plain", "text/html")
	}

	mes := commands.HttpResponse{
		Uuid:          r.Uuid,
		ContentLength: st.Size(),
		Proto:         "HTTP/1.1",
		Status:        int32(http.StatusOK),
		ErrorCode:     -1,
	}

	mes.Header = []*commands.Header{
		{Key: "Content-Type", Value: []string{ct}},
	}

	for k, v := range c.opt.ResponseHeaders {
		mes.Header = append(mes.Header, &commands.Header{Key: k, Value: []string{v}})
	}

	_, err = c.conn.Send(&mes)
	if err != nil {
		c.sendJsonMessage(
			r.Uuid,
			map[string]interface{}{
				"error":   true,
				"message": "Client response error",
			},
			http.StatusBadRequest,
		)
	} else {
		re := bufio.NewReader(f)
		buf := make([]byte, 0, tunl.ReaderSize)

		for {
			n, _ := re.Read(buf[:cap(buf)])
			buf = buf[:n]
			if n == 0 {
				break
			} else {
				c.conn.Send(&commands.BodyChunk{
					Uuid: r.Uuid,
					Body: buf,
				})
			}
		}
	}
}

func (c *Client) processRequestCommand(r *commands.HttpRequest) {
	ch, ok := c.bodySync.Load(r.Uuid)
	if ok {
		defer func(ch chan *commands.BodyChunk) {
			c.bodySync.Delete(r.Uuid)
			close(ch)
		}(ch.(chan *commands.BodyChunk))
	}

	if c.opt.BasicAuth != nil {
		if !r.CheckBasicAuth(c.opt.BasicAuth.Login, c.opt.BasicAuth.Pass) {
			// ToDo - move to Connection messages
			c.conn.Send(&commands.HttpResponse{
				Uuid: r.Uuid,
				Header: []*commands.Header{
					{Key: "WWW-Authenticate", Value: []string{`Basic realm="restricted", charset="UTF-8"`}},
				},
				ContentLength: 0,
				Status:        http.StatusUnauthorized,
				ErrorCode:     -1,
			})
			return
		}
	}

	switch c.opt.LocalAddr.Type() {
	case tunl.DIR:
		c.processDir(r)
	case tunl.PORT, tunl.IP:
		c.processWeb(r)
	case tunl.URL:
		c.sendJsonMessage(
			r.Uuid,
			map[string]interface{}{
				"error":   true,
				"message": "Can't process URL",
			},
			http.StatusBadRequest,
		)
	}
}
