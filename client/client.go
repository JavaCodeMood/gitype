// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package client 负责客户端的渲染
package client

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/caixw/typing/data"
	"github.com/caixw/typing/feed"
	"github.com/caixw/typing/vars"
	"github.com/issue9/mux"
)

// Client 展示给用户的前端页面。
type Client struct {
	mux     *mux.Mux
	path    *vars.Path
	updated int64      // 更新时间，一般为重新加载数据的时间
	etag    string     // 所有页面都采用相同的 etag
	data    *data.Data // 加载的数据，每次加载都会被重置
	tpl     *template.Template
	routes  []string // 记录路由项，释放时，需要删除这些路由项
}

// New 声明一个新的 Client 实例
func New(path *vars.Path, mux *mux.Mux) (*Client, error) {
	d, err := data.Load(path)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	c := &Client{
		mux:     mux,
		path:    path,
		data:    d,
		updated: now,
		etag:    strconv.FormatInt(now, 10),
		routes:  make([]string, 0, 10),
	}

	if err = c.initTemplate(); err != nil {
		return nil, err
	}

	c.initRoutes()

	if err = c.initFeeds(); err != nil {
		return nil, err
	}

	return c, nil
}

// Free 释放所有的数据
func (c *Client) Free() {
	c.removeFeeds()
	c.removeRoutes()

	c.tpl = nil
	c.data = nil
}

func (c *Client) initFeeds() error {
	conf := c.data.Config

	if conf.RSS != nil {
		rss, err := feed.BuildRSS(c.data)
		if err != nil {
			return err
		}

		c.mux.GetFunc(conf.RSS.URL, c.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(rss.Bytes())
		}))
	}

	if conf.Atom != nil {
		atom, err := feed.BuildAtom(c.data)
		if err != nil {
			return err
		}

		c.mux.GetFunc(conf.Atom.URL, c.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(atom.Bytes())
		}))
	}

	if conf.Sitemap != nil {
		sitemap, err := feed.BuildSitemap(c.data)
		if err != nil {
			return err
		}

		c.mux.GetFunc(conf.Sitemap.URL, c.prepare(func(w http.ResponseWriter, r *http.Request) {
			w.Write(sitemap.Bytes())
		}))
	}

	return nil
}

func (c *Client) removeFeeds() {
	conf := c.data.Config

	if conf.RSS != nil {
		c.mux.Remove(conf.RSS.URL)
	}

	if conf.Atom != nil {
		c.mux.Remove(conf.Atom.URL)
	}

	if conf.Sitemap != nil {
		c.mux.Remove(conf.Sitemap.URL)
	}
}
