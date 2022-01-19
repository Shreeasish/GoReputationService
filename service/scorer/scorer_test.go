package scorer

import (
	"strconv"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		lines string
		count int
	}{
		{
			`# Checking comment
http://www.google.com,0`,
			1,
		},
		{
			`https://www.java.com,10
https://www.log4j.com,10`,
			2,
		},
	}

	for _, test := range tests {
		s := strings.NewReader(test.lines)
		scorer, err := New(s)
		if err != nil {
			t.Fatalf("Got err:%v", err)
		}
		if scorer.Count != test.count {
			t.Fatalf("Want %v Got %v", test.count, scorer.Count)
		}
	}
}

func TestGetDomain(t *testing.T) {
	tests := []struct {
		lines string
		url   string
		score int
	}{
		{
			`http://www.bad.com,2`,
			`www.bad.com`,
			2,
		},
		{
			`https://www.known.com,10`,
			`www.unknown.com`,
			0,
		},
		{
			`https://www.known.com,5
https://www.bad.com,10`,
			`www.known.com`,
			5,
		},
	}

	for _, test := range tests {
		s := strings.NewReader(test.lines)
		scorer, _ := New(s)
		ds := scorer.GetDomainScore(test.url)
		if ds.Score != strconv.Itoa(test.score) {
			t.Fatalf("Want %v Got %v", ds.Score, test.score)
		}
	}
}

func TestUpdateDomain(t *testing.T) {
	type update struct {
		domain string
		score  int64
	}
	tests := []struct {
		lines   string
		updates []update
		url     string
		score   int
	}{
		{
			lines: `http://www.known.com,2`,
			updates: []update{
				{
					`www.known.com`, 5,
				},
			},
			url:   `www.known.com`,
			score: 5,
		},
		{
			lines: `http://www.known.com,2`,
			updates: []update{
				{
					`www.new.com`, 5,
				},
			},
			url:   `www.known.com`,
			score: 2,
		},
		{
			lines: `http://www.known.com,2`,
			updates: []update{
				{
					`www.new.com`, 5,
				},
				{
					`www.new.com`, 3,
				},
				{
					`www.new.com`, 9,
				},
			},
			url:   `www.new.com`,
			score: 9,
		},
	}

	for _, test := range tests {
		s := strings.NewReader(test.lines)
		scorer, _ := New(s)
		for _, nd := range test.updates {
			scorer.AddDomain(nd.domain, nd.score)
		}
		ds := scorer.GetDomainScore(test.url)
		if ds.Score != strconv.Itoa(test.score) {
			t.Fatalf("Want %v Got %v", test.score, ds.Score)
		}
	}
}
