# cms

[![Test Status](http://ci.dmalone.io/api/v1/teams/spearwind/pipelines/cms-service/jobs/test-app/badge)](http://ci.dmalone.io/teams/spearwind/pipelines/cms-service)

[API Docs](http://docs.spearwind.apiary.io/#)


## Run locally with defaults

`go build && ./cms`

Optionally, use [fresh](https://github.com/pilu/fresh) to auto-reload changes to speed up your dev cycles.

## System Configuration

The system is configured via environment variables. These are the available environment variables used to configure this system:

1. AWS_ENDPOINT - the Amazon SES email endpoint. i.e. https://email.us-east-1.amazonaws.com/
1. AWS_ACCESS_KEY_ID - your AWS Access Key ID, with SES rights
1. AWS_SECRET_ACCESS_KEY - your AWS Secret Access Key, with SES rights
1. EMAIL_TEMPLATE_DIR - the location of the directory containing all of the email templates
1. FB_APP_ID - Facebook Application ID, for use with Facebook Login
1. FB_APP_SECRET  - Facebook Application Secret, for use with Facebook Login
1. MONGO_URL - Mongo DB Connection URL; e.g. mongodb://127.0.0.1:27017/cms-admin


## Develop

This project uses [Glide](https://github.com/Masterminds/glide)

To setup your local workspace, first clone this project, and then run `glide install`

To run the project test suite, run `go test $(glide novendor)`

## Pipelines

Concourse pipelines are hosted at ci.dmalone.io under a team named Spearwind, with Github auth setup to allow anyone associated with the spear-wind org in Github to use this team. To target and login to Concourse:

    fly -t dmalone login -c http://ci.dmalone.io:8080
    fly login -n spearwind -t dmalone -c http://ci.dmalone.io

### Modify the pipeline and update in Concourse

    fly -t dmalone sp -p cms-sevice -c ci/pipeline.yml

### Testing individual tasks using fly

Params for individual tasks will be read from your local environment variables. To run a task that requires params:

```
cms $ fly -t spearwind execute -c \
  ./ci/test.yml \
  --input cms=.
```

Note that the params defined in the tasks yml file must be *blank* in order for them to be replaced via the environment variables set for your environment

## Build Docker Image

`GOOS=linux GOARCH=amd64 go build -ldflags "-X main.VERSION=1.0" && mkdir -p build/linux64 && mv cms build/linux64`

`docker build -t dmalone/cms-admin .`
