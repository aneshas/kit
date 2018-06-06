# http
### `import “github.com/tonto/kit/http”`
Package http provides server implementation with lifecycle control, graceful shutdown, commonly used 
adapters, easy service registration and response/error handling, endpoints implementation, tls support...

See the [example](example/) package for example usage and to see how you can integrate twirp.

# Creating and running the server
With minimal setup:
```go
server := http.NewServer()

server.RegisterServices(
	...
)

log.Fatal(server.Run(8080))
```

With server options (see the [adapter](adapter/) package):
```go
logger := log.New(os.Stdout, "http/example => ", log.Ldate|log.Ltime|log.Lshortfile)

server := http.NewServer(
  http.WithLogger(logger),
  http.WithTLSConfig("cert.pem", "key.pem"),
  http.WithMux(customMux), // Override default gorilla mux router
  http.WithNotFoundHandler(notFoundHdlr),
  http.WithAdapters(
    adapter.WithRequestLogger(logger),
    adapter.WithCORS(
      adapter.WithCORSAllowOrigins("*"),
      adapter.WithCORSAllowMethods("PUT", "DELETE"),
      adapter.WithCORSMaxAge(86400),
    ),
  ),
)

server.RegisterServices(
  ...
)

log.Fatal(server.Run(8080))
```

You can use `server.Stop()` to explicitly stop the server.

# Services
With this package, there is a notion of service which is simply a type that implements `http.Service`

## Creating basic http service
Creating a service is as easy as embeding `http.BaseService`, and registering an endpoint.
```go
type OrderService struct {
  http.BaseService
}
```

## Service routing prefix
Implement `Prefix` method to return service routing prefix ("/" by default)
```go
func (os *OrderService) Prefix() string { return "order" }
```

## HandlerFuncs and Endpoints
With services there is a notion of HandlerFunc and an Endpoint.
HandlerFunc is basically a standard http handler func with context added as first param:

`type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)`

You can use it the same way you would a normal http handler func with no restrictions,
or caveats (you only have the extra context passed in).

Extending service to register HandlerFunc:

```go
func NewOrderService() *Order {
  svc := Order{}

  svc.RegisterHandler("POST", "/create", svc.create)

  return &svc
}

type OrderService struct {
  http.BaseService
}

func (os *OrderService) Prefix() string { return "order" }

func (os *OrderService) create(c context.Context, w ghttp.ResponseWriter, r *ghttp.Request) {
  respond.WithJSON(w, r, http.NewResponse("order created", ghttp.StatusOK))
}
```

See the [respond](respond/) package for more info on it's usage.

### Endpoints
Endpoints however, are specificaly designed to be used for json api endpoints.
They provide easy request decoding and response encoding.

Endpoint is a func of the following signature:

`func(c context.Context, w http.ResponseWriter, req *CustomType) (*http.Response, error)`

Instead of go http request the third parameter is a custom type that you choose to which 
the request body will be json decoded.

Return parameters are `http.Response` and an `error`

Let's extend our service with two new endpoints:
```go
func NewOrderService() *Order {
  svc := Order{}

  svc.RegisterHandler("POST", "/create", svc.create)
  svc.RegisterEndpoint("POST", "/add_item", svc.addItem)

  return &svc
}

type OrderService struct {
  http.BaseService
}

func (os *OrderService) Prefix() string { return "order" }

func (os *OrderService) create(c context.Context, w ghttp.ResponseWriter, r *ghttp.Request) {
  respond.WithJSON(w, r, http.NewResponse("order created", ghttp.StatusOK))
}

type addItemReq struct {
  ItemID  int64 `json:"item_id"`
  OrderID int64 `json: "order_id"`
  UserID  int64 `json:"user_id"`
}

func (os *OrderService) addItem(c context.Context, w ghttp.ResponseWriter, r *addItemReq) (*http.Response, error) {
  // r now contains the decoded json request
  return http.NewResponse("item added", ghttp.StatusOK), nil
}
```

If you don't expect any request body from an endpoint you can omit the `*http.Response` return value,
and only keep the error (which cannot be omited).

You can return a regular go `error` or use `http.NewError` to create a composite error with custom status.
Both will be correctly json encoded.

## Service adapters
You can use `svc.Adapt(...adapters)` to register per service adapters.
Check out [example](example/) package for an example.

## Putting it all together
The only thing that is left is to register our service with the server:
```go
server := http.NewServer()

server.RegisterServices(
	NewOrderService(),
)

log.Fatal(server.Run(8080))
```
