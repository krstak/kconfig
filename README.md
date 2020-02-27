# Kconfig

Yaml file `my-file.yaml`:

```
addr: ":8080"
database:
  dialect: "postgres"
  timeout: 12
list: 
  - "first"
  - "second"
  - "third"
gtm: "env[GTM_ENV]"
```

Usage:

```go
import "github.com/krstak/kconfig"

c, err := kconfig.Load("./my-file.yaml")

// returns string ":8080"
c.Get("addr")

// returns string "postgres"
c.Get("database.dialect")

// returns int 12
c.GetInt("database.timeout")

// returns slice []string{"first", "second", "third"}
c.GetSlice("list")

// returns a value of environment variable GTM_ENV as a string
c.Get("gtm")
```