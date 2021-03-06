package tile38

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/Unknwon/com"
	"github.com/garyburd/redigo/redis"
	"github.com/tidwall/tile38/controller"
	"github.com/zhuharev/go-osm"
)

var (
	tilePort           = 7091
	tilePortString     = fmt.Sprint(tilePort)
	tileBuildingsTable = "buildings"
)

func StartTileServer() error {
	if err := controller.ListenAndServe("", tilePort, "_tileData", true); err != nil {
		return err
	}
	return nil
}
func GetBuildings() (map[int64][]float64, error) {

	res := make(map[int64][]float64)

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
				_, err := c.Do("SET", tileBuildingsTable, v.StringID(), "POINT", fmt.Sprint(longlat[0]), fmt.Sprint(longlat[1]))
				if err != nil {
					panic(err)
				}
				res[v.ID] = []float64{longlat[0], longlat[1]}
			}
		}
	}
	return res, nil
}

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type geoPoint struct {
	Coordinates []float64 `json:"coordinates"`
}

func Nearby(lat, long float64) (map[int64][]float64, error) {
	//nearby fleet fence point 33.462 -112.268 6000
	c, err := redis.Dial("tcp", ":"+tilePortString)
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer c.Close()
	ret, err := c.Do("NEARBY", tileBuildingsTable, "POINT", fmt.Sprint(lat), fmt.Sprint(long), "250")
	//ret, err := c.Do("GET", tileBuildingsTable, "78411860")
	if err != nil {
		panic(err)
	}

	var buildings = make(map[int64][]float64)

	for _, resp := range ret.([]interface{}) {
		switch resp.(type) {
		case int64:
			fmt.Println(resp)
		case []interface{}:
			for _, item := range resp.([]interface{}) {
				switch item.(type) {
				case []interface{}:
					p := geoPoint{}
					var id int64
					for _, iface := range item.([]interface{}) {
						if id == 0 {
							id = com.StrTo(string(iface.([]byte))).MustInt64()
						} else {
							err := json.Unmarshal(iface.([]byte), &p)
							if err != nil {
								panic(err)
							}
							buildings[id] = []float64{p.Coordinates[1], p.Coordinates[0]}
						}
					}
				}
			}
		}
	}
	//fmt.Printf("%v\n", buildings)
	return buildings, nil
}

func centerFromNodes(nodes []osm.Point) []float64 {
	var points [][]float64
	for _, v := range nodes {
		points = append(points, []float64{v.Lat, v.Lng})
	}
	return Center(points)
}

//  возвращает длину отрезка с координатами (x1,y1)-(x2,y2)
func length(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((x1-x2)*(x1-x2) + (y1-y2)*(y1-y2))
}

func Center(points [][]float64) []float64 {
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
