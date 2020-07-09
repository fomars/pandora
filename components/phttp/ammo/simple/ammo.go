// Copyright (c) 2017 Yandex LLC. All rights reserved.
// Use of this source code is governed by a MPL 2.0
// license that can be found in the LICENSE file.
// Author: Vladimir Skipor <skipor@yandex-team.ru>

package simple

import (
	"net/http"

	"github.com/yandex/pandora/components/phttp"
	"github.com/yandex/pandora/core/aggregator/netsample"
)


type BaseAmmo struct {
	tag       string
	id        int
	isInvalid bool
}

func (a *BaseAmmo) SetId(id int) {
	a.id = id
}

func (a *BaseAmmo) Id() int {
	return a.id
}

func (a *BaseAmmo) Invalidate() {
	a.isInvalid = true
}

func (a *BaseAmmo) IsInvalid() bool {
	return a.isInvalid
}

func (a *BaseAmmo) IsValid() bool {
	return !a.isInvalid
}


type Ammo struct {
	data	  string
	BaseAmmo
}

func (a *Ammo) Reset(data string, tag string) {
	*a = Ammo{data: data, BaseAmmo: BaseAmmo{tag, -1, false}}
}


type HttpAmmo struct {
	// OPTIMIZE(skipor): reuse *http.Request.
	// Need to research is it possible. http.Transport can hold reference to http.Request.
	req       *http.Request
	BaseAmmo
}

func (a *HttpAmmo) Request() (*http.Request, *netsample.Sample) {
	sample := netsample.Acquire(a.tag)
	sample.SetId(a.id)
	return a.req, sample
}

func (a *HttpAmmo) Reset(req *http.Request, tag string) {
	*a = HttpAmmo{req: req, BaseAmmo: BaseAmmo{tag, -1, false}}
}


var _ phttp.Ammo = (*HttpAmmo)(nil)
