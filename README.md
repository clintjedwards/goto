# GoTo: URL Shortener

An implementation of https://github.com/kellegous/go. An internal focused URL shortener.

```
go get github.com/clintjedwards/goto
```

## DNS/DHCP Setup

To enable functionality such as `go/mylink`, you'll need two pieces of functionality.

1. The application must be given an A record: Ex. `go.clintjedwards.home`
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

```golang
// Normal links work just how you expect
http POST localhost:8080/create url="https://github.com" id="github" // normal link
http GET localhost:8080/github                                       // Use ID to redirect to full URL
http GET localhost:8080/github?tab=repositories                      // query params are passed to the full URL

// Formatted links allow you to substitute variables that might be in the middle of a link
http POST localhost:8080/create url="https://github.com/clintjedwards/{}/issues" id="github"
http GET localhost:8080/github/release  // Returns a link to: https://github.com/clintjedwards/release/issues

http GET localhost:8080/links           // View all links
http GET localhost:8080/links/test      // View specific link details
http DELETE localhost:8080/links/test   // Remove a link
```

### Reserved links

The following short names are reserved for app use: ["links", "create", "version", "status", "health", "edit", "api"]

## Authors

- **Clint Edwards** - [Github](https://github.com/clintjedwards)
