[Edge Conductor]: https://github.com/intel/edge-conductor
[Tutorials]: ../index.md
[Sample Applications]: ./index.md
[Hello Cluster!]: ./hello-cluster.md
[Hello Cluster! Helm]: ./hello-cluster-helm.md
[Web Indexing Application]: ./web-indexing.md
[Get Started]: ../../guides/get-started.md

[Edge Conductor] / [Tutorials] / [Sample Applications] / [Hello Cluster! Helm]

# Hello Cluster! Helm Version

In this tutorial, you will re-package the [Hello Cluster!] application as a
Helm chart and learn how to deploy it using Helm.



## Contents

*   [Prerequisites](#prerequisites)
*   [Build a Docker Image](#build-a-docker-image)
*   [Package the Application](#package-the-application)
*   [Deploy the Application](#deploy-the-application)
*   [Clean Up](#clean-up)
*   [What's Next](#whats-next)



## Prerequisites

You must have completed the [Get Started] steps before
you can run this tutorial.

1.  Install Helm to package the Hello Cluster! application as a helm chart.
Follow the [Installing Helm](https://helm.sh/docs/intro/install/) steps to
install the Helm CLI.


## Build a Docker Image

First, you need to get a Docker image that will respond with "Hello Cluster!"
when queried. In this tutorial, you will build a simple image for the pod. You
can also find a similar image on [Docker Hub](https://registry.hub.docker.com/).

### Create a directory

From the `_workspace` directory, run the command:

```bash
mkdir hellocluster && cd hellocluster
```
### Create a "Hello Cluster!" webpage

Run the following command to create a .js that, when queried, will respond with
"Hello Cluster!":

```bash

cat >server.js <<EOF
    var http = require('http');
    var handleRequest = function(request, response) {
        console.log('Received request for URL: ' + request.url);
        response.writeHead(200);
        response.end('Hello Cluster!');
    };
    var www = http.createServer(handleRequest);
    www.listen(8080);
EOF
```


### Create a Dockerfile

Run the following command to create a Dockerfile for the Docker image that
you will build in the next step:

```bash
cat >Dockerfile<<EOF
    FROM node:6.9.2
    EXPOSE 8080
    COPY server.js .
    CMD node server.js
EOF
```

### Build a Docker image

Run the following command to create a Docker image:

```bash
docker build -t hello-cluster:v1 .
```

To check the image, run the command:

```bash
docker images | grep hello-cluster
```

If the image is successfully created, you will see output similar to the
following:

```bash
hello-cluster    v1    017c026c43a3   15 seconds ago   655MB
```


### Login to the local Harbor repository

The Edge Conductor tool will generate a local registry when it is initialized.

Log in to the local registry with the following command, using your host IP
address in place of `<nnn.nnn.nnn.nnn>` and the username and password configured
in the [Prepare Custom
Config](../../guides/get-started.md#prepare-custom-config) step of Get Started.

```bash
docker login <nnn.nnn.nnn.nnn>:9000
```


### Tag and push this image to Harbor repository

Run the following commands to tag and push your image to the Harbor repository:

```bash
docker tag hello-cluster:v1 <nnn.nnn.nnn.nnn>:9000/library/hello-cluster:v1
docker push <nnn.nnn.nnn.nnn>:9000/library/hello-cluster:v1
```

where:

*  `<nnn.nnn.nnn.nnn>` is your host IP address.

*  `9000` is a goharbor/nginx-photon port that is created after completing the
   [Deploy KIND Cluster](../../guides/get-started.md#deploy-a-kind-cluster) step
   of Get Started.

*  The Harbor username and password are the same ones you created in the
   [Prepare Custom Config](../../guides/get-started.md#prepare-custom-config)
   step of Get Started.

For example:

```bash
docker tag hello-cluster:v1 10.67.106.156:9000/library/hello-cluster:v1
docker push 10.67.106.156:9000/library/hello-cluster:v1
```

To check that the image was pushed, run the command:

```bash
docker pull 10.67.106.156:9000/library/hello-cluster:v1
```


## Package the Application

In this step, you will create a simple chart called `hello-cluster-helm`.

### Create a Helm chart template

From the `_workspace` directory, run the following command to create a new chart
with the given name:

```bash
helm create hello-cluster-helm
```

After the command finishes, there is a chart in `./hello-cluster-helm`. You can
edit it and create your own templates.

You can also follow the Helm [Chart Development Guide](https://helm.sh/docs/topics/charts/)
to develop your own charts.


### Change to the Helm chart directory

```bash
cd hello-cluster-helm
```

The file structure generated under the `hello-cluster-helm` folder is:

```bash
├── charts
├── Chart.yaml
├── templates
│   ├── deployment.yaml
│   ├── _helpers.tpl
│   ├── hpa.yaml
│   ├── ingress.yaml
│   ├── NOTES.txt
│   ├── serviceaccount.yaml
│   ├── service.yaml
│   └── tests
│       └── test-connection.yaml
└── values.yaml
```

### Modify the chart template


Helm uses Go templates for templating your resource files.

The `templates/` directory is for template files. When Helm evaluates a chart,
it will send all of the files in the `templates/` directory through the template
rendering engine. It then collects the results of those templates and sends them
on to Kubernetes.

The `values.yaml` file is also important to templates. This file contains the
default values for a chart. These values may be overridden by users during
Helm install or Helm upgrade.

The `Chart.yaml` file contains a description of the chart. You can access it
from within a template. The `charts/` directory may contain other charts (which
are called subcharts). Later in this guide, you will see how those work when it
comes to template rendering.

This [Helm Getting Started](https://helm.sh/docs/chart_template_guide/getting_started/)
guide will show you more about how to create a chart.

In this section of the tutorial, you will modify 3 files:

*  ./templates/deployment.yaml
*  ./templates/service.yaml
*  ./values.yaml

First, you need to modify the template slightly using your preferred text
editor, using the following file as an example.

Replace the `./templates/deployment.yaml` file with the following content:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "hello-cluster-helm.fullname" . }}
  labels:
    {{- include "hello-cluster-helm.labels" . | nindent 4 }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "hello-cluster-helm.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      {{- with .Values.podAnnotations }}
      annotations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      labels:
        {{- include "hello-cluster-helm.selectorLabels" . | nindent 8 }}
    spec:
      {{- with .Values.imagePullSecrets }}
      imagePullSecrets:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      serviceAccountName: {{ include "hello-cluster-helm.serviceAccountName" . }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}
      containers:
        - name: {{ .Chart.Name }}
          securityContext:
            {{- toYaml .Values.securityContext | nindent 12 }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: 8080
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /
              port: http
          readinessProbe:
            httpGet:
              path: /
              port: http
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}

```

Second, replace the `./templates/service.yaml` file with the following content:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: {{ include "hello-cluster-helm.fullname" . }}
  labels:
    {{- include "hello-cluster-helm.labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: {{ .Values.service.port }}
      nodePort: {{ .Values.service.nodeport }}
      protocol: TCP
      name: http
  selector:
    {{- include "hello-cluster-helm.selectorLabels" . | nindent 4 }}

```

Third, replace the `./values.yaml` file with the following content. Replace the
text marked <nnn.nnn.nnn.nnn> with the IP of your image repository.

```yaml
# Default values for hello-cluster-helm.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 2

image:
  repository: <nnn.nnn.nnn.nnn>:9000/library/hello-cluster
  pullPolicy: IfNotPresent
  # Overrides the image tag whose default is the chart appVersion.
  tag: "v1"

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""

serviceAccount:
  # Specifies whether a service account should be created
  create: true
  # Annotations to add to the service account
  annotations: {}
  # The name of the service account to use.
  # If not set and create is true, a name is generated using the fullname template
  name: ""

podAnnotations: {}

podSecurityContext: {}
  # fsGroup: 2000

securityContext: {}
  # capabilities:
  #   drop:
  #   - ALL
  # readOnlyRootFilesystem: true
  # runAsNonRoot: true
  # runAsUser: 1000

service:
  type: NodePort
  port: 8080
  nodeport: 30003
ingress:
  enabled: false
  className: ""
  annotations: {}
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  hosts:
    - host: chart-example.local
      paths:
        - path: /
          pathType: ImplementationSpecific
  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources: {}
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #   cpu: 100m
  #   memory: 128Mi
  # requests:
  #   cpu: 100m
  #   memory: 128Mi

autoscaling:
  enabled: false
  minReplicas: 1
  maxReplicas: 100
  targetCPUUtilizationPercentage: 80
  # targetMemoryUtilizationPercentage: 80

nodeSelector: {}

tolerations: []

affinity: {}

```


### Check the chart

Run the following command to examine a chart for possible issues:

```bash
helm lint --strict ../hello-cluster-helm/
```

You will see output similar to:

```bash
==> Linting ./hello-cluster-helm/
[INFO] Chart.yaml: icon is recommended

1 chart(s) linted, 0 chart(s) failed
```

### Package the chart

Run the following command to package the chart:

```bash
cd ..
helm package hello-cluster-helm
```

You will see output similar to:

```bash
Successfully packaged chart and saved it to: <helm directory>/hello-cluster-helm-0.1.0.tgz
```

## Deploy the Application

In this step, you will learn how to deploy the application using Helm.

### Install the chart

Run the following command to install the chart:

```bash
helm install hello-cluster-helm  hello-cluster-helm-0.1.0.tgz
```

You will see output similar to:

```bash
NAME: hello-cluster-helm
LAST DEPLOYED: Wed Dec 22 09:52:34 2021
NAMESPACE: default
STATUS: deployed
REVISION: 1
NOTES:
1. Get the application URL by running these commands:
  export NODE_PORT=$(kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services hello-cluster-helm)
  export NODE_IP=$(kubectl get nodes --namespace default -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT

```

In this example, `hello-cluster-helm` is your release name.


### Check the release


Run the following command to list your releases:

```bash
helm list
```

You will see output similar to:

```bash
NAME                    NAMESPACE       REVISION        UPDATED                                 STATUS          CHART                           APP VERSION
hello-cluster-helm      default         1               2021-12-22 09:52:34.894836393 +0000 UTC deployed        hello-cluster-helm-0.1.0        1.16.0
```

Run the following commands to check the deployments and services deployed to
the Helm chart.

```bash
kubectl get deployments.apps
kubectl describe services hello-cluster-helm
```

You will see output similar to:

```bash
NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
hello-cluster-helm   2/2     2            2           3m57s

Name:                     hello-cluster-helm
Namespace:                default
Labels:                   app.kubernetes.io/instance=hello-cluster-helm
                          app.kubernetes.io/managed-by=Helm
                          app.kubernetes.io/name=hello-cluster-helm
                          app.kubernetes.io/version=1.16.0
                          helm.sh/chart=hello-cluster-helm-0.1.0
Annotations:              meta.helm.sh/release-name: hello-cluster-helm
                          meta.helm.sh/release-namespace: default
Selector:                 app.kubernetes.io/instance=hello-cluster-helm,app.kubernetes.io/name=hello-cluster-helm
Type:                     NodePort
IP Family Policy:         SingleStack
IP Families:              IPv4
IP:                       10.96.167.18
IPs:                      10.96.167.18
Port:                     http  8080/TCP
TargetPort:               8080/TCP
NodePort:                 http  30003/TCP
Endpoints:                10.244.2.8:8080,10.244.3.10:8080
Session Affinity:         None
External Traffic Policy:  Cluster
Events:                   <none>

```

> *NOTE:*  Make a note of the port and NodePort value for the service. For
example, in the preceding output, the port value is 8080 and the NodePort value
is 30003.


Run the following command to list the pods that are running the Hello Cluster
application:

```bash
kubectl get pods --output=wide
```

You will see output similar to:

```bash
NAME                                  READY   STATUS    RESTARTS   AGE     IP            NODE           NOMINATED NODE   READINESS GATES
hello-cluster-helm-847fdfb79f-lwk77   1/1     Running   0          6m47s   10.244.2.8    kind-worker    <none>           <none>
hello-cluster-helm-847fdfb79f-pxq8m   1/1     Running   0          6m47s   10.244.3.10   kind-worker3   <none>           <none>
```

### Access the Application

#### (For KIND) Use port-forward to access the Hello World application

Run the following command:

```bash
kubectl port-forward -n default  service/hello-cluster-helm  5999:8080
```

If your pod has a different name, change `hello-cluster-helm` in the command
above to the name of your pod.

You will see output similar to:

```bash
Forwarding from 127.0.0.1:5999 -> 8080
Forwarding from [::1]:5999 -> 8080
```

Open a new terminal and run the following command to access the
`hello-cluster-helm` application:

```bash
curl http://<localhost-ip>:<forward-port>
```

In the preceding output, the `localhost-ip` value is 127.0.0.1 and the
`forward-port` value is 5999.

So in this tutorial we can run like this:

```bash
curl http://127.0.0.1:5999
```

The response to a successful request is a hello message:

```bash
Hello Cluster!
```

## Clean Up

If you want to remove `hello-cluster-helm` from the cluster, run the following
command:

```bash
helm uninstall hello-cluster-helm
```

where:

  `hello-cluster-helm` is your chart name.

This command removes all of the resources associated with the last release of
the chart as well as the release history, freeing it up for future use.


## What's Next

Congratulations! You have deployed an application using a Helm chart.
Next you can try to deploy a web indexing service on the Kubernetes cluster.


-----\
Previous Tutorial: [Hello Cluster!]\
Next Tutorial: [Web Indexing Application]\
\
Back to: [Tutorials](/docs/tutorials/index.md)


Copyright (C) 2022 Intel Corporation

SPDX-License-Identifier: Apache-2.0
