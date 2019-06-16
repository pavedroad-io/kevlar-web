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
