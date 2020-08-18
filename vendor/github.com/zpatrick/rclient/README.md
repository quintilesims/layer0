# RClient

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/zpatrick/rclient/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/zpatrick/rclient)](https://goreportcard.com/report/github.com/zpatrick/rclient)
[![Go Doc](https://godoc.org/github.com/zpatrick/rclient?status.svg)](https://godoc.org/github.com/zpatrick/rclient)

## Getting Started
Checkout the [Examples](https://github.com/zpatrick/rclient/tree/master/examples) folder for some working examples.
The following snippet shows RClient interacting with Github's API:
```
package main

import (
        "github.com/zpatrick/rclient"
        "log"
)

type Repository struct {
        Name        string `json:"name"`
}

func main() {
        client := rclient.NewRestClient("https://api.github.com")

        var repos []Repository
        if err := client.Get("/users/zpatrick/repos", &repos); err != nil {
                log.Fatal(err)
        }

        log.Println(repos)
}
```

## Request Options
Requests can be configured using a [Request Option](https://godoc.org/github.com/zpatrick/rclient#RequestOption).
A `RequestOption` is simply a function that manipulates an `http.Request`.
You can create request options like so:
```
setProto := func(req *http.Request) error {
    req.Proto = "HTTP/1.0"
    return nil
}

client.Get("/path", &v, setProto)
```

The built-in request options are described below.

#### Header / Headers
The `Header()` and `Headers()` options add header(s) to a `*http.Request`.
```
// add a single header
client.Get("/path", &v, rclient.Header("name", "val"))

// add multiple headers
client.Get("/path", &v, rclient.Header("name1", "val1"), rclient.Header("name2", "val2"))
client.Get("/path", &v, rclient.Headers(map[string]string{"name1": "val1", "name2":"val2"}))
```

#### Basic Auth
The `BasicAuth()` option adds basic auth to a `*http.Request`.
```
client.Get("/path", &v, rclient.BasicAuth("user", "pass"))
```

#### Query
The `Query()` options adds a query to a `*http.Request`.
```
query := url.Values{}
query.Add("name", "John")
query.Add("age", "35")

client.Get("/path", &v, rclient.Query(query))
```

**NOTE**: This can also be accomplished by adding the raw query to the `path` argument
```
client.Get("/path?name=John&age=35", &v)
```

## Client Configuration
The `RestClient` can be configured using the [Client Options](https://godoc.org/github.com/zpatrick/rclient#ClientOption) described below.

#### Doer
The `Doer()` option sets the `RequestDoer` field on the `RestClient`. 
This is the `http.DefaultClient` by default, and it can be set to anything that satisfies the [RequestDoer](https://godoc.org/github.com/zpatrick/rclient#RequestDoer) interface. 
```
client, err := rclient.NewRestClient("https://api.github.com", rclient.Doer(&http.Client{}))
```

#### Request Options
The `RequestOptions()` option sets the `RequestOptions` field on the `RestClient`.
This will manipulate each request made by the `RestClient`.
This can be any of the options described in the [Request Options](#request-options) section. 
A typical use-case would be adding headers for each request.
```
options := []rclient.RequestOption{
    rclient.Header("name", "John Doe").
    rclient.Header("token", "abc123"),
}

client, err := rclient.NewRestClient("https://api.github.com", rclient.RequestOptions(options...))
```

#### Builder
The `Builder()` option sets the `RequestBuilder` field on the `RestClient`.
This field is responsible for building `*http.Request` objects. 
This is the [BuildJSONRequest](https://godoc.org/github.com/zpatrick/rclient#BuildJSONRequest) function by default, and it can be set to any [RequestBuilder](https://godoc.org/github.com/zpatrick/rclient#RequestBuilder) function.
```
builder := func(method, url string, body interface{}, options ...RequestOption) (*http.Request, error){
    req, _ := http.NewRequest(method, url, nil)
    for _, option := range options {
		if err := option(req); err != nil {
			return nil, err
		}
	}
	
    return nil, errors.New("I forgot to add a body to the request!")
}

client, err := rclient.NewRestClient("https://api.github.com", rclient.Builder(builder))
```

#### Reader
The `Reader()` option sets the `ResponseReader` field on the `RestClient`.
This field is responsible for reading `*http.Response` objects. 
This is the [ReadJSONResponse](https://godoc.org/github.com/zpatrick/rclient#ReadJSONResponse) function by default, and it can be set to any [ResponseReader](https://godoc.org/github.com/zpatrick/rclient#ResponseReader) function.
```
reader := func(resp *http.Response, v interface{}) error{
    defer resp.Body.Close()
    return json.NewDecoder(resp.Body).Decode(v)
}

client, err := rclient.NewRestClient("https://api.github.com", rclient.Reader(reader))
```

# License
This work is published under the MIT license.

Please see the `LICENSE` file for details.
