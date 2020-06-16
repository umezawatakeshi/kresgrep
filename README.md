# kresgrep -- Kubernetes Resource Grep

**kresgrep** is a tool to filter Kubernetes resource YAML file.


## Usage

```
kresgrep [options...] [filenames...]
  -k, --kind        search pattern for kind
  -a, --name        search pattern for name
  -n, --namespace   search pattern for namespace
```

kresgrep searches for resource that matches all the specified patterns.

- A filename of `-` stands for standard input.
- If no filenames are given, kresgrep will read standard input.


## Examples

- `kresgrep --name '^nginx' deployment.yaml`  
  search for resources whose name begin with `nginx`.
- `kresgrep --kind '^ConfigMap$' --namespace 'myapp' common.yaml`  
  search for ConfigMap resources whose namespace contains `myapp`.

## Caveat

- The first non-flag argument is *not* search pattern unlike ordinal grep.
- If the namespace is omitted in resource, kresgrep assumes `default` namespace, even in case of non-namespaced resource.


## Todo

- Label based search
- Ignore case match
- Fixed string match
- Inverted match


## Copyright Notice and License

Copyright &copy; 2020 UMEZAWA Takeshi

**kresgrep** is licensed under GNU General Public License (GNU GPL) version 3 or later. See  [LICENSE](LICENSE) for more information.