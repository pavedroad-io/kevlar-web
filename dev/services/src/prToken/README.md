[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_kevlar-web&metric=alert_status)](https://sonarcloud.io/dashboard?id=pavedroad-io_kevlar-web)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_kevlar-web&metric=bugs)](https://sonarcloud.io/dashboard?id=pavedroad-io_kevlar-web)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_kevlar-web&metric=code_smells)](https://sonarcloud.io/dashboard?id=pavedroad-io_kevlar-web)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_kevlar-web&metric=coverage)](https://sonarcloud.io/dashboard?id=pavedroad-io_kevlar-web)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_kevlar-web&metric=duplicated_lines_density)](https://sonarcloud.io/dashboard?id=pavedroad-io_kevlar-web)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_kevlar-web&metric=ncloc)](https://sonarcloud.io/dashboard?id=pavedroad-io_kevlar-web)
[![Reliability Rating](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_kevlar-web&metric=reliability_rating)](https://sonarcloud.io/dashboard?id=pavedroad-io_kevlar-web)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_kevlar-web&metric=security_rating)](https://sonarcloud.io/dashboard?id=pavedroad-io_kevlar-web)
[![Technical Debt](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_kevlar-web&metric=sqale_index)](https://sonarcloud.io/dashboard?id=pavedroad-io_kevlar-web)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=pavedroad-io_kevlar-web&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=pavedroad-io_kevlar-web)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fpavedroad-io%2Fkevlar-web.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fpavedroad-io%2Fkevlar-web?ref=badge_shield)

# prToken
Stores user token in a cockroachDB.  These tokens are used to manage access to remote sites like GitHub.

## API
OpenAPI integration coming

## Commands and files

### make
- make                Test, compile, and deploy microservice
- make check          Run all tests, code coverage, and linters
- make generate       Enforce dependencies, generate dependency graph
- make compile        Compile the binary and docker container
- make show-coverage  Open a browser window showing code highlighed by coverage
- make start          Start a copy of the microservice locally
- make stop           Stop local microservice
- make clean          Remove executables and run go-clean
- make install        Install microserice in /usr/local/bin
- make kompose        Generate Kubernetes manifests
- make fmt            Enforce go coding standard
- make simplify       Simplify go code 

### skaffold integration
Durring the `make compile` phase of a build several things happen:

- After the go binary is built, a copy is made in the local directory.  The local binary is necessary for building the docker image.  Docker expects all files to be in the current build context.  The build context is the directory pass to the docker deaemon.  In this case that is "." for the current working directory.
- It then issues executes the `skaffold run` command.  This command builds the docker image, tags it with the current git commit-id, executes any specified container-tests, and then pushes the deployment to the local microk8s cluster.  The requires enabling the microk8s repository service with microk8s.enable repository.

Sample output:
```
skaffold run
WARN[0000] config version (skaffold/v1beta9) out of date: upgrading to latest (skaffold/v1beta11)
Generating tags...
 - localhost:32000/prtoken -> localhost:32000/prtoken:8b1064b-dirty
Tags generated in 11.022957ms
Starting build...
Building [localhost:32000/prtoken]...
Sending build context to Docker daemon   14.8MB
Step 1/11 : FROM golang:latest
 ---> 9fe4cdc1f173
Step 2/11 : LABEL "vendor": "PavedRoad.io"       "microservice": "prToken"       "description": "Stores OAUTH access tokens"       "version": "0.0.1"       "env": "dev"
 ---> Using cache
 ---> 71daff69908f
Step 3/11 : MAINTAINER "support@pavedroad.io"
 ---> Using cache
 ---> 5d30e621e516
Step 4/11 : ENV ms prToken
 ---> Using cache
 ---> fdfa5c9ea833
Step 5/11 : ENV kevlar /kevlar
 ---> Using cache
 ---> 647b5d1d300b
Step 6/11 : ENV kevlarbin $kevlar/$ms
 ---> Using cache
 ---> ab273d45fae4
Step 7/11 : RUN mkdir ${kevlar}
 ---> Using cache
 ---> 90f18ba0f5e8
Step 8/11 : WORKDIR ${kevlar}
 ---> Using cache
 ---> 85644e19edda
Step 9/11 : COPY $ms $kevlar
 ---> Using cache
 ---> 7d62287c7b07
Step 10/11 : EXPOSE 8081
 ---> Using cache
 ---> 6fc1d9521224
Step 11/11 : CMD ["/bin/sh", "-c", "$kevlarbin"]
 ---> Using cache
 ---> f50a616e1f70
Successfully built f50a616e1f70
Successfully tagged localhost:32000/prtoken:8b1064b-dirty
The push refers to repository [localhost:32000/prtoken]
8bb9d43faf58: Preparing
6ca26f3cae44: Preparing
39ba5a88a3c4: Preparing
56df4f4c91ea: Preparing
510e5f32af35: Preparing
2c8d31157b81: Preparing
7b76d801397d: Preparing
f32868cde90b: Preparing
0db06dff9d9a: Preparing
2c8d31157b81: Waiting
7b76d801397d: Waiting
f32868cde90b: Waiting
0db06dff9d9a: Waiting
56df4f4c91ea: Layer already exists
510e5f32af35: Layer already exists
8bb9d43faf58: Layer already exists
6ca26f3cae44: Layer already exists
39ba5a88a3c4: Layer already exists
2c8d31157b81: Layer already exists
7b76d801397d: Layer already exists
f32868cde90b: Layer already exists
0db06dff9d9a: Layer already exists
8b1064b-dirty: digest: sha256:358438715d27deec4962dc38f02721c1a86eae0563507e027435742335a2ff8b size: 2214
Build complete in 624.922056ms
Starting test...
Test complete in 7.569µs
Starting deploy...
kubectl client version: 1.14
deployment.extensions/prtoken configured
service/prtoken configured
persistentvolumeclaim/roach-ui-claim0 configured
deployment.extensions/roach-ui configured
service/roach-ui configured
Deploy complete in 2.144968707s
You can also run [skaffold run --tail] to get the logs
```

### Handy commands to know:

##### Find the kubernetes dashboard endpoint
To find the IP address for your local microk8s cluster issue the following kubectl command.  Then use the IP address, 10.152.183.68 in this case, to access the Kubernetes dashboard in your browser.  For `example, https://10.152.183.68`.
```
microk8s.kubectl --namespace=kube-system describe service kubernetes-dashboard
Name:              kubernetes-dashboard
Namespace:         kube-system
Labels:            k8s-app=kubernetes-dashboard
Annotations:       kubectl.kubernetes.io/last-applied-configuration:
                     {"apiVersion":"v1","kind":"Service","metadata":{"annotations":{},"labels":{"k8s-app":"kubernetes-dashboard"},"name":"kubernetes-dashboard"...
Selector:          k8s-app=kubernetes-dashboard
Type:              ClusterIP
IP:                10.152.183.68
Port:              <unset>  443/TCP
TargetPort:        8443/TCP
Endpoints:         10.1.1.32:8443
Session Affinity:  None
Events:            <none>
```

##### skaffold
Use `skaffold run` to build, tag, test, and deploy the current microservice to your local cluster.
Use `skaffold run --tail` to display logs
Use `skaffold delete` to remove an prior deployed version of your microservice.
Use `skaffold build` to generate a docker image.

Sample of skaffold delete output:
```
skaffold delete
WARN[0000] config version (skaffold/v1beta9) out of date: upgrading to latest (skaffold/v1beta11)
Cleaning up...
deployment.extensions "prtoken" deleted
service "prtoken" deleted
persistentvolumeclaim "roach-ui-claim0" deleted
deployment.extensions "roach-ui" deleted
service "roach-ui" deleted
Cleanup complete in 3.936636125s
```

##### docker/microk8s repository
The microk8s repository runs on port 32000.  To see the current list of deployed images use the following command:
```
docker images localhost:32000/prtoken
REPOSITORY                TAG                 IMAGE ID            CREATED             SIZE
localhost:32000/prtoken   04b5e67-dirty       f50a616e1f70        2 hours ago         783MB
```
If you are using a VM, change localhost to the IP address of your VM.

### Docker files
- docker-compose.yaml; used with docker-compose up to start all microservices
- docker-db-only.yaml; start cockroachdb only for executing tests

### Manifests generation
(kcompse)[http://kompose.io/] Create Kubernetes services and deployments deployment YMAL

```
INFO Kubernetes file "prtoken-service.yaml" created 
INFO Kubernetes file "roach-ui-service.yaml" created 
INFO Kubernetes file "prtoken-deployment.yaml" created 
INFO Kubernetes file "roach-ui-deployment.yaml" created 
INFO Kubernetes file "roach-ui-claim0-persistentvolumeclaim.yaml" created 
*PV resources still needs to be provided for cockroach to start*
```

### Documentation
- doc/microserice
    - prToken.png dependecy graph
    - prTokenAPI.html Open API specification

## CI Status
- TODO: github release
- TODO: DOCKER pull requets
- TODO: Go Report Card
- TODO: Slack
- TODO: Twitter following
