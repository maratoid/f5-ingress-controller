all: push

TAG = latest
PREFIX = gcr.io/k8s-work/f5-ingress

controller: controller.go
	CGO_ENABLED=0 GOOS=linux godep go build -a -installsuffix cgo -ldflags '-w' -o controller ./controller.go ./configmanager.go /.watcher.go

container: controller
	docker build -t $(PREFIX):$(TAG) .

push: container
	gcloud docker push $(PREFIX):$(TAG)

clean:
	rm -f controller