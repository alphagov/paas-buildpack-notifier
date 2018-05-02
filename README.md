# paas-buildpack-notifier

Send notifications to Cloud Foundry users about buildpack changes.

NOTE: This is a proof-of-concept. It's not currently in use.

Fetch buildpacks from updated environment:
```
cf curl /v2/buildpacks > buildpacks.new.json
```

Build and run against production environment:
```
go build
./paas-buildpack-notifier -a https://api.cloud.service.gov.uk -t "$(cf oauth-token | awk '{print $2}')" -b buildpacks.new.json
```
