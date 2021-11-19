# Cosigned

** THE CODE IN THIS REPO HAS BEEN MOVED TO THE OFFICIAL [COSIGN REPO](github.com/sigstore/cosign) **
** THIS IS ARCHIVED **

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

> Don't forget to change Go module name <br/>
> **module github.com/dlorenc/cosigned --> module github.com/<your_github_name>/cosigned**

```shell
$ export SECRET_KEY_REF=k8s://default/mysecret
$ envsubst \
    < config/manager/kustomization.template.yaml \
    > config/manager/kustomization.yaml
$ export PROJECT_ID=$(gcloud config get-value project)
$ export KO_DOCKER_REPO=gcr.io/$PROJECT_ID
$ export GITHUB_NAME="dlorenc"
$ IMG=ko://github.com/$GITHUB_NAME/cosigned make deploy
```

## Usage

`cosigned` only watches namespaces with the label `cosigned=true` on them, so set that up:

```shell
NS=default
kubectl label ns $NS cosigned=true --overwrite
```

Grab a container and try to run it:

```shell
$ IMAGE=$KO_DOCKER_REPO/demo
$ crane cp --platform=linux/amd64 ubuntu $IMAGE
$ kubectl run -it unsigned --image=$IMAGE
Error from server (invalid signatures): admission webhook "cosigned.sigstore.dev" denied the request: invalid signatures
```

Sign a container:

```
$ cosign generate-key-pair $SECRET_KEY_REF
$ cosign sign -key $SECRET_KEY_REF $IMAGE
Enter password for private key:
Pushing signature to: gcr.io/dlorenc-vmtest2/cosigned:sha256-fb607a5a85c963d8efe8f07b5935861aea06748f2a740617f672c6f75a35552e.cosign
```

Now run it:

```shell
$ kubectl run -it signed --image=$IMAGE
If you don't see a command prompt, try pressing enter.
/ # 
```

## Configuration

Cosigned uses a single Secret for configuration right now. Because `cosign` now supports to store pub/private key pair in Kubernetes secrets.
There is one field called `cosign.pub`, which contains a PKIX-formatted public key to trust.
All images must be signed by the key to run in the cluster.

Enforcement is opt-in at the namespace-level.
Namespaces with the label `cosigned=true` will be enforced.
