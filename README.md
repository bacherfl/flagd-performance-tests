###

```shell
kubectl testkube create test --file test.js --type "k6/script" --name flagd-load-test
kubectl testkube run test flagd-load-test -f
```