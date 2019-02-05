package tile38

import (
	"fmt"

	"github.com/go-redis/redis"
)

// RedisConfig adapter configuration
type RedisConfig struct {
	HostName string
	Port     uint
	Password string
	DB       int
}

// Tile38Adapter struct
type Tile38Adapter struct {
	cfg    *RedisConfig
	client *redis.Client
}

type SubMessage struct {
	Channel string
	Payload string
}

// New ititialize the adapter
func (tile38 *Tile38Adapter) New(config *RedisConfig) error {

	tile38.cfg = config

	tile38.client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", tile38.cfg.HostName, tile38.cfg.Port),
		Password: tile38.cfg.Password,
		DB:       tile38.cfg.DB,
	})

	_, err := tile38.client.Ping().Result()

	if err != nil {
		return err
	}

	//Set output type
	cmd := tile38.client.Do("OUTPUT", "JSON")

	res, err := cmd.Result()

	if err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

// RunCommand runs a redis command
func (tile38 *Tile38Adapter) RunCommand(q *Command) (interface{}, error) {

	cmd := redis.NewCmd(q.QueryParts...)

	tile38.client.Process(cmd)

	result, err := cmd.Result()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (tile38 *Tile38Adapter) SubscribeChannel(receiver chan *SubMessage, redisChan string) {

	fmt.Println(redisChan)
	pubsub := tile38.client.Subscribe(redisChan)

	for {
		msgi, err := pubsub.Receive()
		if err != nil {
			fmt.Println(err)
		}

		switch msg := msgi.(type) {
		case *redis.Message:
			receiver <- &SubMessage{
				Channel: msg.Channel,
				Payload: msg.Payload,
			}
		default:
			fmt.Println(nil, "Got control message", msg)
		}

	}

}
