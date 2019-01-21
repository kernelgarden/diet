package constant

var CtxDbName = "db_ctx"
var CtxConfig = "global_config"

type Config struct {
	Database struct {
		Driver   string
		Username string
		Password string
		Host	 string
		Port 	 string
		Name     string
		Logger   string
	}

	Behaviorlog struct {
		Kafka string
	}

	Debug    bool
	Service  string
	Httpport string
}
