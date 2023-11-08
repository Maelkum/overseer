package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/labstack/echo/v4"
	"github.com/spf13/pflag"
)

// Super dumb HTTP server, something for the `observer` to start and control.
func main() {

	var (
		flagAddress string
		flagName    string
		flagExec    bool
	)

	pflag.StringVarP(&flagAddress, "address", "a", ":8080", "address for the server to use")
	pflag.StringVarP(&flagName, "name", "n", "example-server", "name for the server to use")
	pflag.BoolVarP(&flagExec, "exec", "e", false, "execute an external command or not")

	pflag.Parse()

	srv := echo.New()
	srv.HidePort = true
	srv.HideBanner = true

	srv.GET("/version", versionHandler)
	srv.GET("/name", nameHandler(flagName))

	if flagExec {
		log.Printf("executing external command")

		out, err := exec.Command("ls", "-lat").Output()
		if err != nil {
			log.Printf("executing external command failed: %s", err)
		} else {
			fmt.Printf("listing current directory: \n======================\n%s\n======================\n", string(out))
		}
	}

	workdir, _ := os.Getwd()
	fmt.Printf("workdir: %s\n", workdir)

	err := srv.Start(flagAddress)
	if err != nil {
		log.Fatalf("could not start server: %s", err)
	}

}

func versionHandler(ctx echo.Context) error {
	log.Printf("received version request")
	return ctx.String(http.StatusOK, "1.0.0")
}

func nameHandler(name string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		log.Printf("received name request")
		return ctx.String(http.StatusOK, name)
	}
}
