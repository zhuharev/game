package nearbydb

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/boltdb/bolt"
	"github.com/tidwall/tile38/geojson/geo"
	"github.com/tidwall/tile38/geojson/geohash"
	"github.com/zhuharev/boltutils"
)

var (
	// DefaultDistance represent radius to get nearby in metres
	DefaultDistance = 500.0

	db *boltutils.DB

	dbPath     = "_state/nearbydb.bolt"
	bucketName = []byte("alo")
)

// NewContext open boltdb file with points
func NewContext() error {
	var err error
	db, err = boltutils.Open(dbPath, 0777, &bolt.Options{ReadOnly: true})
	if err != nil {
		return err
	}
	return nil
}

// Nearby return nearby points by given coordinates
func Nearby(lat, lon float64) (map[int64][]float64, error) {
	points, err := nearby(lat, lon, DefaultDistance)
	if err != nil {
		return nil, err
	}
	m := make(map[int64][]float64)
	for _, p := range points {
		m[p.id] = p.Coordinates
	}
	return m, nil
}

func nearby(lat, lon float64, distance float64) (res []*point, err error) {
	id, err := geohash.Encode(lat, lon, 5)
	if err != nil {
		return
	}

	t := 4
	prefx := []byte(id)[:t]
	db.IteratePrefix(bucketName, prefx, func(k []byte, v []byte) error {
		rdr := bytes.NewReader(v)
		gzReader, err := gzip.NewReader(rdr)
		if err != nil {
			return err
		}
		err = iterReader(gzReader, func(p *point) error {
			if distance >= geo.DistanceTo(lat, lon, p.Lat(), p.Lon()) {
				res = append(res, p)
			}
			return nil
		})
		return err
	})
	return
}

func iterReader(r io.Reader, cb func(p *point) error) error {
	br := r
	for {
		var (
			id  int64
			lat float64
			lon float64
		)
		err := binary.Read(br, binary.BigEndian, &id)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		err = binary.Read(br, binary.BigEndian, &lat)
		if err != nil {
			return err
		}
		err = binary.Read(br, binary.BigEndian, &lon)
		if err != nil {
			return err
		}

		p := newPoint(id, lat, lon)
		if cb != nil {
			err = cb(p)
			if err != nil {
				return err
			}
		}

	}
	return nil
}

func newPoint(id int64, lat, lon float64) *point {
	return &point{id: id, Coordinates: []float64{lat, lon}}
}

type point struct {
	id          int64
	Coordinates []float64 `json:"coordinates"`
}

func (p point) ID() string {
	return fmt.Sprint(p.id)
}

func (p point) Lat() float64 {
	return p.Coordinates[0]
}

func (p point) Lon() float64 {
	return p.Coordinates[1]
}

func (p point) LatString() string {
	return fmt.Sprint(p.Coordinates[0])
}

func (p point) LonString() string {
	return fmt.Sprint(p.Coordinates[1])
}

func (p point) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, p.id)
	if err != nil {
		return 0, err
	}

	err = binary.Write(w, binary.BigEndian, p.Coordinates[0])
	if err != nil {
		return 8, err
	}
	err = binary.Write(w, binary.BigEndian, p.Coordinates[1])
	if err != nil {
		return 16, err
	}
	return 24, nil
}

type points []point

func (p points) Len() int { return len(p) }

func (p points) Less(i, j int) bool { return p[i].id < p[i].id }

func (p points) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
