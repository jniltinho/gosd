// To run:
// go get github.com/githubnemo/CompileDaemon
// CompileDaemon -command="./gervice"

package gosd

import "os"
import "fmt"
import "time"
import "strconv"
import redis "gopkg.in/redis.v3"

type DriverRedis struct {}

var RedisClient *redis.Client

func (this DriverRedis) Start(name, url string) string {
  // start
  redisDB,_ := strconv.Atoi(os.Getenv("gosdRedisDB"))
  RedisClient = redis.NewClient(&redis.Options{
        Addr:     os.Getenv("gosdRedisAddr"),
        Password: os.Getenv("gosdRedisPassword"), // no password set
        DB:       int64(redisDB),  // use default DB
    })
  _, err := RedisClient.Ping().Result()
  if err != nil {
    fmt.Println("Error connecting with Redis.")
    fmt.Println(err.Error())
    return "standalone-" + name
  }

  // set
  currentName := registerService(name, url)
  return currentName
}

func (this DriverRedis) Get() (map[string]string, error) {
  return RedisClient.HGetAllMap("gosd").Result()
}

func (this DriverRedis) Delete(currentName string) {
  RedisClient.HDel("gosd", currentName)
}


func registerService(basicName, url string) string {
  finalServiceName := ""
  created := false
  for created != true {
    finalServiceName = basicName + "-" + time.Now().Format("20060102150405.99999999")
    created,_ = RedisClient.HSetNX("gosd", finalServiceName, url).Result()
  }
  return finalServiceName
}
