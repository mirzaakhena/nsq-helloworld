
For build the image 
```
./build.sh
```

For see the image
```
docker images
```

For delete the image
```
docker image rm <the_image_name>
```

For run the 4 container (mirserverapi, mirservercon, nsqd, nsqlookupd)
```
docker-compose up
```

For run the 4 container with daemon mode
```
docker-compose up -d
```

For testing the apps
```
Open browser type
localhost:8080/test?message=Helloworld
The message will received by mirserverapi (see the printed log) then published to NSQ
The message later will forward by NSQ to mirservercon (see the printed log)
```

Sample output
```
mirserverapi_1  | {"func":"main.(*MBProducer).CallProducer","level":"info","msg":"raw data {\"message\":\"Helloooow\"}","time":"0331 023749.649"}
mirservercon_1  | {"func":"main.(*MBConsumer).TestRequest","level":"info","msg":"Receive message: {\"message\":\"Helloooow\"}","time":"0331 023749.661"}
```

For lookup the running container
```
docker ps
```

For stop container
```
docker stop <the_container_name_from_ps>
```

For remove container
```
docker rm <the_container_name_from_ps>
```

For ssh into container
```
docker exec -it <the_container_name_from_ps> sh
```

