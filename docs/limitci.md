# LimitCI

limitci is use for add concurrency for pipelinerun

## Configuring a `limitci`

```yaml
apiVersion: ci.w6d.io/v1alpha1
kind: LimitCi
metadata:
  name: limitci-sample
spec:
  concurrent: 2
```