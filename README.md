# Mango Runner

## Todo

- [x] [061501] Provide a http call for cloning from the Git repo/branch. 
- [ ] [062101] Refactor methods basis on Git Clone.

## Build

```bash
$ go build
```

## Use Guider

Add it the content below to your `[project-root]/.mango/mango-ci.yaml`.
```yaml
Version: "abc"
Stages:
  - start
  - build

job-dev:
  stage: start
  scripts:
    - echo "dev success"

build-job:
  stage: build
  scripts:
    - echo "build success"
```

Execute at the command line.
```bash
$ mango-cli serve start 1234
```

### How to Test

[Http Test](./http_test/README.md)