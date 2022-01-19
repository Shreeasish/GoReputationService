package scorer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	radix "github.com/armon/go-radix"
)

// DomainScore
type DomainScore struct {
	Domain string `json:"Domain"`
	Score  string `json:"Score"`
	// TODO: Add status as enums
}

// Document Struct/public member
type DomainScorer struct {
	Count int
	// LastUpdated ..
	LastUpdated time.Time
	rt          *radix.Tree
	mu          sync.Mutex
}

// errStr is a templated string to produce error messages
const errStr = `could not parse line as url,score: %v`

// New initializes an iradix tree from a list of domains from a file
func New(reader io.Reader) (*DomainScorer, error) {
	r := radix.New()

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		l := scanner.Text()
		if strings.HasPrefix(l, `#`) {
			continue
		}

		u, s, err := parseLine(l)
		if err != nil {
			return nil, err
		}
		// Ignore returned value and found bool
		_, _ = r.Insert(u, s)
	}

	log.Printf("Tree initialized with %d domains", r.Len())

	return &DomainScorer{
		Count:       r.Len(),
		LastUpdated: time.Now(),
		rt:          r,
	}, nil
}

// Get the score for a domain that's already been set
// TODO: Switch to LongestPrefix match in order to respond to partial matches
// net/url can be used to match fragments of the incoming url to provide
// more robust checks such hostnames etc.
func (rp DomainScorer) GetDomainScore(d string) *DomainScore {
	v, ok := rp.rt.Get(d)
	// Handle uknown domains with 0

	// Return both domain and score even though domain
	// is known to the caller for ease of use
	if !ok {
		return &DomainScore{
			Domain: d,
			Score:  "0",
		}
	}

	return &DomainScore{
		Domain: d,
		Score:  strconv.FormatInt(v.(int64), 10),
	}
}

// AddDomain updates or adds a domain to the radix tree
func (rp DomainScorer) AddDomain(domain string, score int64) {
	// TODO: Replace with a CAS operation instead of locking
	rp.mu.Lock()
	prev, _ := rp.rt.Insert(domain, score)
	rp.LastUpdated = time.Now()
	rp.Count++
	rp.mu.Unlock()

	log.Printf("Domain %v with Score %v updated at %v with new Score %v", domain, prev, time.Now().String(), score)

}

func parseLine(line string) (string, int64, error) {
	// Format: http://www.baddomain.com,3
	s := strings.Split(line, ",")
	if len(s) != 2 {
		return "", 0, errors.New(fmt.Sprintf(errStr, line))
	}
	i, err := strconv.ParseInt(s[1], 10, 64)
	if err != nil {
		return "", 0, errors.New(fmt.Sprintf(errStr, err))
	}
	uri, err := url.ParseRequestURI(s[0])
	if err != nil {
		return "", 0, errors.New(fmt.Sprintf(errStr, err))
	}
	return uri.Hostname(), i, nil
}
