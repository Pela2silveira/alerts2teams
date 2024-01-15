# alert2teams
Simple middleware between alertmanager and msteams. You can set it a template for transform alertmanager payload to teams accepted payload.

This was tested with alertmanager 0.24. In the last version exist native support for this integration.

## Install
Just set DESTINATION_ENDPOINT env vars with your webhook endpoint

for example:
```export DESTINATION_ENDPOINT=my_webhook_url```

You can set your own template, you can modify it in configmap manifest.

Deploy:
```kubectl apply -n your_namespace -f manifests/.```

In alertmanager.yaml set http://alert2teams:8080 as receiver:

```
...
# A list of notification receivers.
receivers:
    webhook_configs:
      - name: 'msteams'
         webhook_configs:
           - send_resolved: true
             url: http://alert2teams:8080
...
```