# cms


[![wercker status](https://app.wercker.com/status/475e09b299697263c1d546fc24e9b5d7/m "wercker status")](https://app.wercker.com/project/bykey/475e09b299697263c1d546fc24e9b5d7)

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
