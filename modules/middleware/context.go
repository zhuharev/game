package middleware

import (
	"github.com/monoculum/formam"
	"gopkg.in/macaron.v1"
)

type Context struct {
	*macaron.Context
}

func Contexter() macaron.Handler {
	return func(c *macaron.Context) {
		ctx := &Context{
			Context: c,
		}

		c.Map(ctx)
	}
}

func (c *Context) ReadForm(obj interface{}) error {
	if e := c.Context.Req.ParseForm(); e != nil {
		return e
	}
	dec := formam.NewDecoder(&formam.DecoderOptions{
		TagName: "form",
	})
	vals := c.Context.Req.Form
	return dec.Decode(vals, obj)
}
