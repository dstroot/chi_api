// handlers.article_test.go

package main

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	env "github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
	. "github.com/smartystreets/goconvey/convey"
)

// Test that check(e) exits
// this is how you test fails
func TestCheck(t *testing.T) {
	err := errors.New("test check function")

	if os.Getenv("CRASHED") == "1" {
		check(err)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestCheck")
	cmd.Env = append(os.Environ(), "CRASHED=1")
	err = cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestSetupDatabase(t *testing.T) {
	Convey("Bad configuration", t, func() {
		Convey("Bad configuration will not connect", func() {

			// no configuration

			// Connect to database
			err1 := setupDatabase()
			So(err1, ShouldEqual, nil)

			// Ping database
			err2 := db.Ping()
			So(err2, ShouldNotBeNil)
		})
	})
	Convey("Good configuration", t, func() {
		Convey("It can connect to a database", func() {

			// Read configuration
			err := env.Decode(&cfg)
			So(err, ShouldEqual, nil)

			// Connect to database
			err1 := setupDatabase()
			So(err1, ShouldEqual, nil)

			// Ping database
			err2 := db.Ping()
			So(err2, ShouldEqual, nil)
		})
	})

}
