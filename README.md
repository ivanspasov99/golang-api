
[Task Definition](task.md)

| Table of Contents                                     |
|:------------------------------------------------------|
| * [API Docs](#api-docs)                               |
 | * [Full Software Lifecycle](#full-software-lifecycle) |
| * [Architecture](#architecture)                       |
| * [Package](#package)                                 |
| * [CI](#ci)                                           |
| * [CD](#cd)                                           |
| * [Monitor](#monitor)                                 |
| * [Job Handler](#job-handler)                         |
| * [Testing](#testing)                                 |
| * [Logging Package](#logging-package)                 |
| * [Graph Algorithm](#graph-algorithm)                 |
| * [Security](#security)                               | 

## API Docs
<details>
<summary>
<code>POST</code>
<code><b>/job?mode={mode}</b></code>
<code>Accepts job with tasks and returns ordered commands as different format depending on the `mode` passed as query parameter</code>
</summary>

##### Query

| name | type     | data type | description                                                | default |
|------|----------|-----------|------------------------------------------------------------|---------|
| mode | optional | string    | represents required response format - JSON, Bash supported | JSON    |

##### Responses

| http code | Content-Type       | Request                                  | Response                                   |
|-----------|--------------------|------------------------------------------|--------------------------------------------|
| `200`     | `application/json` | [Example Request](#example-json-request) | [Example Response](#example-json-response) | 
| `200`     | `text`             | [Example Request](#example-bash-request) | [Example Response](#example-bash-response) |

###### Example JSON Request
```curl -d @testing/input.json http://localhost:8080```

###### Example JSON Response
```json
[
  {
    "name":"task-1",
    "command":"touch /tmp/file1"
  },
  {
    "name":"task-3",
    "command":"echo 'Hello World!' > /tmp/file1"
  },
  {
    "name":"task-2",
    "command":"cat /tmp/file1"
  },
  {
    "name":"task-4",
    "command":"rm /tmp/file1"
  }
]
```

###### Example Bash Request
```curl -d @testing/input.json http://localhost:8080?mode=bash | bash```

###### Example Bash Response
```bash
#!/usr/bin/env bash
touch /tmp/file1
echo "Hello World!" > /tmp/file1
cat /tmp/file1
rm /tmp/file1
```

</details>

## Full Software Lifecycle 
What should be added to be production ready.

API Docs could be added using Swagger

### Architecture
Diagrams should be added
- Block Diagram
- Flow Diagrams

### Package
- All Configurations will be passed as `env` variables through charts. Example how should they be handled in [Config](pkg/config/config.go). Also, there is example for some cluster communication client 
- Using Helm & Kubernetes

### CI 
Executed on Every PR/Merge
- Build - Using Docker
- Integration - Jenkins or Azure. Execution of Unit, Performance testing
- Generate Version and Release

Executed on Merge
- Push Image to Container Registry (GCP)
- Push Helm Chart to Artifact Registry (GCP)
- Security, Compliant checks on built image

### CD
Depends on whether Continuous Delivery/Continues Deployment/Progressive Delivery is set and what Service Mesh is used
- Configure Environment Configuration Repository (keeping all microservices version and used for release/promotion)
- Deploy with Argo CD 
- Configure [Cluster Bootstrapping](https://argo-cd.readthedocs.io/en/stable/operator-manual/cluster-bootstrapping/) if multiple clusters are used for different environments
- Secrets Management - [External Secrets Operator](https://external-secrets.io/v0.7.1/)
- Centralized Secret Management Platform - [HashiCorp Vault](https://www.vaultproject.io/)

### Monitor 
Tools that have specified in [Job Handler](#job-handler)

## Job Handler
Job Handler handles every job in separate goroutine, it is highly possible real big data scenario so 
performance is crucial. Therefore, it is taken into account and time complexity is linear - Graph Implementation O(n + e).

- Encoding/Decoding special symbols use-cases are not taken into account
- Job Processing is separated to two middlewares using chain of responsibility pattern - job.Handle and job.HandleError as both we grow in the future so they should be separated as abstractions
- Monitoring/Alerting is out of scope. Could be done with different tools depending on requirements
  - Sentry - Error Alerting, could alert the DoD (developer on duty) for errors which should be process immediately 
  - Kibana - Logging Analyse tool
  - Prometheus - Resource/Performance analyse tool 

## Testing
Testing is created using Table Driven Testing over BDT (Behavior Driven Testing). Output could be improved when test fails, as it would 
bring big value in debugging faster. 

## Logging Package
Package encapsulate productive json requirement logging which is required by a lot of analysing log tools

Logging package could be extended with dynamic logging and log level state which represent the option to change the level of logging (debug, warn, info, error)
This help in generating fewer logs when not needed and set more logs when problem arise for debugging purposes
This could be implemented through a configmap in the which is deployed in k8s cluster for examples separately from
the application, then you can consume/read it as `env` variable in the code  

## Graph Algorithm
**It is better to use already implemented packages which are community adopted and tested**, but I have decided to refresh my skills a little bit

**It is best to be implemented using generics as now it is very limited to one type/struct etc**

**Time Complexity - O(n + e)**

The implementation is using maps which does no guarantee order of the topological sorting. It can be 
implemented with arrays or can be improved with following custom map key logic
```go
// key will be used as map key
type key struct {
	name string
}

// graph will look like
type graph struct {
	vertices map[key]*Vertex
	edges    map[key]*Edge
}

// comparison function
func (k key) Less(other key) bool {
	// Return true if k should be sorted before the other.
	// Return false otherwise.
}
```

## Security
Security not part of the task