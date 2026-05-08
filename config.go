package main
import("log"
		"github.com/spf13/viper")
type Configuration struct{
	Port string `mapstructure:"Port"`
	DBURL string `mapstructure:"DB_URL"`
}
func LoadConfiguration()(Configuration Config,err error)
{
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err= viper.ReadInConfig()
	if err != nil
	{
		log.Println(".env file not found, using system environment variables")
	}

	err = viper.Unmarshal(&Configuration)
}