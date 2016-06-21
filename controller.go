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
)

func main() {
  var ingClient client.IngressInterface
  if kubeClient, err := client.NewInCluster(); err != nil {
    log.Fatalf("Failed to create client: %v.", err)
  } else {
    ingClient = kubeClient.Extensions().Ingress(api.NamespaceAll)
  }

  rateLimiter := flowcontrol.NewTokenBucketRateLimiter(0.1, 1)
  known := &extensions.IngressList{}

  // TODO: pull BIG-IP address and login credentials from the K8s ConfigMap

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