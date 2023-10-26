![Test Status](https://github.com/henders/writing-an-envoy-wasm-plugin/actions/workflows/test.yml/badge.svg)

# Envoy WASM Plugin
For retrieving JWTs for service-to-service API requests. 
Follow along with the code from https://medium.com/zendesk-engineering/writing-an-istio-wasm-plugin-in-go-for-migrating-100s-of-services-to-new-auth-strategy-part-1-cd551e1455d7

### Deploying to local K8s cluster

Build the Docker Image:
```shell
$ docker buildx build . -t shender/wasmplugin:v1
$ docker push shender/wasmplugin:v1
```

Deploy to K8s:
```shell
$ kubectl apply -f k8s_deploy.yml
```

Wasm binary
`tinygo build -o main.wasm -gc=custom -tags=custommalloc -scheduler=none -target=wasi main.go`
