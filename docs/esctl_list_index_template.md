## esctl list index template

get one or more index templates

```
esctl list index template [template name pattern] [flags]
```

### Aliases

```
templates
```

### Examples

```
# List all index templates and their index patterns
esctl list index template

# Get list index templates that match template pattern
esctl list index template .monit*

```

### Options

```
  -h, --help     help for template
      --legacy   list only legacy index templates
```

### Options inherited from parent commands

```
      --config string   config file (default is $HOME/.esctl.yaml)
```

### SEE ALSO

* [esctl list index](esctl_list_index.md)	 - show information about one or more index

