### Description
A simple api service that provides reputations for queried domains. The service
uses a radix tree as the backing datastructure. The radix tree provides an
efficient storage solution for domains which can vary by a few characters as we
often see with malicious domains. 

The service implicitly allows domains that are unrecognized. The design
decision here is to only track and stop accesses to known malicious domains
instead of building reputations for all possible domains.

### Routes
A `GET` at `/urlinfo/url/{url}` can be used to query for a domain and its score
A `POST` at `/urlinfo/update` with the request body containing ``{"url":"<url>", "score":"<score>"}`` (both as strings) will update or add a new domain and its score.

#### Examples

Request: `GET /urlinfo/url/www.google.com`
Response:
```
{
    "Domain": "www.google.com",
    "Score": "0"
}
```

Request:
```
POST http://Infra-Farga-U6XG7DG3WFYF-764537501.us-east-1.elb.amazonaws.com/urlinfo/update
...
Body:
{
    "url": "www.google.com",
    "score": "10"
}
```

Response:
```
{
    "Message": "Domain updated successfully"
}
```

### Files
1. `domains.go` Abstraction around [Hashicorp's immutable radix trie](https://pkg.go.dev/github.com/hashicorp/go-immutable-radix)
datastructure. Exposes read and update APIs only
2. `handler.go` Intermediary layer between `domains.go` and `main.go` 
that provides handlers for incoming api calls
3. `main.go` Sets up an api path router and wires up with routes with handlers.
4. `router.go` Abstraction layer around [Gorilla Mux](https://github.com/gorilla/mux)



### TODO
1. Currently only supports exact matches. Use available prefix match apis to support subdomains.
2. Updates are not persistent. Introduce a persistent storage solution like DynamoDB.
3. Use RPC for queries instead of HTTP. Since we expect this service to be called for any outbound request, performance is critical.
4. Switch to an immutable radix tree and use CAS operations to improve performance. Currently the system synchronizes access to the tree when
   updates are made to it. Instead of using a lock, build the tree outside a synchronization block and only swap the pointer as a CAS operation.
