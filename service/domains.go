package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	iradix "github.com/hashicorp/go-immutable-radix"
)

type ReputationProvider struct {
	Count       int
	LastUpdated string
	rt          *iradix.Tree
}

// Initialize an iradix tree from a list of domains from a file
func InitializeFromPath(path string) *ReputationProvider {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Init tree
	r := iradix.New()
	var c int = 0

	// Read line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		l := scanner.Text()
		// Format: www.baddomain.com,<score: int>
		s := strings.Split(l, ",")
		// Insert into tree.
		r, _, _ = r.Insert([]byte(s[0]), s[1])
		c++
	}

	log.Printf("Inserted %d domains", c)

	return &ReputationProvider{
		Count:       c,
		LastUpdated: time.Now().String(),
		rt:          r,
	}
}

type DomainScore struct {
	Domain string `json:"Domain"`
	Score  string `json:"Score"`
}

// Get the score for a domain that's already been set
func (rp ReputationProvider) GetDomainScore(d string) (*DomainScore, bool) {
	v, ok := rp.rt.Get([]byte(d))
	if !ok {
		return nil, false
	}

	// Return both domain and score even though domain
	// is known to the caller for ease of use
	return &DomainScore{
		Domain: d,
		Score:  v.(string),
	}, ok
}

/**
 * Update or add a new domain and score to an existing trie.
 *
 * Updates do not need any synchronization here since we are using an immutable
 * Data Structure. Any queries made while nodes are being added will be sound
 * since no changes will be made to the original
 */
func (rp ReputationProvider) AddDomain(ds *DomainScore) {
	rp.rt, _, _ = rp.rt.Insert([]byte(ds.Domain), ds.Score)
}
