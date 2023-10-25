## esctl kbn

Send HTTP requests Kibana REST API

### Synopsis

Interact with the REST APIs for Kibana using http requests.

```
esctl kbn [METHOD] [ENDPOINT] [flags]
```

### Examples

```

# basic example
esctl kbn GET /api/spaces/space

```

### Options

```
  -d, --data string       data body to be sent with http request
  -f, --filename string   file that contains data to be sent with request. --data takes precedence
  -h, --help              help for kbn
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl](esctl.md)	 - CLI tool for interacting with elasticsearch API

