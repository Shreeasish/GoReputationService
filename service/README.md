### Files
1. `domains.go` Abstraction around [Hashicorp's immutable radix trie](https://pkg.go.dev/github.com/hashicorp/go-immutable-radix)
datastructure. Exposes read and update APIs only
2. `handler.go` Intermediary layer between `domains.go` and `main.go` 
that provides handlers for incoming api calls
3. `main.go` Sets up an api path router and wires up with routes with handlers.
4. `router.go` Abstraction layer around [Gorilla Mux](https://github.com/gorilla/mux)

### TODO
1. Restructure into `domains` and `router` into their own modules
2. Rename `main.go` to `wiring.go` or similar
3. Integration tests for
  1.`router.go` and `main.go`
  2.`main.go` and `handler.go`
4. Currently only supports exact matches. Use available prefix match apis to support subdomains.
