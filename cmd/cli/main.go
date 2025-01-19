package main

import (
	"log"
	"os"

	"github.com/hoitek/Maja-Service/config"
)

func main() {
	// Load configuration
	config.LoadDefault()

	// Args from command line
	args := os.Args

	// Check if there is any argument
	if len(args) < 2 {
		log.Fatal("Missing argument")
	}

	// Check if argument is valid
	switch args[1] {
	case "migrate", "seed":
		// Set gorm instance
		//migrator.SetGorm(database.ConnectGorm())
		//// Set migrations
		//migrator.SetMigrations(
		//	migrations.UsersTable{},
		//)
		//// Set arguments
		//migrator.SetArgs(os.Args)

		// case "start":
		// 	build := path.Join(config.GetRootPath(), "httpMaja")
		// 	cmd := exec.Command(build)
		// 	cmd.Stdout = os.Stdout
		// 	cmd.Stderr = os.Stderr
		// 	if err := cmd.Run(); err != nil {
		// 		log.Fatal(err)
		// 	}
	}
}
