package main

import (
	"flag"
	"log"
	"os"

	"github.com/Shreeasish/reputation/handler"
	"github.com/Shreeasish/reputation/router"
	"github.com/Shreeasish/reputation/scorer"
)

var (
	domainList = flag.String("domain_list", "/app/resources/domains.csv", "Path to csv of domains to load on startup")
)

func main() {
	f, err := os.Open(*domainList)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer f.Close()

	scorer, err := scorer.New(f)
	if err != nil {
		log.Fatalf("Unable to initialize scorer: %v", err)
	}
	api := handler.New(scorer)
	router.NewRouter([]router.Handler{
		{"/health", "GET", api.HandleHealthCheck},
		{"/urlinfo/url/{url}", "GET", api.GetScoreHandler},
		{"/urlinfo/update/score/{score}/url/{url}", "GET", api.UpdateDomainHandler},
	})
}
