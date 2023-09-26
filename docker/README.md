# Multinode elastic cluster with kibana in docker<!-- omit in toc -->

The vm.max_map_count kernel setting must be set in order for the elastic nodes in docker containers to start up.

Example error message of elastic containers failing to start elasticsearch
```
docker-es01-1  | {"@timestamp":"2023-08-31T22:10:23.254Z", "log.level":"ERROR", "message":"node validation exception\n[1] bootstrap checks failed. You must address the points described in the following [1] lines before starting Elasticsearch.\nbootstrap check failure [1] of [1]: max virtual memory areas vm.max_map_count [65530] is too low, increase to at least [262144]", "ecs.version": "1.2.0","service.name":"ES_ECS","event.dataset":"elasticsearch.server","process.thread.name":"main","log.logger":"org.elasticsearch.bootstrap.Elasticsearch","elasticsearch.node.name":"es01","elasticsearch.cluster.name":"docker-cluster"}
```

To view current vm.max settings
```bash
sudo systctl -a | grep vm.max
```

To set vm.max_map_count settings live (non-persistent)
```bash
sysctl -w vm.max_map_count=262144
```

To set vm.max_map_count persistent
```bash
sudo su
echo "vm.max_map_count=262144" >> /etc/sysctl.conf
```

- [Docker compose operations](#docker-compose-operations)
  - [To start up the cluster](#to-start-up-the-cluster)
  - [To view logs from a container/service in docker compose](#to-view-logs-from-a-containerservice-in-docker-compose)
  - [To tear down docker compose project](#to-tear-down-docker-compose-project)
  - [To just stop the docker compose project so that you can resume later](#to-just-stop-the-docker-compose-project-so-that-you-can-resume-later)
  - [To resume previously stopped docker compose project](#to-resume-previously-stopped-docker-compose-project)

## Docker compose operations

### To start up the cluster
```bash
flynshue@flynshue-Latitude-7430:~/github.com/flynshue/esctl/docker$ docker compose up -d
[+] Running 11/11
 ✔ Network docker_default         Created                                                                                                                                                                                                          0.1s 
 ✔ Volume "docker_kibanadata"     Created                                                                                                                                                                                                          0.0s 
 ✔ Volume "docker_esdata01"       Created                                                                                                                                                                                                          0.0s 
 ✔ Volume "docker_esdata02"       Created                                                                                                                                                                                                          0.0s 
 ✔ Volume "docker_certs"          Created                                                                                                                                                                                                          0.0s 
 ✔ Volume "docker_esdata03"       Created                                                                                                                                                                                                          0.0s 
 ✔ Container docker-setup-1       Started                                                                                                                                                                                                          0.8s 
 ✔ Container docker-es-data-01-1  Healthy                                                                                                                                                                                                         21.6s 
 ✔ Container docker-es-data-02-1  Healthy                                                                                                                                                                                                         22.0s 
 ✔ Container docker-es-data-03-1  Healthy                                                                                                                                                                                                         22.0s 
 ✔ Container docker-kibana-1      Started 
```

### To view logs from a container/service in docker compose

**Note: You need to target the service name that is listed in the docker compose file and not the actual container name**
```bash
$ docker compose logs kibana
docker-kibana-1  | [2023-09-12T14:24:13.074+00:00][INFO ][node] Kibana process configured with roles: [background_tasks, ui]
docker-kibana-1  | [2023-09-12T14:24:19.397+00:00][INFO ][plugins-service] Plugin "cloudChat" is disabled.
docker-kibana-1  | [2023-09-12T14:24:19.400+00:00][INFO ][plugins-service] Plugin "cloudExperiments" is disabled.
docker-kibana-1  | [2023-09-12T14:24:19.400+00:00][INFO ][plugins-service] Plugin "cloudFullStory" is disabled.
docker-kibana-1  | [2023-09-12T14:24:19.400+00:00][INFO ][plugins-service] Plugin "cloudGainsight" is disabled.
docker-kibana-1  | [2023-09-12T14:24:19.453+00:00][INFO ][plugins-service] Plugin "profiling" is disabled.
docker-kibana-1  | [2023-09-12T14:24:19.463+00:00][INFO ][plugins-service] Plugin "serverless" is disabled.
docker-kibana-1  | [2023-09-12T14:24:19.463+00:00][INFO ][plugins-service] Plugin "serverlessObservability" is disabled.
docker-kibana-1  | [2023-09-12T14:24:19.463+00:00][INFO ][plugins-service] Plugin "serverlessSearch" is disabled.
docker-kibana-1  | [2023-09-12T14:24:19.463+00:00][INFO ][plugins-service] Plugin "serverlessSecurity" is disabled.
docker-kibana-1  | [2023-09-12T14:24:19.546+00:00][INFO ][http.server.Preboot] http server running at http://0.0.0.0:5601
docker-kibana-1  | [2023-09-12T14:24:19.617+00:00][INFO ][plugins-system.preboot] Setting up [1] plugins: [interactiveSetup]
docker-kibana-1  | [2023-09-12T14:24:19.638+00:00][WARN ][config.deprecation] The default mechanism for Reporting privileges will work differently in future versions, which will affect the behavior of this cluster. Set "xpack.reporting.roles.enabled" to "false" to adopt the future behavior before upgrading.


To list all docker compose projects
```bash
$ docker compose ls -a
NAME                STATUS                  CONFIG FILES
docker              exited(1), running(4)   /home/flynshue/github.com/flynshue/esctl/docker/compose.yaml
```

### To tear down docker compose project
**Note: This will stop and remove the containers and volumes associated with the compose project**
```bash
$ docker compose down -v
[+] Running 6/6
 ✔ Container docker-kibana-1  Removed                                                                                                                                                                                                              0.6s 
 ✔ Container docker-es03-1    Removed                                                                                                                                                                                                              2.7s 
 ✔ Container docker-es02-1    Removed                                                                                                                                                                                                              2.9s 
 ✔ Container docker-es01-1    Removed                                                                                                                                                                                                             10.6s 
 ✔ Container docker-setup-1   Removed                                                                                                                                                                                                              0.2s 
 ✔ Network docker_default     Removed
```

### To just stop the docker compose project so that you can resume later
```bash
docker compose stop
```

### To resume previously stopped docker compose project
```bash
docker compose start
```
