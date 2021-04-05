# Cosigned

A Kubernetes admission controller to verify images have been signed by `cosign`!

![intro](images/demo.gif)

## Installation

### Prereqs

* install [ko](https://github.com/google/ko)
* install [cert-manager](https://cert-manager.io/docs/installation/kubernetes/)
* install [kustomize](https://kustomize.io/)
* install [cosign](https://github.com/sigstore/cosign)

### Install

Run `make deploy`!

## Usage

`cosigned` only watches namespaces with the label `cosigned=true` on them, so set that up:

```shell
NS=default
kubectl label ns $NS cosigned=true --overwrite
```

Grab a container and try to run it:

```shell
$ IMG=$KO_DOCKER_REPO/demo
$ crane cp --platform=linux/amd64 ubuntu $IMG
$ kubectl run -it unsigned --image=$IMG
Error from server (invalid signatures): admission webhook "cosigned.sigstore.dev" denied the request: invalid signatures
```

Sign a container:

```
$ cosign generate-key-pair
$ cosign sign -key cosign.key $IMG
Enter password for private key:
Pushing signature to: gcr.io/dlorenc-vmtest2/cosigned:sha256-fb607a5a85c963d8efe8f07b5935861aea06748f2a740617f672c6f75a35552e.cosign
```

Upload the key:

```
$ kubectl create configmap cosigned-config -n cosigned-system --dry-run -o=yaml --from-file=keys=cosign.pub | kubectl apply -f -
```

Now run it:

```shell
$ kubectl run -it signed --image=$IMG
If you don't see a command prompt, try pressing enter.
/ # 
```

## Configuration

Cosigned uses a single configmap for configuration right now.
There is one field called `keys`, which contains a concatenated list of PKIX-formatted public keys to trust.
All images must be signed by one of these keys to run in the cluster.

You can create and update this with a command like this:

```
$ kubectl create configmap cosigned-config -n cosigned-system --dry-run -o=yaml --from-file=keys=cosign.pub | kubectl apply -f -
```

Enforcement is opt-in at the namespace-level.
Namespaces with the label `cosigned=true` will be enforced.