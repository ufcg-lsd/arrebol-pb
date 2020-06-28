package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

type Storage struct {
	driver *gorm.DB
}

const dbDialect string =  "postgres"
var DB *Storage

func New(host string, port string, user string, dbname string, password string) *Storage {
	dbConfig := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, dbname, password)
	driver, err := gorm.Open(dbDialect, dbConfig)

	if err != nil {
		log.Fatalln(err.Error())
	}

	err = driver.DB().Ping()

	if err != nil {
		log.Fatalln(err.Error())
	}

	DB = &Storage{
		driver,
	}

	return DB
}

func (s *Storage) Setup() {
	s.CreateSchema()
	createDefaults(s)
}

func (s *Storage) Driver() *gorm.DB {
	return s.driver
}

func createDefaults(storage *Storage) {
	q := &Queue{
		Name: "Default",
	}

	queue, err := storage.GetDefaultQueue()

	if queue == nil {
		err = storage.SaveQueue(q)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Println("Default queue already exists")
	}
}
