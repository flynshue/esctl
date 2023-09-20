## esctl console

Send HTTP requests Elasticsearch REST API

### Synopsis

Interact with the REST APIs of Elasticsearch using http requests. This is useful for sending http requests to elasticsearch when we don't have commands built out for it yet.

```
esctl console [METHOD] [ENDPOINT] [flags]
```

### Aliases

```
esc
```

### Examples

```

# basic example
esctl console GET /my-index-000001

# command alias
esctl esc GET /my-index-000001

# without leading "/"
esctl esc GET my-index-000001

# supplying request data
esctl esc put /customer/_doc/1 -d \
'{
	"name": "John Doe"
}'

# supplying request data from file
esctl esc put /customer/_doc/2 -f /tmp/test-doc.json 
```

### Options

```
  -d, --data string       data body to be sent with http request
  -f, --filename string   file that contains data to be sent with request. --data takes precedence
  -h, --help              help for console
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl](esctl.md)	 - CLI tool for interacting with elasticsearch API

