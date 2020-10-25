# serviceMonitor
kubernetes health check using service


## customer kubernetes readiness health check using web domain

- step 1: binding k8s service name with web domain

- step 2: group application with the same label

- step 3: Using wather to loop service event

- step 4: send message to dingtalk

## Rules:

- 5 * Seconds do one health check.Concurrency http get request
- send every 20 * Seconds if continues' error happend

