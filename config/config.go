package config

import "flag"

var DataBucket = flag.String(`data_bucket`, ``, `The S3 bucket name`)
