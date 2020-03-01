# GoTo: URL Shortner

An implementation of https://github.com/kellegous/go. A internal focused URL shortener.

```
go get github.com/clintjedwards/goto
```

## DNS/DHCP Setup

To enable functionality such as `go/mylink`, you'll need two pieces of functionality.

1. The application must be be given an A record: Ex. `go.clintjedwards.home`
2. Through DHCP you must set up search domains that include the (example) `clintjedwards.home` domain
3. Now simply typing `go/{something here}` will take you to your shortened link

## API Documentation

| Route       | Methods     | Payload   | Returns                       |
| ----------- | ----------- | --------- | ----------------------------- |
| /links      | GET         | None      | [{url, id, hits, created}]    |
| /links/{id} | GET, DELETE | None      | {url, id, hits, created}, nil |
| /create     | POST        | {url, id} | {url, id, hits, created}      |
| /{id}       | GET         | None      | 302/Redirect                  |

## Usage

Examples use [httpie](https://httpie.org/)

```
http POST localhost:8080/create url=https://yahoo.com id=test
http GET localhost:8080/links
http GET localhost:8080/test
http GET localhost:8080/links/test
http DELETE localhost:8080/links/test
```

### Reserved links

The following short names are reserved for app use: ["links", "create", "version", "status", "health", "edit", "api"]

## Authors

- **Clint Edwards** - [Github](https://github.com/clintjedwards)
