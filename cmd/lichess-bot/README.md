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
