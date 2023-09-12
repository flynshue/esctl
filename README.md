# esctl
CLI tool written in Go to for interacting with elasticsearch API. This tool aims to be the GO replacement of https://github.com/slmingol/escli.

## Configuration
By default, the tool will use `$USER` as the elasticsearch username.  To override that:

```bash
export ES_USERNAME=fooBar
```

To set the elasticsearch password:
```bash
export ES_PASSWORD=$(lpass show --password example.com)
```

To set the hosts:
```bash
export ES_HOSTS="https://es-data-01d.lab1.example:9200 https://es-data-01e.lab1.example:9200"
```

You can also use a configuration file instead.
```bash
vim ~/.esctl-lab.yaml 
```

```yaml
# ~/.esctl-lab.yaml 
hosts:
  - https://es-data-01d.lab1.example:9200
  - https://es-data-01e.lab1.example:9200
  - https://es-data-01f.lab1.example:9200
  - https://es-data-01g.lab1.example:9200
# skip tls verification, default is false
insecure: true
# not recommended but you can do it
username: fooBar
password: fakePass
```