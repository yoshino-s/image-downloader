package s3

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/yoshino-s/go-framework/common"
	"github.com/yoshino-s/go-framework/configuration"
)

var _ configuration.Configuration = (*config)(nil)

type config struct {
	Endpoint                  string        `mapstructure:"endpoint"`
	AccessKeyID               string        `mapstructure:"access_key_id"`
	SecretAccessKey           string        `mapstructure:"secret_access_key"`
	Insecure                  bool          `mapstructure:"insecure"`
	Region                    string        `mapstructure:"region"`
	Bucket                    string        `mapstructure:"bucket"`
	PresignedGetObjectExpires time.Duration `mapstructure:"presigned_get_object_expires"`
	Public                    bool          `mapstructure:"public"`
}

func (c *config) Register(set *pflag.FlagSet) {
	set.String("s3.endpoint", "", "s3 endpoint")
	set.String("s3.access_key_id", "", "s3 access key id")
	set.String("s3.secret_access_key", "", "s3 secret access key")
	set.Bool("s3.insecure", false, "s3 insecure")
	set.String("s3.region", "", "s3 region")
	set.String("s3.bucket", "", "s3 bucket")
	set.Duration("s3.presigned_get_object_expires", 7*24*time.Hour, "s3 presigned get object expires")
	set.Bool("s3.public", false, "s3 public")

	common.MustNoError(viper.BindPFlags(set))
	configuration.Register(c)
}

func (c *config) Read() {
	common.MustDecodeFromMapstructure(viper.AllSettings()["s3"], c)
}
