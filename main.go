package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/garyburd/redigo/redis"
	"github.com/mholt/binding"
	"github.com/tidwall/tile38/controller"
	"github.com/zhuharev/go-osm"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

var (
	tilePort           = 7091
	tilePortString     = fmt.Sprint(tilePort)
	tileBuildingsTable = "buildings"
)

type Center struct {
	LongLat string
}

func (c *Center) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&c.LongLat: "center",
	}
}

func main() {
	go func() {
		err := startTileServer()
		if err != nil {
			panic(err)
		}
	}()

	SetDb()

	time.Sleep(1 * time.Second)

	nearby(54.7779274, 32.0219039)
	app := iris.New()
	// output startup banner and error logs on os.Stdout
	app.Adapt(iris.DevLogger())
	// set the router, you can choose gorillamux too
	app.Adapt(httprouter.New())

	api := app.Party("/api/v1")
	api.Get("/buildings", handleBuildings)
	api.Get("/users/me", me)
	api.Get("/auth", handleAuth)
	api.Get("/games/new", handleNewGame)
	api.Get("/games/check", handleCheck)
	api.Get("/location/update", func(ctx *iris.Context) {
		user, err := getUserFromCtx(ctx)
		if err != nil {
			handleError(ctx, err)
			return
		}

		arr := strings.Split(ctx.FormValue("location"), ",")
		if len(arr) != 2 {
			if err != nil {
				handleError(ctx, fmt.Errorf("not 2"))
				return
			}
		}

		lat, err := strconv.ParseFloat(arr[0], 64)
		if err != nil {
			handleError(ctx, err)
			return
		}

		lon, err := strconv.ParseFloat(arr[1], 64)
		if err != nil {
			handleError(ctx, err)
			return
		}

		err = setLocation(user.Id, lat, lon)
		if err != nil {
			handleError(ctx, err)
			return
		}
		ctx.JSON(200, "ok")
	})

	app.Get("/", func(ctx *iris.Context) {
		ctx.JSON(iris.StatusOK, iris.Map{"name": "iris"})
	})

	app.Listen(":7000")
}

func startTileServer() error {
	if err := controller.ListenAndServe("", tilePort, "_tileData", true); err != nil {
		return err
	}
	return nil
}

func getBuildings() (map[int][]float64, error) {

	//res = make(map[int][]float64)

	c, err := redis.Dial("tcp", ":"+tilePortString)
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer c.Close()

	m, e := osm.DecodeFile("map.osm")
	if e != nil {
		panic(e)
	}

	for _, v := range m.Ways {
		for _, t := range v.RTags {
			if t.Key == "building" && t.Value == "yes" {
				longlat := centerFromNodes(v.Nds)
				b := Building{
					Id:   v.ID,
					Long: longlat[1],
					Lat:  longlat[0],
				}
				ret, err := c.Do("SET", tileBuildingsTable, v.StringID(), "POINT", fmt.Sprint(b.Lat), fmt.Sprint(b.Long))
				if err != nil {
					panic(err)
				}
				err = createBuilding(&b)
				if err != nil {
					panic(err)
				}
				fmt.Printf("%s\n", ret)
			}
		}
	}
	return nil, nil
}

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type geoPoint struct {
	Coordinates []float64 `json:"coordinates"`
}

func nearby(lat, long float64) ([]Building, error) {
	//nearby fleet fence point 33.462 -112.268 6000
	c, err := redis.Dial("tcp", ":"+tilePortString)
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer c.Close()
	ret, err := c.Do("NEARBY", tileBuildingsTable, "POINT", fmt.Sprint(lat), fmt.Sprint(long), "200")
	//ret, err := c.Do("GET", tileBuildingsTable, "78411860")
	if err != nil {
		panic(err)
	}

	var buildings []Building

	for _, resp := range ret.([]interface{}) {
		switch resp.(type) {
		case int64:
			fmt.Println(resp)
		case []interface{}:
			for _, item := range resp.([]interface{}) {
				switch item.(type) {
				case []interface{}:
					b := Building{}
					p := geoPoint{}
					for _, iface := range item.([]interface{}) {
						if b.Id == 0 {
							b.Id = com.StrTo(string(iface.([]byte))).MustInt64()
						} else {
							err := json.Unmarshal(iface.([]byte), &p)
							if err != nil {
								panic(err)
							}
							b.Long = p.Coordinates[0]
							b.Lat = p.Coordinates[1]
						}
					}
					buildings = append(buildings, b)
				}
			}
		}
	}
	//fmt.Printf("%v\n", buildings)
	return buildings, nil
}

func testCenter() {
	points := [][]float64{
		{32.0493702, 54.7805089},
		{32.0497840, 54.7805128},
		{32.0493675, 54.7803668},
		{32.0497892, 54.7803764},
		{32.0493680, 54.7803805},
	}
	fmt.Println(center(points))
}

func centerFromNodes(nodes []osm.Point) []float64 {
	var points [][]float64
	for _, v := range nodes {
		points = append(points, []float64{v.Lat, v.Lng})
	}
	return center(points)
}

//  возвращает длину отрезка с координатами (x1,y1)-(x2,y2)
func length(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}

func center(points [][]float64) []float64 {
	var (
		xc float64
		yc float64
		P  float64
		n  = len(points)
	)

	for i, xy := range points {
		l := length(xy[0], xy[1], points[(i+1)%n][0], points[(i+1)%n][1])
		xc += l * (xy[0] + points[(i+1)%n][0]) / 2
		yc += l * (xy[1] + points[(i+1)%n][1]) / 2
		P += l
	}
	xc /= P
	yc /= P
	return []float64{xc, yc}
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	var earthRadiusKm = 6371.0

	var dLat = degreesToRadians(lat2 - lat1)
	var dLon = degreesToRadians(lon2 - lon1)

	lat1 = degreesToRadians(lat1)
	lat2 = degreesToRadians(lat2)

	var a = math.Sin(dLat/2.0)*math.Sin(dLat/2.0) +
		math.Sin(dLon/2.0)*math.Sin(dLon/2.0)*math.Cos(lat1)*math.Cos(lat2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1.0-a))
	return earthRadiusKm * c
}
