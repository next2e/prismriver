package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error
var once sync.Once

func GetDatabase() (*gorm.DB, error) {
	once.Do(func () {
		newDb, openErr := gorm.Open("postgres", "host=localhost port=5432 user=postgres " +
			"password=postgres dbname=postgres sslmode=disable")
		err = openErr
		db = newDb
		if err != nil {
			return
		}
		db.AutoMigrate(&Media{})
	})
	return db, err
}

func AddMedia(media Media) error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}
	db.Create(&media)
	return nil
}

func FindMedia(query string, limit int) []Media {
	db, err := GetDatabase()
	if err != nil {
		logrus.Fatal("Error loading database:", err)
	}
	var media []Media
	db.Limit(limit).Where("title ILIKE ?", "%" + query + "%").Find(&media)
	return media
}

func GetMedia(id string, kind string) (Media, error) {
	db, err := GetDatabase()
	if err != nil {
		logrus.Fatal("Error loading database:", err)
	}
	var media []Media
	db.Where(Media{ID: id, Type: kind}).First(&media)
	if len(media) > 0 {
		return media[0], nil
	}
	return Media{}, errors.New("Media not found in DB.")
}

func GetRandomMedia(limit int) []Media {
	db, err := GetDatabase()
	if err != nil {
		logrus.Fatal("Error loading database:", err)
	}
	var media []Media
	db.Order("random()").Limit(limit).Find(&media)
	return media
}