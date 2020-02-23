# Sif Docker Swarm Example

![Logo](https://raw.githubusercontent.com/go-sif/sif/master/media/logo-128.png)

This example covers the deployment of a simple [Sif](github.com/go-sif/sif) job to a [Docker Swarm](https://docs.docker.com/engine/swarm/).

## Building the Example

To build this example, you will need Go version 1.13 and Docker installed.

```bash
# Downloads dependencies and test data
$ make dependencies
# Compiles the example, and builds a Docker image
$ make build
```

## Running the Example

To run the example, it is recommended that you temporarily create a local swarm:

```bash
$ docker swarm init
```

To deploy the example:

```bash
$ docker stack deploy -c docker-compose.yml sif-example
```

To monitor log output from the Coordinator:

```bash
$ docker service logs -f sif-example_sif-coordinator
# You should see something similar to this:
2020/02/23 21:16:33 Starting Sif Coordinator at sif-coordinator:1643
2020/02/23 21:16:33 Waiting for 2 workers to connect...
2020/02/23 21:16:35 Registered worker 1a9fecdd-9314-4807-a54b-f66a3d09dfbc at 10.0.1.5:1643
2020/02/23 21:16:35 Registered worker 81a79f07-b8db-4fc3-b214-8327864210ef at 10.0.1.4:1643
2020/02/23 21:16:35 Running job...
2020/02/23 21:16:35 Asking worker 1a9fecdd-9314-4807-a54b-f66a3d09dfbc to run stage stage-0
2020/02/23 21:16:35 Asking worker 81a79f07-b8db-4fc3-b214-8327864210ef to run stage stage-0
2020/02/23 21:16:37 Asking worker 81a79f07-b8db-4fc3-b214-8327864210ef to run stage stage-1
2020/02/23 21:16:37 Asking worker 1a9fecdd-9314-4807-a54b-f66a3d09dfbc to run stage stage-1
2020/02/23 21:16:37 Asking worker 81a79f07-b8db-4fc3-b214-8327864210ef to supply prepared partitions to coordinator
2020/02/23 21:16:37 Asking worker 1a9fecdd-9314-4807-a54b-f66a3d09dfbc to supply prepared partitions to coordinator
2020/02/23 21:16:37 Stopping worker 1a9fecdd-9314-4807-a54b-f66a3d09dfbc...
2020/02/23 21:16:37 Stopping worker 81a79f07-b8db-4fc3-b214-8327864210ef...
2020/02/23 21:16:37 There have been 331952 system discoveries in Elite Dangerous in the last 7 days.
# Press CTRL+C to stop following the logs
```

## Cleaning Up

Remove the stack:

```bash
$ docker stack rm sif-example
```

Then, tear down the local swarm like this:

```bash
$ docker swarm leave -f
```
