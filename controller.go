package main

import (
  "log"
  "os"
  "os/exec"
  "reflect"
  "text/template"

  "k8s.io/kubernetes/pkg/api"
  "k8s.io/kubernetes/pkg/apis/extensions"
  client "k8s.io/kubernetes/pkg/client/unversioned"
  "k8s.io/kubernetes/pkg/util/flowcontrol"
  "github.com/scottdware/go-bigip"
  "gopkg.in/yaml.v2"
)

const CONFIG_FILE = "/etc/config/bigip-config.yaml"

// Config map struct
type Config struct {
  Url string `yaml:"url"`
  Username string 'yaml:"username"'
  Password string 'yaml:"password"'
}

// simple yaml file loader
func loadConfig(configFile string) *Config {
  conf := &Config{}
  configData, err := ioutil.ReadFile(configFile)
  check(err)

  err = yaml.Unmarshal(configData, conf)
  check(err)
  return conf
}

func main() {
  var ingClient client.IngressInterface
  if kubeClient, err := client.NewInCluster(); err != nil {
    log.Fatalf("Failed to create client: %v.", err)
  } else {
    ingClient = kubeClient.Extensions().Ingress(api.NamespaceAll)
  }

  rateLimiter := flowcontrol.NewTokenBucketRateLimiter(0.1, 1)
  known := &extensions.IngressList{}

  // pull BIG-IP address and login credentials from the K8s ConfigMap file
  confManager := NewChannelConfigManager(loadConfig(CONFIG_FILE))

  // Watch the file for modification and update the config manager with the new config when it's available
  watcher, err := WatchFile(CONFIG_FILE, time.Second, func() {
    log.Printf("Configfile Updated")
    conf := loadConfig(CONFIG_FILE)
    confManager.Set(conf)
  })
  check(err)

  // Controller loop
  for {
    rateLimiter.Accept()
    ingresses, err := ingClient.List(api.ListOptions{})
    if err != nil {
      log.Printf("Error retrieving ingresses: %v", err)
      continue
    }
    if reflect.DeepEqual(ingresses.Items, known.Items) {
      continue
    }

    // Todo: interact with F5 REST api based on delta between ingresses and known

    // update known
    known = ingresses
  }
}