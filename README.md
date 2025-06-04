# builder-playground-operator
The Builder Playground Operator is a Kubernetes Operator designed to simplify the deployment of multi-container applications built using the Builder Playground CLI. This operator supports:
1. Declarative Deployment Automation – Apply a single Custom Resource (CR) manifest to automatically generate and manage a StatefulSet.
2. Standalone YAML Manifest Generation – Use the CLI tool to generate Kubernetes manifests without needing to run the operator.

## Description
1. Operator Mode
Deploy your application by applying a BuilderPlaygroundDeployment Custom Resource.
The Operator automatically translates the **CR** into a Kubernetes StatefulSet, deploy and manage the resources.
Supports:
    - Multiple containers with custom ports, env vars, volume mounts
    - Init containers
    - Persistent storage using hostPath

2. CLI Mode
    - Use "builder-playground" cli to generate all the configuration files
    - Use "builder-playground-operator" cli to generate Kubernetes manifest YAML files
    - Users then can deploy the StatefulSet.yaml file manually

## Getting Started

### Prerequisites
- go version v1.23.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### Use it as a BuilderPlayground K8S **StatefulSet Generator** for local K8S deployment testing
```
# Generate all the configuration files using builder-playground
sudo ./builder-playground cook l1 \
  --latest-fork \
  --use-reth-for-validation \
  --output ~/my-builder-testnet-01 \
  --genesis-delay 15 \
  --log-level debug \
  --dry-run
```
You can see all the necessary files and configurations are generated as below,
```
$ cd ~/my-builder-testnet-01
$ ll

drwxr-xr-x  4 root   root      4096 May 22 03:04 data_validator/
-rw-r--r--  1 root   root    518877 May 22 03:04 genesis.json
-rw-r--r--  1 root   root       538 May 22 03:04 graph.dot
-rw-r--r--  1 root   root        64 May 22 03:04 jwtsecret
-rw-r--r--  1 root   root   9337132 May 22 03:04 l2-genesis.json
-rw-r--r--  1 root   root      5508 May 22 03:04 manifest.json
-rw-r--r--  1 root   root      1133 May 22 03:04 rollup.json
drwxr-xr-x  2 root   root      4096 May 22 03:04 scripts/
drwxr-xr-x  2 root   root      4096 May 22 03:04 testnet/
```

Then builder-playground-operator will use the file **"manifest.json"** to generate the StatefulSet and CR files,
- `builder-config-dir` in command below is the directory that the **builder-playground** used above to geneate the configurations. Here it is used as the input of builder-playground-operator CLI.
- `k8s-manifests-dir` in command below is the output directory that **builder-playground-operator** uses to geneate the manifests.
```
# Compile the Go source code in the current directory into a binary named **"builder-playground-operator"**
go build -o builder-playground-operator .

# Run the compiled "builder-playground-operator" binary with CLI mode and specify paths for configuration and Kubernetes manifests
./builder-playground-operator --cli \                            # Run in CLI mode
  --manifest /home/ubuntu/my-builder-testnet-01/manifest.json \  # Use the specified manifest file for deployment instructions
  --builder-config-dir /home/ubuntu/my-builder-testnet-01 \      # Directory containing builder-specific configuration files
  --k8s-manifests-dir /home/ubuntu/my-builder-operator-01        # Directory containing Kubernetes manifest templates to apply
```

There will be 3 files geneated. 
- You can use "BuilderPlaygroundStatefulSet.yaml" to deploy builder playground to your local cluster.
- The StatefulSet use the folder **"builder-config-dir"** to retrieve all the configurations and write Data
```
$ cd my-builder-operator-01
$ ll

-rw-r--r--  1 ubuntu ubuntu 5038 May 22 03:11 BuilderPlaygroundStatefulSet.yaml
-rw-r--r--  1 ubuntu ubuntu 3673 May 22 03:11 CR-BuilderPlaygroundDeployment.yaml
-rw-r--r--  1 ubuntu ubuntu 6744 May 22 03:11 processed.json

```
Once you get file "BuilderPlaygroundStatefulSet.yaml", you can deploy Builder Playground on Kubernetes manually.

### To Deploy Builder Playground on Kubernes using a real Kubernetes Operator
- Use existing BuilderPlaygroundOperator and CR to deploy Builder Playground on Kubernes (This part will be added soon)
- Build and push your own BuilderPlaygroundOperator image to the location specified by `IMG`

The part below shows show how to build and push your own BuilderPlaygroundOperator image and use it deploy Builder Playground on Kubernes
```sh
# Build the Docker image using the Dockerfile and Makefile (typically builds and tags the image as 'controller:latest')
sudo make docker-build

# Tag the built Docker image with your Docker Hub repository name and 'latest' tag
sudo docker tag controller:latest samuelest/builder-playground-operator:latest

# Push the tagged Docker image to your Docker Hub repository
sudo docker push samuelest/builder-playground-operator:latest

# Deploy the operator to the Kubernetes cluster (usually applies manifests like RBAC, CRDs, and the operator Deployment)
sudo make deploy
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

When you push the Operator to a different registry,
You need to change the file `/manager/manager.yaml` **Manager to the cluster with the image specified by `IMG`:**

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the CR and the Operator will handle all the left for StatefulSet deployment:**

```sh
# Apply a custom resource (CR) definition to the cluster, which triggers the operator to reconcile and act
kubectl apply -f /home/ubuntu/my-builder-operator-01/CR-BuilderPlaygroundDeployment-LocalPath.yaml
or
kubectl apply -f /home/ubuntu/my-builder-operator-01/CR-BuilderPlaygroundDeployment-PVC.yaml
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.





**The part below will be changed soon.**
**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/builder-playground-operator:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/builder-playground-operator/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.


**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

