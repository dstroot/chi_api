// handlers.article_test.go

package main

import (
	"fmt"
	"testing"

	"github.com/dstroot/chi_api/database"
	_ "github.com/joho/godotenv/autoload"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInitialize(t *testing.T) {
	Convey("Initialize", t, func() {
		Convey("Should initialize our configuration", func() {
			fmt.Printf("\n\n")
			err := initialize()
			So(err, ShouldEqual, nil)
		})
	})
}

func TestSetupDatabase(t *testing.T) {
	Convey("Good configuration", t, func() {
		Convey("It can connect to a database", func() {

			// Connect to database
			err1 := setupDatabase()
			So(err1, ShouldEqual, nil)

			// Ping database
			err2 := database.DB.Ping()
			So(err2, ShouldEqual, nil)
		})
	})

}
