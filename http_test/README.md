## commands

cretea pipeline:
```bash
curl -X POST -d "name=mango" -d "command=npm info react" http://localhost:1234/pipeline/create
```

pipeline console stdout:
```bash
curl -X POST -d "name=mango" -d "filename=mango_20240629_113957.txt" http://localhost:1234/pipeline/stdout
```

pipeline list:
```bash
curl -X POST -d "name=mango" http://localhost:1234/pipeline/list
```

Is service healthy?
```bash
curl -X POST http://localhost:1234/service/status
```

Clone a project (for Gitlab):
```bash
curl -X POST -d "name=<your_repo_name>" -d "repo=<your_repo>" -d "branch=<your_repo_branch_name>" -d "user=<your_username>" -d "pwd=<your_pwd>" http://localhost:1234/git/clone
```

