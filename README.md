### Requirements

1. Docker
2. AWS CDK and Typescript
3. AWS CLI configured with an IAM Admin Role for an AWS Account
4. Go 1.17

### Deployment

1. Run `go build` under service
2. Run `npm install` to install all dependencies
3. Run `cdk bootstrap` within `infra`
4. Run `cdk deploy` within `infra`


