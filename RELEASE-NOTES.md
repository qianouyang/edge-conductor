# Release Notes for Edge Conductor

This document provides Edge Conductor system requirements, release location,
issues and limitations, and legal information. To learn more, refer to the
following sections.


## Release Notes for v0.5.0

<details>
  <summary>Click for Release Notes details</summary>

## Contents v0.5.0

  * [New in this Release v0.5.0](#new-in-this-release-v050)
  * [Feature & Version List for Cluster Type v0.5.0](#feature---version-list-for-cluster-type-v050)
  * [Known Issues v0.5.0](#known-issues-v050)
  * [Related Documentation v0.5.0](#related-documentation-v050)
  * [Where to find the release v0.5.0](#where-to-find-the-release-v050)
  * [Hardware and Software Requirements v0.5.0](#hardware-and-software-requirements-v050)
  * [Legal Disclaimers](#legal-disclaimers)


## New in this Release v0.5.0

Edge Conductor Release 0.5.0 includes the following:

*  Support declarative node joining.

*  Replace "experience kit" terms by "package" or "kit"

*  Reduce the time for EC cluster build and service deployment

*  Upgrade to the latest versions on cluster  infra components:(Need a list)

*  Upgrade to the latest versions on supported services: (Need a list)

*  Enable Debian 11.3 OS support

*  Support Edge node with Ubuntu 22.04 as host OS & preempt_RT

*  integrated "cert-Management" service from upstream

*  Enable error code and indexing for the error messages

*  Add 3 user tutorial documents: "Hello Cluster Helm", "Web-info" and "RKE deployment"

## Feature & Version List for Cluster Type v0.5.0

The following table lists available features and versions for each cluster type
provided in the Edge Conductor v0.5.0 release.

| No.| EC v0.5.0 Features            | Version | Kind | RKE  | CAPI |
|----| -----------------------------|---------|------|------|----- |
| 1  |`Calico`                 |  3.23.1    | No  | No | Yes   |
| 2  |`CRI-O`                       |  1.23.3 | No   | No   | Yes  |
| 3  |`Edge Software Provisioner`   |  2.0.3  | No   | Yes  | Yes† |
| 4  |`Grafana Dashboard`           |  8.3.6  | Yes  | Yes  | Yes  |
| 5  |`Intel-GPU-Plugin`            |  0.24   | No   | Yes  | No   |
| 6  |`Multus`                      |  0.3.1   | Yes  | Yes  | Yes  |
| 7  |`Offline deployment`          |  N/A    | No   | Yes  | Yes†   |
| 8  |`Portainer-ce`                |  1.0.32 | Yes  | Yes  | Yes  |
| 9  |`Prometheus`                  |  36.2.1 | Yes  | Yes  | Yes  |
| 10 |`Rook-ceph`                   |  1.9.6  | No   | Yes  | Yes  |
| 11 |`Rook-ceph-cluster`           |  1.9.6  | No   | Yes  | Yes  |
| 12 |`Intel-sriov-network`                      |  3.5.1    | No   | Yes  | Yes  |
| 13 |`Akri`                        |  0.8.4  | Yes  | Yes  | Yes  |
| 14 |`NFD`                         |  0.11.0 | No   | Yes  | Yes  |
| 15 |`rt-linux-detection`          |  N/A    | Yes  | Yes  | Yes  |
| 16 |`nginix-ingress`          |  4.2.0    | No  | No  | Yes  |
| 17 | `CPU Manager`                 |  N/A    | Yes  | Yes  | No   |
| 18 |`kubevirt-operator`          |  0.54.0    | Yes  | Yes  | Yes  |
| 19 |`kubevirt-cr`          |  0.54.1    | Yes  | Yes  | Yes  |

†  Note that Edge Software Provisioner and offline deployment feature is enabled on CAPI clusters when using the BYOH deployment framework.

The following table lists the versions for each cluster type provided in the Edge Conductor v0.5.0 release.
| No.| Cluster            | Version | K8S version | CRI Version  | Other Version | Comments|
|----| ------------|---------|---------|-----------|---------------|------------|
| 1  |`KIND`               |  0.14.0    | 1.24.2  | -- | -- |-- |
| 2  |`RKE`                 |  1.3.12    | 1.23.7-rancher1-1(default)  | -- | -- |-- |
| 3  |`CAPI(core)`      |  1.2.0 | ---   | ---   | Cert-manager v1.9.0 | -- |
| 4  |`CAPI-BYOH`    |  1.2.0 | 1.23.5   |  Crio v1.23.3 Containerd v1.6.6 | -- |---|
| 5  |`CAPI-METAL3`|  1.2.0 | 1.23.5   |  Crio v1.23.5 Containerd v1.6.6 | -- |---|


## Known Issues v0.5.0
**EPJ-2755**

Description: The system time of RKE node is inconsistent with the time of day-0 node.

Workaround: No workaround, it will be fixed in next release.
  

**EPJ-2007**

Description: 

Workaround: 
**EPJ-2126**

Description: When using ESP as the OS provider to profile the OS of the
nodes, the user name and password is not set from the Package Kit
config file, but set from the ESP profile.

Workaround: Use the user name and password set in the official ESP profile to
config the node list in a Package Kit config.

**EPJ-2376**

Description: When using CAPI deploy cluster with crio as container runtime, there is error message "overlayfs: unrecognized mount option "volatile" or missing value" on the provisioned nodes. This message is by design of the upstream project. Refer to the [known issue](https://github.com/cri-o/cri-o/issues/4773) of CRI-O.

Workaround: No workaround required.

**EPJ-2849**

Description: failed to pull image "10.10.10.1:9000/docker.io/library/rt-linux-detection:latest on CAPI cluster

Workaround: No workaround, it will be fixed in next release.

**EPJ-2436**

Description: Edge Conductor, version 0.5.0 does not include the latest functional and security updates. Customers should update to the latest version as it becomes available.

| Image                                                                                                                       | #Critical | #High | #Medium | #Low | Action                                                                |
|-----------------------------------------------------------------------------------------------------------------------------|-----------|-------|---------|------|-----------------------------------------------------------------------|
| quay.io/cephcsi/cephcsi:v3.5.1                                                                                              | 0         | 18    | 103     | 80   | No action                                                             |
| rook/ceph:v1.8.5                                                                                                            |           |       |         |      | See [rook-ceph-188](docs/assets/vulnerabilities.md#rook-ceph-188)     |
| k8s.gcr.io/ingress-nginx/controller:v0.47.0                                                                                 | 0         | 8     | 4       | 1    | See [nginx-0.14.0](docs/assets/vulnerabilities.md#nginx-0.14.0)       |
| docker.io/nfvpe/multus:stable                                                                                               | 0         | 4     | 243     | 288  | No action                                                             |
| ghcr.io/k8snetworkplumbingwg/multus-cni:stable                                                                              | 0         | 4     | 243     | 288  | No action                                                             |
| quay.io/ceph/ceph:v16.2.7                                                                                                   |           |       |         |      | See [ceph-ceph-v17.1](docs/assets/vulnerabilities.md#ceph-ceph-v17.1) |
| k8s.gcr.io/build-image/debian-base:buster-v1.7.2                                                                            | 0         | 1     | 2       | 45   | No action                                                             |
| k8s.gcr.io/ingress-nginx/controller:v1.1.1 <br>@sha256:0bc88eb15f9e7f84e8e56c14fa5735aaa48<br>8b840983f87bd79b1054190e660de | 0         | 1     | 0       | 0    | No action                                                             |
| docker.io/calico/node:v3.22.1                                                                                               | 0         | 0     | 16      | 7    | No action                                                             |
| quay.io/kubevirt/virt-api:v0.41.0                                                                                           | 0         | 0     | 1       | 26   | No action                                                             |
| quay.io/kubevirt/virt-controller:v0.41.0                                                                                    | 0         | 0     | 1       | 26   | No action                                                             |
| quay.io/kubevirt/virt-operator:v0.41.0                                                                                      | 0         | 0     | 1       | 26   | No action                                                             |
| grafana/grafana:8.3.6                                                                                                       | 0         | 0     | 0       | 1    | No action                                                             |
| quay.io/kiwigrid/k8s-sidecar:1.15.6                                                                                         | 0         | 0     | 0       | 1    | No action                                                             |


**EPJ-1965**

Description: Prometheus Alertmanager exposes a cluster service while implementing a web application on port 80, which does not use TLS.

**EPJ-1973**

Description: Prometheus server enables web user interface as a cluster service on port 80 without TLS authentication.

## Related Documentation v0.5.0

The [README file](README.md)
has an overview of the Edge Conductor tool and its capabilities.

Use the [Get Started](docs/guides/get-started.md) guide to try out Edge
Conductor.

Next, follow the [Edge Conductor Tutorials](docs/tutorials/index.md) to learn
how to deploy some simple applications on the Kubernetes cluster you built with
the [Get Started](docs/guides/get-started.md) guide.

See the [Guides index](/docs/guides/index.md)
for a current list of Edge Conductor user guides.


## Where to find the release v0.5.0

Download the package from [Release Tag for v0.5.0](https://github.com/intel/edge-conductor/releases/tag/v0.5.0).


## Hardware and Software Requirements v0.5.0

Be sure your host meets the following requirements.

Hardware:

*   CPU: 2+ cores
*   Memory: 2+ GB

OS and System:

*   Ubuntu 18.04+ LTS
*   DockerCE
    * 20.10.3+ (for DockerCE v20)
      

</details>

# Release Notes for Edge Conductor

This document provides Edge Conductor system requirements, release location,
issues and limitations, and legal information. To learn more, refer to the
following sections.


## Release Notes for v0.4.0

<details>
  <summary>Click for Release Notes details</summary>

## Contents v0.4.0

  * [New in this Release v0.4.0](#new-in-this-release-v040)
  * [Feature & Version List for Cluster Type v0.4.0](#feature---version-list-for-cluster-type-v040)
  * [Known Issues v0.4.0](#known-issues-v040)
  * [Related Documentation v0.4.0](#related-documentation-v040)
  * [Where to find the release v0.4.0](#where-to-find-the-release-v040)
  * [Hardware and Software Requirements v0.4.0](#hardware-and-software-requirements-v040)
  * [Legal Disclaimers](#legal-disclaimers)


## New in this Release v0.4.0

Edge Conductor Release 0.4.0 includes the following:

*  DEK (Development Experience Kit) framework and DEK instances created for CAPI (K8s Cluster API), RKE, and KIND clusters.

*  DCE (Distributed Command Executor) framework enabled for operation automation.

*  Offline deployment solution for KIND DEK, RKE DEK (only cluster deployment for RKE) and CAPI DEK.

*  [Intel ESP](https://github.com/intel/Edge-Software-Provisioner) integrated for OS provisioning.

*  [BYOH (Bring Your Own Host)](https://github.com/vmware-tanzu/cluster-api-provider-bringyourownhost) as CAPI provider integrated with ESP.

*  SR-IOV NIC supported end-2-end software stack integrated for container workloads. (Support for VMs is planned for a future release.)

*  CRI-O as an optional supported CRI, in addition to container.

*  CPU manager, Prometheus stack and Grafana are integrated.

*  Rook operator and Ceph as the integrated backend storage solution.

*  TLS enabled for Orchestrated service endpoints.

*  Akri supported both USB and ONVIF camera.

*  NFD integrated for node feature discovery.

*  [Edge Conductor Tutorials](docs/tutorials/index.md) for first-time users.

*  RT Linux Detection on node.


## Feature & Version List for Cluster Type v0.4.0


The following table lists available features and versions for each cluster type
provided in the Edge Conductor v0.4.0 release.



| No.| EC v0.4.0 Feature            | Version | Kind | RKE  | CAPI |
|----| -----------------------------|---------|------|------|----- |
| 1  |`CPU Manager`                 |  N/A    | Yes  | Yes  | No   |
| 2  |`CRI-O`                       |  1.23.2 | No   | No   | Yes  |
| 3  |`Edge Software Provisioner`   |  2.0.3  | No   | Yes  | Yes† |
| 4  |`Grafana Dashboard`           |  8.3.6  | Yes  | Yes  | Yes  |
| 5  |`Intel-GPU-Plugin`            |  0.24  | No   | Yes  | No   |
| 6  |`Multus`                      |  v3.8   | Yes  | Yes  | Yes  |
| 7  |`Offline deployment`          |  N/A    | No   | Yes  | No   |
| 8  |`Portainer-ce`                |  1.0.22 | Yes  | Yes  | Yes  |
| 9  |`Prometheus`                  |  2.33.4 | Yes  | Yes  | Yes  |
| 10 |`Rook-ceph`                   |  1.8.5  | No   | Yes  | Yes  |
| 11 |`Rook-ceph-cluster`           |  1.8.5  | No   | Yes  | Yes  |
| 12 |`SR-IOV`                      |  N/A    | No   | Yes  | Yes  |
| 13 |`Akri`                        |  0.8.4  | Yes  | Yes  | Yes  |
| 14 |`NFD`                         |  0.11.0 | No   | Yes  | Yes  |
| 15 |`rt-linux-detection`          |  N/A    | Yes  | Yes  | Yes  |


†  Note that Edge Software Provisioner can be run on CAPI clusters when using
   the BYOH deployment framework.

## Known Issues v0.4.0

**EPJ-2133**

Description: It is not possible to modify the node list in the Experience Kit, after executing `conductor init`.

Workaround: Add all nodes in the Experience Kit configuration file before executing `conductor init`.

**EPJ-2126**

Description: When using ESP as the OS provider to profile the OS of the
nodes, the user name and password is not set from the Experience Kit
config file, but set from the ESP profile.

Workaround: Use the user name and password set in the official ESP profile to
config the node list in a Experience Kit config.

**EPJ-2397**

Description: For capi-metal3, in some cases, the work node provisioning fails, and the command "kubectl get machine -n metal3" shows that the worker node status is always in provisioning state.

Workaround: Run the command "kubectl delete machine -n metal3" to delete this machine and wait for capi-metal3 to provision the worker node again.

**EPJ-2376**

Description: When using CAPI deploy cluster with crio as container runtime, there is error message "overlayfs: unrecognized mount option "volatile" or missing value" on the provisioned nodes. This message is by design of the upstream project. Refer to the [known issue](https://github.com/cri-o/cri-o/issues/4773) of CRI-O.

Workaround: No workaround required.

**EPJ-2155**

Description: The NO_PROXY variable existing on the OS deployed node is not the same as that of the NO_PROXY variable set in EK.

Workaround: Add all private networks in OS deployed node to "global_settings" -> "no_proxy" in EK yml.

**EPJ-2849**

Description: failed to pull image "10.10.10.1:9000/docker.io/library/rt-linux-detection:latest on CAPI cluster

Workaround: No workaround, it will be fixed in next release.

**EPJ-2432**

Description: On ClusterAPI clusters, Prometheus helm deployment failed during the first and second deployments.

Workaround: There is no workaround.

**EPJ-2436**

Description: Edge Conductor, version 0.4.0 does not include the latest functional and security updates. Customers should update to the latest version as it becomes available.

| Image                                                                                                                       | #Critical | #High | #Medium | #Low | Action                                                                |
|-----------------------------------------------------------------------------------------------------------------------------|-----------|-------|---------|------|-----------------------------------------------------------------------|
| quay.io/cephcsi/cephcsi:v3.5.1                                                                                              | 0         | 18    | 103     | 80   | No action                                                             |
| rook/ceph:v1.8.5                                                                                                            |           |       |         |      | See [rook-ceph-188](docs/assets/vulnerabilities.md#rook-ceph-188)     |
| k8s.gcr.io/ingress-nginx/controller:v0.47.0                                                                                 | 0         | 8     | 4       | 1    | See [nginx-0.14.0](docs/assets/vulnerabilities.md#nginx-0.14.0)       |
| docker.io/nfvpe/multus:stable                                                                                               | 0         | 4     | 243     | 288  | No action                                                             |
| ghcr.io/k8snetworkplumbingwg/multus-cni:stable                                                                              | 0         | 4     | 243     | 288  | No action                                                             |
| quay.io/ceph/ceph:v16.2.7                                                                                                   |           |       |         |      | See [ceph-ceph-v17.1](docs/assets/vulnerabilities.md#ceph-ceph-v17.1) |
| k8s.gcr.io/build-image/debian-base:buster-v1.7.2                                                                            | 0         | 1     | 2       | 45   | No action                                                             |
| k8s.gcr.io/ingress-nginx/controller:v1.1.1 <br>@sha256:0bc88eb15f9e7f84e8e56c14fa5735aaa48<br>8b840983f87bd79b1054190e660de | 0         | 1     | 0       | 0    | No action                                                             |
| docker.io/calico/node:v3.22.1                                                                                               | 0         | 0     | 16      | 7    | No action                                                             |
| quay.io/kubevirt/virt-api:v0.41.0                                                                                           | 0         | 0     | 1       | 26   | No action                                                             |
| quay.io/kubevirt/virt-controller:v0.41.0                                                                                    | 0         | 0     | 1       | 26   | No action                                                             |
| quay.io/kubevirt/virt-operator:v0.41.0                                                                                      | 0         | 0     | 1       | 26   | No action                                                             |
| grafana/grafana:8.3.6                                                                                                       | 0         | 0     | 0       | 1    | No action                                                             |
| quay.io/kiwigrid/k8s-sidecar:1.15.6                                                                                         | 0         | 0     | 0       | 1    | No action                                                             |


**EPJ-1965**

Description: Prometheus Alertmanager exposes a cluster service while implementing a web application on port 80, which does not use TLS.

**EPJ-1973**

Description: Prometheus server enables web user interface as a cluster service on port 80 without TLS authentication.

**EPJ-2403**

Description: When deploying rook ceph on CAPI cluster, there is a possibility that pod ceph-objectstore might fail to launch. 

Workaround: No workaround

## Related Documentation v0.4.0

The [README file](README.md)
has an overview of the Edge Conductor tool and its capabilities.

Use the [Get Started](docs/guides/get-started.md) guide to try out Edge
Conductor.

Next, follow the [Edge Conductor Tutorials](docs/tutorials/index.md) to learn
how to deploy some simple applications on the Kubernetes cluster you built with
the [Get Started](docs/guides/get-started.md) guide.

See the [Guides index](/docs/guides/index.md)
for a current list of Edge Conductor user guides.


## Where to find the release v0.4.0

Download the package from [Release Tag for v0.4.0](https://github.com/intel/edge-conductor/releases/tag/v0.4.0).


## Hardware and Software Requirements v0.4.0

Be sure your host meets the following requirements.

Hardware:

*   CPU: 2+ cores
*   Memory: 2+ GB

OS and System:

*   Ubuntu 18.04+ LTS
*   DockerCE
    * 20.10.3+ (for DockerCE v20)
      

</details>

## Legal Disclaimers

Copyright (C) 2022 Intel Corporation

SPDX-License-Identifier: Apache-2.0

