# RoboCall v1

### Whats needed?
- HTTP Server for serving the WebUI
- Running Asterisk as a child

### Building & Running
To build the image, run:
```shell
docker build -t asterisk -f Dockerfile.asterisk .
```

To run the image, run:
```shell
docker run -it -p 8084:8080 -v ./web/tmpl:/opt/robocall/web/tmpl -v ./cnf:/opt/robocall/cnf asterisk
```