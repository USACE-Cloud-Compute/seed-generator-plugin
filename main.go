package main

import (
	"log"

	"github.com/usace-cloud-compute/cc-go-sdk"
	tiledb "github.com/usace-cloud-compute/cc-go-sdk/tiledb-store"
	_ "github.com/usace-cloud-compute/seed-generator/internal/actions"
)

var version string
var commit string
var date string

func main() {
	//register tiledb
	cc.DataStoreTypeRegistry.Register("TILEDB", tiledb.TileDbEventStore{})

	pm, err := cc.InitPluginManager()
	if err != nil {
		log.Fatalf("Unable to initialize the plugin manager: %s\n", err)
	}

	pm.Logger.Info("Seed Generator", "version", version, "commit", commit, "build-date", date)

	err = pm.RunActions()
	if err != nil {
		pm.Logger.Fatalf("error running actions: %s", err.Error())
	}

}
