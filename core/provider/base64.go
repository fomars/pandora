package provider

import (
	"bufio"
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/yandex/pandora/components/phttp/ammo/simple"
	"github.com/yandex/pandora/components/phttp/ammo/simple/raw"
	"github.com/yandex/pandora/core"
	"io"
	"net/http"
	"strings"
	"sync"
)

type Provider struct {
	simple.Provider
	Config
}

type Config struct {
	File string `validate:"required"`
	// Limit limits total num of ammo. Unlimited if zero.
	Limit int `validate:"min=0"`
	// Redefine HTTP headers
	Headers []string
	ContinueOnError bool
}

func Base64Provider(fs afero.Fs, conf Config) *Provider {
	var p Provider
	p = Provider{
		Provider: simple.NewProvider(fs, conf.File, p.start),
		Config:   conf,
	}
	return &p
}

func (p *Provider) start(ctx context.Context, ammoFile afero.File) error {
	var ammoNum, line int = 0, 0
	decodedConfigHeaders, err := raw.DecodeHTTPConfigHeaders(p.Config.Headers)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(ammoFile)

	for scanner.Scan() {
		line++
		if p.Limit != 0 && ammoNum >= p.Limit{
			break
		}
		base64_data := scanner.Text()
		request, err := ToRequest(base64_data, decodedConfigHeaders)
		ammo := p.Pool.Get().(*simple.Ammo)
		ammo.Reset(request, "")  // TODO: implement tags
		if err != nil {
			if p.Config.ContinueOnError == true {
				ammo.Invalidate()
			} else {
				return errors.Wrapf(err, "failed to decode ammo at line: %v; data: %q", line, base64_data)
			}
		}
		select {
		case p.Sink <- ammo:
			ammoNum++
		case <-ctx.Done():
			return nil
		}
	}
	ammoFile.Seek(0, 0)
	return nil
}

func ToRequest(data string,headers []raw.Header) (req *http.Request, err error) {
	uri := "http://" + data.Host + data.Uri
	req, err = http.NewRequest(data.Method, uri, strings.NewReader(d.Body))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for k, v := range data.Headers {
		req.Header.Set(k, v)
	}
	return
}