# lichess-bot

[Lichess compensate for lag](https://lichess.org/lag). Nevertheless we run in
`eu-central-1` - Frankfurt - as that's the geographically closest AWS region to
the Lichess servers, which are currently hosted by OVH in Strasbourg.

```
aws cloudformation deploy `
    --region eu-central-1 `
    --stack-name lichess-bot `
    --template-file .\infrastructure.yaml
```

## Publishing

You'll need both [Docker][docker] and the [AWS CLI][awscli] installed.

[docker]: <https://docs.docker.com/get-docker/>
[awscli]: <https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2.html>

 1. `$AWS_REGION = "eu-central-1"`
 2. `$AWS_ACCT_ID = (aws sts get-caller-identity --query Account --output text)`
 3. `docker build . -f .\cmd\lichess-bot\Dockerfile -t lichess-bot -t "$AWS_ACCT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/lichess-bot"`
 4. `aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin "$AWS_ACCT_ID.dkr.ecr.$AWS_REGION.amazonaws.com"`
 5. `docker push "$AWS_ACCT_ID.dkr.ecr.$AWS_REGION.amazonaws.com/lichess-bot"`
