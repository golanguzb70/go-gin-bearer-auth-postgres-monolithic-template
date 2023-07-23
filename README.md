# go-gin-basicauth-monolithic-template
Being developed...

# Instructions to run the source code
Follow the instractions to run the code in your local machine.
MakeFile, Docker engine and docker compose should be installed in your local machine.
Make sure that you write all needed environment variables into /config/.env file 

To run the program
```
make compose_up
```
To stop the program
```
make compose_down
```

After program is started successfully, you can check if it is running using this address.
```
http://localhost:8000/v1/swagger/index.html
```

### Setting environment variables for gitlab and github actions
These environment variables can be saved different places according to your OS configurations. It can be stored in .zshrc, .bashrc, .profile files.
```
GITHUB_USERNAME="your_github_username"
GITHUB_PERSONAL_ACCESS_TOKEN="your github personal access token."
GITLAB_USERNAME="your_gitlab_username"
GITLAB_PERSONAL_ACCESS_TOKEN="your gitlab personal access token."
```