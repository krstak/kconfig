# Kconfig

Yaml file `my-file.yaml`:

```
addr: ":8080"
database:
  use: true
  dialect: "postgres"
  timeout: 12
  secret: env[SECRET_PASS|super-secret]
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

// returns bool true
c.GetBool("database.use")

// returns string "postgres"
c.Get("database.dialect")

// returns int 12
c.GetInt("database.timeout")

// returns slice []string{"first", "second", "third"}
c.GetSlice("list")

// returns a value of environment variable GTM_ENV as a string
c.Get("gtm")

// returns a value of environment variable SECRET_PASS as a string. If not exists, it returns a default value super-secret
c.Get("database.secret")
```