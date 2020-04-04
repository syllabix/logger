# Logger

This project is a wrapper around the excellent logging framework [zap](https://github.com/uber-go/zap). It provides some opinionated encoders as well as an implementation for a redis sink.

As an aside - the module also includes a work in progress mechanism for setting log level per package.

# Example Use

For an in depth overview - please check the [zap documentation](https://godoc.org/go.uber.org/zap)

#### Local Development Example

```go
// assumed imports

func main() {
    // configure global settings that will
	// apply to all instance of a logger
	logger.Configure(
		logger.AppName("my-app"),
		logger.Level(zapcore.DebugLevel),
		logger.Mode(mode.Development),
	)

	// New() returns a logger with a strongly typed logging context intended for use in performance critical paths
	log := logger.New()

	log.Info("hello world, this is some info",
		zap.String("field", "value"),
		zap.Int64("count", 12315))

	// Slightly less performant, but simpler to use, the "Sugar" logger provides a loosely typed field context
	sLog := logger.New().Sugar()

	sLog.Warnw("this logger is sugared",
		"count", 1234,
		"flavor", "tasty",
    )
}

```

#### Logging to remote redis sink with JSON encoded log output

```go
// assumed imports

func main() {

    p, err := newpool("localhost:6379", "password123")

	rsink := redis.NewSink("logstash.stg.bolcom", p)
	if err != nil {
		log.Fatal("is redis running?", err)
	}

    // all logs will write to the local console as well as to the provided redis instance
	logger.Configure(
		logger.AppName("test-app"),
		logger.Level(zapcore.DebugLevel),
		logger.Mode(mode.Production),
		logger.JSONWriter(rsink),
    )

    log := logger.New().Sugar()

    // log away!


}

// example func for setting up a redigo redis pool
func newpool(addr, password string) (*redigo.Pool, error) {
	pool := &redigo.Pool{
		Dial: func() (redigo.Conn, error) {
			conn, err := redigo.Dial("tcp",
				addr,
				redigo.DialPassword(password),
				redigo.DialConnectTimeout(time.Second*5),
			)

			if err != nil {
				return nil, err
			}

			return conn, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	conn := pool.Get()
	defer conn.Close()

	err := conn.Send("PING")
	if err != nil {
		return nil, err
	}

	return pool, nil
}


```

#### Changing log level for all loggers in a package

```go
// assumed imports


func main() {

    // ... log config already set up

    err = logger.SetLevelForPackage("main", zapcore.WarnLevel)
	if err != nil {
		log.Warn(err.Error())
	}
}

```