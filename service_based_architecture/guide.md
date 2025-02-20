### At service_based_architecture folder
```commandline
cd service_based_architecture
```

### Run these files in separate terminals
```commandline
python gateway_service.py
```

```commandline
python project_service.py
```

```commandline
python user_service.py
```

```commandline
python repo_service.py
```

Run these commands in a seperate terminal other than the running servers'


### POST request presets

Add a user
```commandline
curl -X POST "http://127.0.0.1:5000/user/username=John"
```

Add a project
```commandline
curl -X POST "http://127.0.0.1:5000/project/add?name=AI_Project&type=Machine_Learning"
```

Add a repo
```commandline
curl -X POST "http://127.0.0.1:5000/repo/add?name=AI_Project"
```

Add a repo to a project
```commandline
curl -X POST "http://127.0.0.1:5000/project_repo/add?project_name=AI_Project&repo_id=1"
```

Add a repo to a user
```commandline
curl -X POST "http://127.0.0.1:5000/user_repo/add?username=John&repo_id=1"
```

### GET request presets

Get a user
```commandline
curl -X GET "http://127.0.0.1:5000/user/1"
```

Get all users
```commandline
curl -X GET "http://127.0.0.1:5000/users"
```

Get a project
```commandline
curl -X GET "http://127.0.0.1:5000/project/1"
```

Get all projects
```commandline
curl -X GET "http://127.0.0.1:5000/projects"
```

Get a repo
```commandline
curl -X GET "http://127.0.0.1:5000/repo/1"
```

Get all repos
```commandline
curl -X GET "http://127.0.0.1:5000/repos"
```