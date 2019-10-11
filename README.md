# sless

Serverless HTTP server with git ops


It is possible to serve go standard http server handler from other git repos:

Example config:

```yaml
repos:
  serverFunctions:
    example:
      sourceURL: https://github.com/bgokden/sless-example.git
      SubFolder: example
      Path: /example
      username: bgokden
      password: ******
    example2:
      sourceURL: https://github.com/bgokden/sless-example.git
      SubFolder: example2
      Path: /example2
      username: bgokden
      password: ******
```

For password you can use github token

Save this example as config.yaml after changing your username and password

And run:

```shell
$ mkdir tempdata
$ go run main.go serve --config ./config.yaml
```

```
$ curl localhost:8080/example
server example from github
$ curl localhost:8080/example2
server example from github 2
```

You can change config and reload to see the effects

```shell
$ go run main.go serve --config ./config.yaml
```
