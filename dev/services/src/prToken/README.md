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
-- prToken.png dependecy graph
-- prTokenAPI.html Open API specification

## CI Status
- TODO: github release
- TODO: DOCKER pull requets
- TODO: Go Report Card
- TODO: Slack
- TODO: Twitter following
