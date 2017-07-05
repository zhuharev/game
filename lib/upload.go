package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/zhuharev/game/models"
	"github.com/zhuharev/game/modules/bloblog"
	"github.com/zhuharev/game/modules/middleware"
	"github.com/zhuharev/game/modules/setting"
)

var (
//bufpool = bpool.NewBufferPool(8)

)

func handleUpload(c *middleware.Context) {

	u, err := models.GetUserByToken(c.Query("token"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, err.Error())
		return
	}
	if u == nil {
		log.Println("Ошибка авторизации")
		c.JSON(200, "Ошибка авторизации")
		return
	}

	c.Req.Request.Body = http.MaxBytesReader(c.Resp, c.Req.Request.Body, 10*1024*1024)

	c.Req.ParseMultipartForm(10 * 1024 * 1024)
	file, _, err := c.Req.FormFile("file")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer file.Close()

	bts, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, err.Error())
		return
	}

	id, err := bloblog.Save(bts)
	if err != nil {
		log.Println(err.Error())
		c.JSON(200, err.Error())
		return
	}

	err = models.UsersSetAvatarUrl(u.Id, fmt.Sprintf("https://%s/images/%d", setting.App.App.Host, id))
	if err != nil {
		log.Println(err.Error())
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
