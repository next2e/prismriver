package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/ttpcodes/prismriver/internal/app/constants"
	"sync"

	// Import Postgres dialect.
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error
var once sync.Once

// BeQuiet is built-in media used for the "Be Quiet!" feature.
var BeQuiet = &Media{
	ID:     "bequiet",
	Length: 3710000000,
	Title:  "Please Be Quiet!",
	Type:   "internal",
}

// GetDatabase gets the instance of the database connection used for the application.
func GetDatabase() (*gorm.DB, error) {
	once.Do(func() {
		dbHost := viper.GetString(constants.DBHOST)
		dbName := viper.GetString(constants.DBNAME)
		dbPassword := viper.GetString(constants.DBPASSWORD)
		dbPort := viper.GetString(constants.DBPORT)
		dbUser := viper.GetString(constants.DBUSER)
		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost,
			dbPort, dbUser, dbPassword, dbName)
		newDb, openErr := gorm.Open("postgres", connString)
		err = openErr
		db = newDb
		if err != nil {
			return
		}
		db.AutoMigrate(&Media{})
		db.FirstOrCreate(BeQuiet)
	})
	return db, err
}

// AddMedia adds a new Media to the database.
func AddMedia(media Media) error {
	db, err := GetDatabase()
	if err != nil {
		return err
	}
	db.Create(&media)
	return nil
}

// FindMedia searches the database for Media items matching the title in query and returns the number of results
// specified by limit.
func FindMedia(query string, limit int) []Media {
	db, err := GetDatabase()
	if err != nil {
		logrus.Fatal("Error loading database:", err)
	}
	var media []Media
	db.Limit(limit).Where("title ILIKE ? AND type <> ?", "%"+query+"%", "internal").Find(&media)
	return media
}

// GetMedia attempts to return the Media identified by id and kind, and returns an error if not found.
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
	return Media{}, errors.New("media not found in DB")
}

// GetRandomMedia returns a number of random Media specified by limit.
func GetRandomMedia(limit int) []Media {
	db, err := GetDatabase()
	if err != nil {
		logrus.Fatal("Error loading database:", err)
	}
	var media []Media
	db.Order("random()").Where("type <> ?", "internal").Limit(limit).Find(&media)
	return media
}
