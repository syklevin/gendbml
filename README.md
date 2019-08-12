GENDBML
======

store procedure models and funcs generator from c# dbml


### Build
```sh

go build

```

### Usage
```sh

./gendbml --file fixtures/SshRpt.dbml --pkg data --dir tmp --datapkg github.com/{username}/{datapkg_destination} --externalDB DBVariable --errorPkg github.com/{username}/{errorpkg_destination}

```
