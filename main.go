package main

import (
	"flag"
	"fmt"
	"log"
	"log_analyser/database"
	"log_analyser/types"
	"math"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	_ "github.com/joho/godotenv/autoload"
)

func main() {

	configPath := flag.String("config", "config.toml", "path to config.toml file")
	flag.Parse()

	tomlData, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	var conf types.Config
	_, err = toml.Decode(string(tomlData), &conf)

	if err != nil {
		fmt.Println("Error decoding config file: ", err)
		panic(err)
	}

	db := database.NewPostgres(conf.PostgresHost, conf.PostgresUser, conf.PostgresPassword, conf.PostgresDb, conf.PostgresPort)

	for {
		outerStartTime := time.Now()
		repos, err := db.GetRepositories()
		if err != nil {
			log.Fatalln(err)
		}
		for {
			innerStartTime := time.Now()

			var wg sync.WaitGroup
			for _, repo := range repos {
				wg.Add(1)

				// create copy to prevent change during loop
				repoCopy := repo
				go func() {
					defer wg.Done()
					fmt.Println(repoCopy)
					analyseRepository(repoCopy, db, conf)
				}()
			}
			wg.Wait()
			innerEndTime := time.Now()

			outerSleepTime := float64(conf.RepoScanFrequency) - innerEndTime.Sub(outerStartTime).Seconds()
			if outerSleepTime <= 0 {
				break
			} else {
				innerSleepTime := math.Max(float64(conf.RunFrequency), innerEndTime.Sub(innerStartTime).Seconds())
				time.Sleep(time.Duration(innerSleepTime) * time.Second)
			}
		}
	}
}
