## commands

cretea pipeline:
```bash
curl -X POST -d "path=/datas/mango-runner" -d "tag=mango" http://localhost:1234/pipeline/create
```

create pipeline same tag:
```bash
curl -X POST -d "path=/datas/mango-runner" -d "tag=mango" http://localhost:1234/pipeline/create
```

pipeline console stdout:
```bash
curl -X POST -d "path=/datas/mango-runner" -d "filename=mango_20231113_231434.txt" http://localhost:1234/pipeline/stdout
```

pipeline list:
```bash
curl -X POST -d "path=/datas/mango-runner" -d "tag=mango" http://localhost:1234/pipeline/list
```

Is service healthy?
```bash
curl -X POST http://localhost:1234/service/status
```

