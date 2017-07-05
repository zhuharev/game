package migrations

import (
	"log"

	"github.com/go-xorm/xorm"
	pb "gopkg.in/cheggaaa/pb.v1"

	"github.com/zhuharev/game/modules/fixdb"
)

type Building struct {
	Id   int64
	Area int
}

func updateArea(x *xorm.Engine) error {

	var buildings []Building
	err := x.Find(&buildings)
	if err != nil {
		log.Println(err)
	}
	bar := pb.StartNew(len(buildings))
	for _, v := range buildings {
		_, area, err := fixdb.Get(v.Id)
		if err != nil {
			log.Printf("Err get building with id %d", v.Id)
			continue
		}
		_, err = x.Exec("update building set area = ? where id = ?", area, v.Id)
		if err != nil {
			log.Println(err)
		}
		bar.Increment()
	}
	return nil
}
