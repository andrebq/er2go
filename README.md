# er2go - Erlang/Elixir 2 Go

This package provides some facilities to simplify data exchange between
Go and BEAM (Erlang, Elixir, etc...).

The most important sub-package is [etf](./etf/) which implements a very
basic encoder that encode/decode data using a subset of the [etf spec](./etf/spec.md).

The SPEC file was copied from [Erlang/OTP OSS Repo](https://github.com/erlang/otp/blob/OTP-27.3.2/erts/doc/guides/erl_ext_dist.md#L1).
To simplify the initial implementation only a sub-set of all possible etf types
can be processed. You can check the [generate.exs](./etf/testdata/generate.exs)
to have an idea of what datatypes can be processed.

Another key detail is that er2go attempts to be round-trip capable, ie.:
after decoding a binary etf, the package should be able to re-encode the
data into the same binary representation.

## Why creating this project?

I think it is philosophically interesting to have two languages, designed to
simplify concurrency, being able to talk to each other using a format that is
native to at least one of them. Given I am a better Go developer than 
an Erlang one, implementing etf in Go was easier than implementing Gob 
in Erlang/Elixir.

Then, from that premise, I started looking for etf encoders/decoders in Go,
but most projects were too old (without any recent activity). While I could have
started by: forking those projects and cleaning them up, instead I choose to
write the encoders/decoders myself to have a better understanding of the code,
as I intend to actually use it to build other things.

## How to build

Using `go build ./...` should be enough to build everything. But there is a
[Makefile](./Makefile) to help with other tasks.

Test data is generated via [generate.exs](./etf/testdata/generate.exs), any
recent version of `elixir` should suffice, but using `make regen-testdata` is
preferred as it will use a specific version of elixir (docker is needed for that).

## Is it production ready?

All tests should pass without issues. All supported types are defined in [generate.exs](./etf/testdata/generate.exs).
The codebase lacks fuzzy tests and it has not been optimized, you can check some of
the benchmarks to get a feeling of how it behaves.

### Round trip guarantees

I will try to keep it for as long as possible. In the event that round-trip is not
achievable anymore, it will be considered a breaking change and major version
will be created.

## Contributing

I work on this codebase during my free time, between work, personal life and other
fun (non-tech) projects. I will read issues and review PRs but response time
could vary from hours to weeks.