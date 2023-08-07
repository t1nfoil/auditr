package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "github.com/lib/pq"
	"github.com/microsoftgraph/msgraph-sdk-go/models/odataerrors"
)

var config auditrConfig
var daLock sync.Mutex
var siLock sync.Mutex
var auditrDb *gorm.DB

func main() {
	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.loadConfigYaml()
	config.checkConfig()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Printf("captured %v, stopping profiler and exiting..\n", sig)
			cancel()
			os.Exit(0)
		}
	}()

	auditrDb, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  config.GetDatabaseDSN(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	go auditLoop(ctx)

	fmt.Println("auditr running..")
	fmt.Printf("trim entry age time: %d seconds\n", config.GetTrimEntryAge())
	fmt.Printf("update interval: %d seconds\n", config.GetUpdateInterval())

	for {
		time.Sleep(time.Duration(config.GetUpdateInterval()) * time.Second)
	}

}

func printOdataError(err error) {
	switch err.(type) {
	case *odataerrors.ODataError:
		typed := err.(*odataerrors.ODataError)
		fmt.Printf("error: %v", typed.Error())
		if terr := typed.GetErrorEscaped(); terr != nil {
			fmt.Printf("code: %s", *terr.GetCode())
			fmt.Printf("msg: %s", *terr.GetMessage())
		}
	default:
		fmt.Printf("%T > error: %#v", err, err)
	}
}
