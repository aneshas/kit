<p align="center">
<img src="docs/img/logo.png" alt="tonto kit logo" title="tonto kit logo" width="400" />
</p>

# kit
[![wercker status](https://app.wercker.com/status/684f3efd53b66300c2470d2f0a6c2bd4/s/master "wercker status")](https://app.wercker.com/project/byKey/684f3efd53b66300c2470d2f0a6c2bd4)
[![codecov](https://codecov.io/gh/tonto/kit/branch/master/graph/badge.svg)](https://codecov.io/gh/tonto/kit)
[![Go Report Card](https://goreportcard.com/badge/github.com/tonto/kit)](https://goreportcard.com/report/github.com/tonto/kit)

This package contains common packages used accross server based projects

# Packages included
* [goapp/cli](goapp/) - Go application package layout reference and cli supporting microservice workflow, DDD oriented with clean code/concerns separation (kubernetes, prometheus, twirp)
* [tx](tx/) - simple transactional abstraction
* [http](http/) - http server implementation with server lifecycle control, commonly used 
adapters, easy service registration and response/error handling, gracefull shutdown, tls support...
