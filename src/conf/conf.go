package conf

import (
	"log"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

type config struct {
	Nats struct {
		URL       string
		ClusterID string
		ClientID  string
		Subject   struct {
			PosPostDB string
		}
		PublisherNum int
	}
	PostDB struct {
		Dbms    string
		User    string
		Pass    string
		Host    string
		Port    uint16
		Charset string
	}
}

var C config

func init() {

	viper.SetConfigName("conf")
	viper.SetConfigType("yml")
	viper.AddConfigPath("conf")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Fatal(err)
	}

	spew.Dump(C)
}
