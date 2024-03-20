## Go open telemetry

- Project of postgraduate course Go Expert at Full Cycle.

### How to run

- First, you need create a .env file in the root of the project following the .env.example file.

```bash

# This command will start the services and run service-a and service-b
docker-compose up -d
```

### How to test

- You can use the file in ./tests/api.http to test the services.

- To access jaeger, go to: http://localhost:16686
- To access zipkin, go to: http://localhost:9411
- To access prometheus, go to: http://localhost:9090
- To access grafana, go to: http://localhost:3000
