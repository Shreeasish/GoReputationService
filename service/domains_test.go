package main

import "testing"

func TestHelloWorld(t *testing.T) {
	// t.Fatal("not implemented")
}

func TestInitializeFromPath(t *testing.T) {
	d := map[string]int{"www.baddomain.csv": 1,
											"www.evildomain.csv": 5,
											"www.java.com": 10,
											"www.log4j.com": 9,
										}
	
	r := InitializeFromPath("./test_resources/test_domains.csv")
	it := r.radixTree.Root().Iterator()
	if (r.Count != 4) {
		t.Fatalf("Assets read incorrectly")
	}

	for key, _, ok := it.Next(); ok; key, _, ok = it.Next() {
		if _, ok := d[string(key)]; !ok {
			t.Fatalf("Assets read incorrectly")
		}
	}
}
