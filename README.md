
# Load Balancer

A Load balancer which distributes requests to many servers, implemented from scratch in go lang.

## Features

- `Static load Balacing`: Acheived using Round robin assumes optimistic requests completion.

- `Periodic Health Checks`: Hits health check end points after a specific timeout can be tuned according to use case.

- `Active Passive Server Management`: Requests wont be redirected to failing servers or servers that are shutdown.




## API Reference

### Loadbalancer

```http
GET /
```

Redirects the request to server according to algorithmn, copies the response and sends back to client. 
### Servers

#### Default endpoint

```http
  GET /
```

Returns the message from servers. For easier identification I explicitly specified Server Id.

#### Healthcheck endpoint

```http
  GET /healthcheck
```

Returns health message with 200 status code.


## Run Locally

Clone the project

```bash
  git clone https://github.com/adityadafe/load-balancer
```

Go to the server directory

```bash
  cd server
```

Run server

```bash
  go run .
```

Go  to loadbal directory

```bash
  cd loadbal
```

Run loadbalancer

```bash
  go run .
```




## Benefits

-  **Improved Availability**: Distributes incoming traffic across multiple servers to ensure no single server becomes a bottleneck, enhancing overall system uptime.
-  **Scalability**: Easily scales application infrastructure by adding or removing servers as demand changes, ensuring consistent performance.
- **Enhanced Performance**: Balances the load to prevent any single server from being overwhelmed, improving response times and throughput.
- **Fault Tolerance**: Automatically detects server failures and reroutes traffic to healthy servers, maintaining service continuity.
- **Efficient Resource Utilization**: Maximizes the utilization of available server resources, reducing idle time and optimizing cost-efficiency.
- **Security**: Provides an additional layer of security by masking the backend server infrastructure, and can also help mitigate DDoS attacks.
- **Simplified Maintenance**: Allows for maintenance or upgrades on individual servers without disrupting the overall service, ensuring continuous availability.
- **Geographic Distribution**: Directs traffic to servers located closest to the user, reducing latency and improving user experience.


## Caution

This is not recommended to put in production enviornment use battle tested lb like [fabio](https://github.com/fabiolb/fabio)
## Future Work

I have created branches with different versions and implementation. In coming time i will add registry endpoint for service registry.