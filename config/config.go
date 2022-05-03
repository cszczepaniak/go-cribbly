package config

import (
	"flag"
	"time"

	"github.com/rakyll/globalconf"
)

var (
	DataBucket       = flag.String(`data_bucket`, ``, `The S3 bucket name`)
	ByteStoreTimeout = flag.Duration(`byte_store_timeout`, time.Second, `The timeout for byte store requests`)
)

func Init() {
	config, err := globalconf.NewWithOptions(&globalconf.Options{
		EnvPrefix: `CRIBBLY_`,
	})
	if err != nil {
		panic(err)
	}
	config.Parse()
}
