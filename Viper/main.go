// go get github.com/spf13/viper
package main

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func main() {
	Example2()
}

func Example1() {
	viper.SetConfigFile("./etc/etc.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %v", err))
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	viper.WatchConfig()

	select {}
}

func Example2() {
	// Name of the config file without an extension (Viper will intuit the type from an extension on the actual file)
	viper.SetConfigName("etc")

	// Add search paths to find the file
	viper.AddConfigPath("./etc/")
	viper.AddConfigPath(".")

	// Find and read the config file
	err := viper.ReadInConfig()

	// Handle errors
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	redisCfg := viper.Get("server.redis")
	fmt.Println(redisCfg)
	redisPort := viper.Get("server.redis.port")
	fmt.Println(redisPort)
	redisDbViper := viper.Sub("server.redis")
	fmt.Println(redisDbViper.GetString("addr"))
	fmt.Println(redisDbViper.GetInt("port"))
	fmt.Println(redisDbViper.GetInt("db"))

}

func Example3() {
	viper.SetConfigFile("./etc/etc.yaml")
	all := viper.AllSettings()
	bs, err := yaml.Marshal(all)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}
	fmt.Println(bs)
}

func Example4() {
	// 读取不到
	fmt.Println(viper.Get("path"))

	// 开始读取环境变量
	viper.AutomaticEnv()
	// 读取环境变量的值
	fmt.Println(viper.Get("path"))
}
