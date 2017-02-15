# envlate
Expand environment variables in template (the templates use Go text/template syntax)

## Installation

Download from https://github.com/tkuchiki/envlate/releases

## Usage

```shell
$ ./envlate --help
usage: envlate [<flags>]

Expand environment variables in template (the templates use Go text/template syntax)

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
  -f, --file=FILE        Template file
  -u, --unset-error      Treat unset variables as an error
  -o, --output=FILENAME  Write the output to the file rather than to stdout
      --mode="0644"      File permission
      --version          Show application version.
```

## Examples

```shell
$ cat test.tpl
- {{.foo}}
- {{.bar}}

$ ./envlate -f examples/test.tpl
-
-

$ foo=baz bar=qux ./envlate -f examples/test.tpl
- baz
- qux

$ echo "{{.foo}}" | foo=bar ./envlate
bar

$ echo '{{or .foo "baz"}}' | ./envlate
baz

$ ./envlate -f examples/test.tpl -u
2017/02/11 21:52:35 line 1 char 4 : no entry for key `foo`

$ echo "{{.foo}}" | foo=bar ./envlate -o test.txt
$ cat test.txt
bar

$ foo="foo,bar,baz" ./envlate -f examples/split.tpl
- foo
- bar
- baz
```
