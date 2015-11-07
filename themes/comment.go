// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package themes

import (
	"strconv"

	"github.com/caixw/typing/core"
)

type Comment struct {
	PostID   int64
	PostName string

	ID          int64
	Created     int64
	IP          string
	Agent       string
	Content     string
	IsAdmin     bool
	AuthorName  string
	AuthorURL   string
	AuthorEmail string
}

func (c *Comment) Permalink() string {
	if len(c.PostName) > 0 {
		return core.CommentURL(opt, c.PostName, c.ID)
	}
	return core.CommentURL(opt, strconv.FormatInt(c.PostID, 10), c.ID)
}
