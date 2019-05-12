package weatherstation

import (
	"context"
	eh "github.com/looplab/eventhorizon"
	"github.com/looplab/eventhorizon/eventhandler/projector"
	"github.com/looplab/eventhorizon/examples/weatherstation/domain"
	"github.com/looplab/eventhorizon/repo/mongodb"
	"github.com/looplab/eventhorizon/repo/version"
	"log"
	"os"
	"time"
)

type Projector struct {
	repo                  *mongodb.Repo
	projectorEventHandler *projector.EventHandler
}

func SetupProjector() *Projector {
	// Local Mongo testing with Docker
	url := os.Getenv("MONGO_HOST")
	if url == "" {
		// Default to localhost
		//url = "localhost:27017"
		url = RepoHost + ":" + RepoPort
	}

	var err error

	var repo *mongodb.Repo
	// Create the read repositories.
	repo, err = mongodb.NewRepo(url, DbPrefix, "projector")
	if err != nil {
		log.Fatalf("could not create projector repository: %s", err)
	}
	repo.SetEntityFactory(func() eh.Entity { return &domain.Temperature{} })

	// A version repo is needed for the projector to handle eventual consistency.
	var versionRepo *version.Repo
	versionRepo = version.NewRepo(repo)

	// Create and register a read model for individual temperature reads.
	projectorEventHandler := projector.NewEventHandler(domain.NewTemperatureReadProjector(), versionRepo)
	projectorEventHandler.SetEntityFactory(func() eh.Entity { return &domain.Temperature{} })

	return &Projector{
	   repo,
	   projectorEventHandler,
	}
}

func (p *Projector) ClearRepo(ctx context.Context) {
	err := p.repo.Clear(ctx)
	if err != nil {
		log.Fatalf("could not clear repo: %s", err)
	}
}

func (p *Projector) ListTemperatureHistory(ctx context.Context) (string, float32, []float32) {
	var resultName string
	var resultTemp float32
	var resultHistory []float32
	var weatherStations []eh.Entity
	var err error
	var i = 0
	for {
		// Projector is updated asynchronously so we need to wait an undefined time to see results,
		// in this case 100ms, 10 retries.
		time.Sleep(100 * time.Millisecond)
		weatherStations, err = p.repo.FindAll(ctx)
		if err != nil {
			log.Println("error:", err)
		}
		i++
		if len(weatherStations) > 0 {
			break
		}
		if i > 10 {
			log.Printf("No weatherStations found. Return.")
			return "", 0, resultHistory // empty result, error is captured by test
		}
	}
	for _, weatherStation := range weatherStations {
		if station, ok := weatherStation.(*domain.Temperature); ok {
			log.Printf("Name %s Temperature %f", station.Name, station.Temperature)

			resultName = station.Name
			resultTemp = station.Temperature
			resultHistory = station.History

			for _, hist := range station.History {
				log.Printf("Hist: %f", hist)
			}
		}
	}

	return resultName, resultTemp, resultHistory
}
