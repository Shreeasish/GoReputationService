package main

import "testing"

func TestInitializeFromPath(t *testing.T) {
	d := map[string]int{"www.baddomain.csv": 1,
		"www.evildomain.csv": 5,
		"www.java.com":       10,
		"www.log4j.com":      9,
	}

	rp := InitializeFromPath("./test_resources/test_domains.csv")
	it := rp.rt.Root().Iterator()
	if rp.Count != 4 {
		t.Fatalf("Assets read incorrectly")
	}

	for key, _, ok := it.Next(); ok; key, _, ok = it.Next() {
		if _, ok := d[string(key)]; !ok {
			t.Fatalf("Assets read incorrectly")
		}
	}
}

func TestGetDomain(t *testing.T) {
	rp := InitializeFromPath("./test_resources/test_domains.csv")
	ds, ok := rp.GetDomainScore("www.java.com")
	if !ok || ds.Domain != "www.java.com" {
		t.Fatalf("Unable to retrieve domain")
	}

	if ds, ok = rp.GetDomainScore("www.gooddomain.com"); ok {
		t.Fatalf("Non existent domain returned ok")
	}
}
