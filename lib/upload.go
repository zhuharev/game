package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/bloblog"
	"github.com/zhuharev/game/modules/middleware"
)

var (
//bufpool = bpool.NewBufferPool(8)

)

func handleUpload(c *middleware.Context) {

	u, err := models.GetUserByToken(c.Query("token"))
	if err != nil {
		c.JSON(200, err.Error())
		return
	}
	if u == nil {
		c.JSON(200, "Ошибка авторизации")
		return
	}

	c.Req.Request.Body = http.MaxBytesReader(c.Resp, c.Req.Request.Body, 10*1024*1024)

	c.Req.ParseMultipartForm(10 * 1024 * 1024)
	file, _, e := c.Req.FormFile("file")
	if e != nil {
		fmt.Println(e)
		return
	}
	defer file.Close()

	bts, e := ioutil.ReadAll(file)
	if e != nil {
		c.JSON(200, err.Error())
		fmt.Println(e)
		return
	}

	id, err := bloblog.Save(bts)
	if e != nil {
		c.JSON(200, err.Error())
		fmt.Println(e)
		return
	}

	err = models.UsersSetAvatarUrl(u.Id, fmt.Sprintf("/images/%d", id))
	if err != nil {
		c.JSON(200, err.Error())
		return
	}

	c.JSON(200, "ok")
}

func handleAvatar(c *middleware.Context) {
	bts, err := bloblog.Get(c.ParamsInt64(":id"))
	if err != nil {
		log.Println(err)
	}
	c.Resp.Header().Set("Content-Type", "image/jpeg")
	c.Resp.WriteHeader(200)
	c.Resp.Write(bts)
}
