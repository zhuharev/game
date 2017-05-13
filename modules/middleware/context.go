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

// JSON put json of data in response
// func (c *Context) JSON(status int, v interface{}) {
// 	var (
// 		result []byte
// 		err    error
// 	)
//
// 	if c.QueryBool("pretty") {
// 		result, err = json.MarshalIndent(v, "", "  ")
// 	} else {
// 		result, err = json.Marshal(v)
// 	}
// 	if err != nil {
// 		http.Error(c.Resp, err.Error(), 500)
// 		return
// 	}
//
// 	// json rendered fine, write out the result
// 	c.Resp.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	c.Resp.WriteHeader(status)
// 	c.Resp.Write(result)
// }
