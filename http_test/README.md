## commands

cretea pipeline:
```bash
curl -X POST -d "path=/datas/mango-runner" http://localhost:1234/pipeline/create
```

create pipeline same tag:
```bash
curl -X POST -d "path=/datas/mango-runner" -d "tag=0c539af420a54f60a55f6d7a0c4be1ec" http://localhost:1234/pipeline/create
```

pipeline status:
```bash
curl -X POST -d "path=/datas/mango-runner" -d "tag=0c539af420a54f60a55f6d7a0c4be1ec" http://localhost:1234/pipeline/status
```

pipeline console stdout:
```bash
curl -X POST -d "path=/datas/mango-runner" -d "filename=0c539af420a54f60a55f6d7a0c4be1ec_20231017_235756.txt" http://localhost:1234/pipeline/stdout
```

pipeline lsit:
```bash
curl -X POST -d "path=/datas/mango-runner" -d "tag=0c539af420a54f60a55f6d7a0c4be1ec" http://localhost:1234/pipeline/list
```


