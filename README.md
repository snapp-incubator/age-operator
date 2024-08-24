# Age-Operator
This repository is a secret manager for kubernetes using [age](https://github.com/FiloSottile/age).

### What is secret manager?
By default, you put all your configuration on a codebase. There may be sensitive data like passwords or tokens among them. a secret manager lets you encrypt your sensitive data and then push it to codebase.

## Description
Instead of pushing raw data to git, encrypt your data, push it to codebase, apply it to kubernetes using gitops, and kubernetes will use this operator to create a [secret](https://kubernetes.io/docs/concepts/configuration/secret/) for you. then you can use that secret for your deployments.

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

## Step-by-Step guide through using age-operator

### Age CLI Installation

You can use [this link](https://github.com/FiloSottile/age#installation) to install `age`. Also, you can build from source if you have a supported version of golang installed on your system. It's mentioned in the given link.

### Generate Age Keys

You need to  generate `age keys` by running the code below simply:

```sh
age-keygen -o {key file name}
# For example
age-keygen -o key.txt
```

If you look inside the file, it has three lines.

```text
# created: 2022-xx-yyT00:00:00+04:30
# public key: age1fjn89y8svr9rdqh6c9h6drxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
AGE-SECRET-KEY-1ZQ3729NNAAP8MAFXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

- The first line shows the creation timestamp of the Age key.
- The second line is the public key. It's used to encrypt data.
- The third line is the secret key. It's used to create an AgeKey object that decrypts the AgeSecret object.

### Encryption with Age

For encryption, you will need a `public key` and a file that contains your sensitive data. your yaml data should be flat. nested yaml files are not supported now.

sample data.yaml:

```yaml
password: my_password
token: my_token
```

and then run:

```sh
age -r {public key} -e -a {plaintext data file} > {encrypted data file}
# for example
age -r age1fjn89y8svr9rdqh6c9h6drxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx -e -a data.yaml > data.age
```

If you look at the `data.age` file, it contains data like this (surely content differs).

```text
-----BEGIN AGE ENCRYPTED FILE-----
YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBXbXlPYUxHRTVMRGxBdkxr
Zk5VZXpSR0Npc2EzYmRUbFNUbTVRRUpvb0ZjCjk1blA2QTFmWHN5akV1aDhGQUJR
RTBwWmRvaFJjUWlYcFBBdS93bFBiaGMKLS0tIE5LZDN4aElNMEhwZXcwWW9ZUmdN
bTI0V0NGWTJkTElaRmFQNjhWREQ3bXcKKHYCUSb/xvPlj5umQRFwwd1ULlXDTYXw
jFRZvb9z4cXANc6Vp6kK8aoXNw0EzT46WId4KtTgCVwl7UDcgj+LXiO/e4J/2Rk0
0z1P3YUb
-----END AGE ENCRYPTED FILE-----
```

### Decryption with Age

To decrypt data, you will need a `secret key` and an encrypted file.

```sh
age -d -i {secret key file} {encrypted data file}
# For example
age -d -i key.txt data.age
```

## Kubernetes Resources

### AgeKey object

To create an `AgeKey` object, you need to create a template similar to the following.

```yaml
apiVersion: gitopssecret.snappcloud.io/v1alpha1
kind: AgeKey
metadata:
    name: {fill name}
    namespace: {fill namespace}
spec:
    ageSecretKey: {fill with secret key}
```

Here is an `AgeKey` sample:

```bash
oc get AgeKey agekey-sample -o yaml
```

```yaml
apiVersion: gitopssecret.snappcloud.io/v1alpha1
kind: AgeKey
metadata:
  name: agekey-sample
  namespace: default
spec:
  ageSecretKey: "AGE-SECRET-KEY-1ZQ3729NNAAP8MAFXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
```

### AgeSecret object

You can create an `AgeSecret` object with the below template. You can find more information about each field in the [AgeSecret Fields](./secret-management.md#agesecret) section:

```yaml
apiVersion: gitopssecret.snappcloud.io/v1alpha1
kind: AgeSecret
metadata:
    name: {fill name}
    namespace: {fill namespace}
    labels:
        ... {fill labels}
    annotations:
        ... {fill annotations}
spec:
    labelsToRemove:
      - {label1}
      - {label2}
    suspend: {fill suspend}
    ageKeyRef: {fill ref name}
    stringData:
        ... {fill stringData}
```

#### AgeSecret

- <h5>name</h5> Name of the AgeSecret object. If the decryption process was successful, the controller will generate a <a href="https://kubernetes.io/docs/concepts/configuration/secret/">Secret</a> object with this name.

- <h5>namespace</h5>
  The namespace of the AgeSecret object, the secret will get created in this namespace.

- <h5>labels</h5>
  A set of "key:value" that will be <b>copied inside the generated secret</b>.

- <h5>annotations</h5>
  A set of "key:value" that will be <b>copied inside the generated secret</b>.

- <h5>labelsToRemove</h5>
  An array of labels to remove while creating the child secret, and not to inherit them. Sample use-case is inside CD on k8s, when you want a label selector to track "AgeSecret" but not the child secret.

- <h5>suspend</h5>
  It's boolean. The default value is <b>false</b>. It determines whether the controller should reconcile on changes and apply changes to secret or you are just testing and the controller should not change anything.

- <h5>ageKeyRef</h5>
  This field indicates the name of the AgeKey object that you have encrypted this AgeSecret with it, so it has to be the same as the "AgeKey metadata.name" field's value.

- <h5>stringData</h5>
  This field is where you should put the encrypted Age message without its "-----BEGIN AGE ENCRYPTED FILE-----" prefix and "-----END AGE ENCRYPTED FILE-----" suffix.

> [!NOTE]
> Be careful not to push files containing "secret key" or "plain configuration" on git during the encryption and decryption steps.

Here is an `AgeSecret` sample:

```bash
oc get AgeSecret agesecret-sample -o yaml
```

```yaml
apiVersion: gitopssecret.snappcloud.io/v1alpha1
kind: AgeSecret
metadata:
  name: agesecret-sample
  namespace: test-age-secret
  labels:
    key_label: value_label
  annotations:
    key_annotation: value_annotation
spec:
  suspend: false
  ageKeyRef: agekey-sample
  stringData: |
    YWdlLWVuY3J5cHRpb24ub3JnL3YxCi0+IFgyNTUxOSBXbXlPYUxHRTVMRGxBdkxr
    Zk5VZXpSR0Npc2EzYmRUbFNUbTVRRUpvb0ZjCjk1blA2QTFmWHN5akV1aDhGQUJR
    RTBwWmRvaFJjUWlYcFBBdS93bFBiaGMKLS0tIE5LZDN4aElNMEhwZXcwWW9ZUmdN
    bTI0V0NGWTJkTElaRmFQNjhWREQ3bXcKKHYCUSb/xvPlj5umQRFwwd1ULlXDTYXw
    jFRZvb9z4cXANc6Vp6kK8aoXNw0EzT46WId4KtTgCVwl7UDcgj+LXiO/e4J/2Rk0
    0z1P3YUb
```

## Installation

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:
	
```sh
make docker-build docker-push IMG=<some-registry>/gitops-secret-manager:tag
```
	
3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/gitops-secret-manager:tag
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing
After forking this repository, add your code and add tests for your code. then make sure that test cases are alright.

```sh
make test

# to test api with verbosity
make test-api

# to test controller with verbosity
make test-controller
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/) 
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the cluster 

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

And in case you don't want to enable webhook in local, run command below and comment [WEBHOOK] sections in config/crd.

```sh
make run ENABLE_WEBHOOKS=false
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

