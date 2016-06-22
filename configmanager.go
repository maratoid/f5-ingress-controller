package main

type ConfigManager interface {
  Set(*Config)
  Get() *Config
  Close()
}

type ChannelConfigManager struct {
  conf *Config
  get  chan *Config
  set  chan *Config
  done chan bool
}

func NewChannelConfigManager(conf *Config) *ChannelConfigManager {
  parser := &ChannelConfigManager{conf, make(chan *Config), make(chan *Config), make(chan bool)}
  parser.Start()
  return parser
}

func (self *ChannelConfigManager) Start() {
  go func() {
    defer func() {
      close(self.get)
      close(self.set)
      close(self.done)
    }()

    for {
      select {
      case self.get <- self.conf:
      case value := <-self.set:
        self.conf = value
      case <-self.done:
        return
      }
    }

  }()
  
}

func (self *ChannelConfigManager) Close() {
  self.done <- true
}

func (self *ChannelConfigManager) Set(conf *Config) {
  self.set <- conf
}

func (self *ChannelConfigManager) Get() *Config {
  return <-self.get
}