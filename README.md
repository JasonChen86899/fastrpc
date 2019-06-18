# Fast rpc<br/>

**Fast rpc** is a **framework** for building microservices
or servermesh in Go. We solve common problems in distributed
systems and application architecture so you can focus on delivering
business value.

- Mailing: [fast rpc](chenhao86899@gmail.com)
- Slack: [gophers.slack.com](https://gophers.slack.com) **#fast rpc** ([invite](https://gophersinvite.herokuapp.com/))

## Motivation

Go has emerged as the language of the server, but it remains underrepresented
in so-called "modern enterprise" companies like Facebook, Twitter, Netflix, and
SoundCloud. Many of these organizations have turned to JVM-based stacks for
their business logic, owing in large part to libraries and ecosystems that
directly support their microservice architectures.

To reach its next level of success, Go needs more than simple primitives and
idioms. It needs a service mesh, for coherent distributed programming
in the large. Fast rpc is a best practices, which provide a
comprehensive, robust, and trustable way of building microservices for
organizations of any size.


## Goals

- Operate in a heterogeneous SOA — expect to interact with mostly non-Go-kit services
- RPC as the primary messaging pattern
- Pluggable serialization and transport — not just JSON over HTTP
- Operate within existing infrastructures — no mandates for specific tools or technologies

## Non-goals

- Supporting messaging patterns other than RPC (for now) — e.g. MPI, pub/sub, CQRS, etc.
- Re-implementing functionality that can be provided by adapting existing software
- Having opinions on operational concerns: deployment, configuration, process supervision, orchestration, etc.

## Contributing

Please see [CONTRIBUTING.md](/CONTRIBUTING.md).

## Dependency management

Fast rpc is a library, designed to be imported into a binary package. Go mod is
the way we recommend to import this package to your project

## Related projects

Projects with a ★ have had particular influence on Go kit's design (or vice-versa).

### Service frameworks

- [go-micro](https://github.com/myodc/go-micro), a microservices client/server library ★

### Individual components

- [afex/hystrix-go](https://github.com/afex/hystrix-go), client-side latency and fault tolerance library
- [armon/go-metrics](https://github.com/armon/go-metrics), library for exporting performance and runtime metrics to external metrics systems
- [codahale/lunk](https://github.com/codahale/lunk), structured logging in the style of Google's Dapper or Twitter's Zipkin
- [eapache/go-resiliency](https://github.com/eapache/go-resiliency), resiliency patterns
- [sasbury/logging](https://github.com/sasbury/logging), a tagged style of logging
- [grpc/grpc-go](https://github.com/grpc/grpc-go), HTTP/2 based RPC
- [inconshreveable/log15](https://github.com/inconshreveable/log15), simple, powerful logging for Go ★
- [mailgun/vulcand](https://github.com/vulcand/vulcand), programmatic load balancer backed by etcd
- [mattheath/phosphor](https://github.com/mondough/phosphor), distributed system tracing
- [pivotal-golang/lager](https://github.com/pivotal-golang/lager), an opinionated logging library
- [rubyist/circuitbreaker](https://github.com/rubyist/circuitbreaker), circuit breaker library
- [sirupsen/logrus](https://github.com/sirupsen/logrus), structured, pluggable logging for Go ★
- [sourcegraph/appdash](https://github.com/sourcegraph/appdash), application tracing system based on Google's Dapper
- [spacemonkeygo/monitor](https://github.com/spacemonkeygo/monitor), data collection, monitoring, instrumentation, and Zipkin client library
- [streadway/handy](https://github.com/streadway/handy), net/http handler filters
- [vitess/rpcplus](https://godoc.org/github.com/youtube/vitess/go/rpcplus), package rpc + context.Context
- [gdamore/mangos](https://github.com/gdamore/mangos), nanomsg implementation in pure Go

## Additional reading

- [Architecting for the Cloud](https://slideshare.net/stonse/architecting-for-the-cloud-using-netflixoss-codemash-workshop-29852233) — Netflix
- [Dapper, a Large-Scale Distributed Systems Tracing Infrastructure](http://research.google.com/pubs/pub36356.html) — Google
- [Your Server as a Function](http://monkey.org/~marius/funsrv.pdf) (PDF) — Twitter

---