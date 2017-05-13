package fixdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	pb "gopkg.in/cheggaaa/pb.v1"
)

var (
	// DefaultDbPath path to db file
	DefaultDbPath = "_state/buildings.fixdb"

	// ErrNotFound error if not found
	ErrNotFound = fmt.Errorf("notfound")

	// Readonly db
	Readonly = true

	db *DB

	bar *pb.ProgressBar
)

// NewContext init db as global variable
func NewContext() error {
	var err error
	db = New(DefaultDbPath)
	if err != nil {
		return err
	}

	// err = db.validate()
	// if err != nil {
	// 	return err
	// }

	return nil
}

// DB wrapper of file
type DB struct {
	*os.File
	cacheID map[int]int64
}

// New return open file and db
func New(path string) *DB {
	flags := os.O_RDWR
	if Readonly {
		flags = os.O_RDONLY
	}
	f, err := os.OpenFile(path, flags, 0777)
	if err != nil {
		panic(err)
	}
	l := &DB{
		File:    f,
		cacheID: make(map[int]int64),
	}
	return l
}

func (db *DB) validate() error {
	_, err := db.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}
	//br := bufio.NewReader(db)
	var last int64
	buf := make([]byte, 24)
	for i := 0; i < db.Len(); i++ {
		_, err := db.Read(buf)
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		var current int64
		binary.Read(bytes.NewReader(buf), binary.BigEndian, &current)
		if last < current {
			last = current
			continue
		} else {
			log.Printf("Error ID value %d > %d at index %d", last, current, i)
			return fmt.Errorf("Error ID value %d > %d at index %d", last, current, i)
		}
	}
	return nil
}

// Get return coordinates by given id
func Get(id int64) ([]float64, error) {
	return db.Get(id)
}

// Get return coordinates by given id
func (db *DB) Get(id int64) ([]float64, error) {
	index := sort.Search(db.Len(), db.makeSearchFunc(int64(id)))
	data := db.get(index)
	br := bytes.NewReader(data)
	var (
		oldID int64
		lat   float64
		lon   float64
	)
	err := binary.Read(br, binary.BigEndian, &oldID)
	if err != nil {
		if err == io.EOF {
			log.Println(index)
			return nil, ErrNotFound
		}
		return nil, err
	}
	if oldID != id {
		return nil, ErrNotFound
	}
	err = binary.Read(br, binary.BigEndian, &lat)
	if err != nil {
		return nil, err
	}
	err = binary.Read(br, binary.BigEndian, &lon)
	if err != nil {
		return nil, err
	}
	return []float64{lat, lon}, nil
}

func (db *DB) makeSearchFunc(needID int64) func(i int) bool {
	return func(i int) bool {
		return db.getID(i) >= needID
	}
}

func (db *DB) Len() int {
	stat, err := db.Stat()
	if err != nil {
		panic(err)
	}
	return int(stat.Size()) / 24
}

func (db *DB) getID(index int) (id int64) {
	// if cachedID, has := db.cacheID[index]; has {
	// 	return cachedID
	// }
	_, err := db.Seek(int64(24*index), os.SEEK_SET)
	if err != nil {
		panic(err)
	}
	err = binary.Read(db.File, binary.BigEndian, &id)
	if err != nil {
		panic(err)
	}
	db.cacheID[index] = id
	return
}

func (db *DB) get(index int) []byte {
	buf := make([]byte, 24)
	_, err := db.ReadAt(buf, int64(24*index))
	if err != nil {
		panic(err)
	}
	return buf
}

func (db *DB) set(i int, data []byte) {
	_, err := db.WriteAt(data, int64(24*i))
	if err != nil {
		panic(err)
	}
}

func (db *DB) Less(i, j int) bool {
	return db.getID(i) < db.getID(j)
}

func (db *DB) Swap(i, j int) {
	delete(db.cacheID, i)
	delete(db.cacheID, j)
	dataI := db.get(i)
	dataJ := db.get(j)
	db.set(i, dataJ)
	db.set(j, dataI)
	bar.Increment()
}
