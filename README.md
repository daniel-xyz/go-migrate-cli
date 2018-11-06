A simple CLI library that let's you easily work with PostgreSQL migration files. It is intended to be wrapped by your own application.

It uses https://github.com/golang-migrate/migrate under the hood. Might supprt the same amount of databases at a later date when demand is there.

![](demo.gif)

### How to use

1. download the package: `go get github.com/Flur3x/go-migrate-cli`

2. use it in your own appliaction:

```go
import "github.com/Flur3x/go-migrate-cli"

func main() {
	db, err := database.Connection()

	if err != nil {
		panic(err)
	}

	migrate.CLI(db, "my-database-name", "cmd/migration/sources")
}
```

3. start your app, the CLI will guide you. :)

### Be warned

This is one of my first go projects and I wouldn't use it in production yet.
