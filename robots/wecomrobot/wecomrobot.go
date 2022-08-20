package wecomrobot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hirpc/arsenal/hihttp"
)

type wecomrobot struct {
	endpoint string
}

// DefaultRobot can be used for general usage.
// If the construction function called, an option can be provided to setup the default robot.
// Currently, it can be used for DirectSend() as a name to illustrate which robot should be used to delivery.
const DefaultRobot = "default"

var (
	registeredRobots = make(map[string]*wecomrobot)
	mu               sync.RWMutex
)

func New(endpoint string, opts ...Option) *wecomrobot {
	opt := Options{
		defaultRegistration: false,
	}
	for _, o := range opts {
		o(&opt)
	}

	r := &wecomrobot{
		endpoint: endpoint,
	}
	if opt.defaultRegistration {
		registeredRobots[DefaultRobot] = r
	}
	return r
}

func (w wecomrobot) Send(ctx context.Context, msg string) error {
	res, err := hihttp.New(
		hihttp.WithRetryWait(time.Millisecond*300),
		hihttp.WithRetryCount(2),
	).Post(
		ctx,
		w.endpoint,
		hihttp.NewJSONPayload(
			fmt.Sprintf(
				`{"msgtype":"markdown","markdown":{"content":"%v"}}`, msg,
			),
		),
	)
	if err != nil {
		return err
	}
	var response struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}
	if err := json.Unmarshal(res, &response); err != nil {
		return err
	}
	if response.Errcode != 0 {
		return errors.New(response.Errmsg)
	}
	return nil
}

// ErrRobotName represents there is no valid named robot that has been registered before.
var ErrRobotName = errors.New("no valid robot found, please check the name of robot you provided")

// Register can be called multiple times so that many rebots can be acceptable message delivers.
func Register(name string, robot *wecomrobot) {
	mu.Lock()
	defer mu.Unlock()
	registeredRobots[name] = robot
}

func Get(name string) *wecomrobot {
	mu.RLock()
	defer mu.RUnlock()
	if v, ok := registeredRobots[name]; ok {
		return v
	}
	return nil
}

type configuration struct {
	Robots []struct {
		Name     string `json:"name"`
		Endpoint string `json:"endpoint"`
		Default  bool   `json:"default"`
	} `json:"robots"`
}

func Load(src []byte) error {
	var cfg configuration
	if err := json.Unmarshal(src, &cfg); err != nil {
		return err
	}
	for _, robot := range cfg.Robots {
		if robot.Default {
			Register(robot.Name, New(robot.Endpoint, WithDefaultRegistration()))
		} else {
			Register(robot.Name, New(robot.Endpoint))
		}
	}
	return nil
}

// DirectSend will send a message directly
func DirectSend(ctx context.Context, name, msg string) error {
	mu.RLock()
	defer mu.RUnlock()
	if r, ok := registeredRobots[name]; ok {
		return r.Send(ctx, msg)
	}
	return ErrRobotName
}
