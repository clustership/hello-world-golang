# Hello world!

This repository only contains code samples to be used for basic kubernetes demo.

It is used to demo kubernetes features not for features of the app itself.

It is tailored to use s2i on OpenShift.

# Deploy on OpenShift

To deploy this application on OpenShift

```bash
export NS=hello-world # Adapt this to fit your needs
export PORT=8080
oc new-project $NS
oc new-app --name=hello-world -e MSG='Hello world in Golang !' -e PORT=${PORT} https://github.com/clustership/hello-world-golang
oc expose deployment/hello-world --port=${PORT} # Service is not automatically configured for Golang application
oc expose svc hello-world

curl $(oc get route hello-world -o jsonpath='{.spec.host}')
```

Then fine tune application configuration.

```bash
oc create sa ${NS}-sa
oc set serviceaccount deploy/hello-world ${NS}-sa
oc scale deployment/hello-world --replicas=3

cat <<EOF > hello-world-security-context-patch.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-world
  namespace: ${NS}
spec:
  template:
    spec:
      containers:
      - name: hello-world
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
          runAsNonRoot: true
          seccompProfile:
            type: RuntimeDefault
      securityContext:
        allowPrivilegeEscalation: false
        capabilities:
          drop:
          - ALL
        runAsNonRoot: true
        seccompProfile:
          type: RuntimeDefault
EOF
oc patch deploy hello-world --type="strategic" -p "$(cat hello-world-security-context-patch.yaml)" --dry-run=client -o yaml | oc replace -f -
# rm -f hello-world-security-context-patch.yaml
```
