# DevOps with Kubernetes — Submissions

Exercise submissions for [DevOps with Kubernetes](https://courses.mooc.fi/org/uh-cs/courses/devops-with-kubernetes) (2026).

Each exercise has a GitHub release tagged with its number (e.g. `1.1`). Each link below points to the relevant app directory at that release.

## Exercises

### Chapter 2 — Kubernetes basics

- [1.1.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.1)
- [1.2.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.2)
- [1.3.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.3)
- [1.4.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.4)
- [1.5.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.5)
- [1.6.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.6)
- [1.7.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.7)
- [1.8.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.8)
- [1.9.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.9)
- [1.10.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.10)
- [1.11.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.11)
- [1.12.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.12)
- [1.13.](https://github.com/milaneuh/KubernetesSubmissions/tree/1.13)

### Chapter 2 — Networking between pods

- [2.1.](https://github.com/milaneuh/KubernetesSubmissions/tree/2.1)
- [2.2.](https://github.com/milaneuh/KubernetesSubmissions/tree/2.2)
- [2.3.](https://github.com/milaneuh/KubernetesSubmissions/tree/2.3)
- [2.4.](https://github.com/milaneuh/KubernetesSubmissions/tree/2.4)
- [2.5.](https://github.com/milaneuh/KubernetesSubmissions/tree/2.5)
- [2.6.](https://github.com/milaneuh/KubernetesSubmissions/tree/2.6)
- [2.7.](https://github.com/milaneuh/KubernetesSubmissions/tree/2.7)
- [2.8.](https://github.com/milaneuh/KubernetesSubmissions/tree/2.8)
- [2.9.](https://github.com/milaneuh/KubernetesSubmissions/tree/2.9)
- [2.10.](https://github.com/milaneuh/KubernetesSubmissions/tree/2.10)

### Chapter 4 — To the cloud

- [3.1.](https://github.com/milaneuh/KubernetesSubmissions/tree/3.1)
- [3.2.](https://github.com/milaneuh/KubernetesSubmissions/tree/3.2)
- [3.3.](https://github.com/milaneuh/KubernetesSubmissions/tree/3.3)
- [3.4.](https://github.com/milaneuh/KubernetesSubmissions/tree/3.4)
- [3.5.](https://github.com/milaneuh/KubernetesSubmissions/tree/3.5)
- [3.6.](https://github.com/milaneuh/KubernetesSubmissions/tree/3.6)
- [3.7.](https://github.com/milaneuh/KubernetesSubmissions/tree/3.7)
- [3.8.](https://github.com/milaneuh/KubernetesSubmissions/tree/3.8)
- [3.9.](https://github.com/milaneuh/KubernetesSubmissions/tree/3.9)

### DBaaS VS DIY
#### DBaaS
Database as a service (DBaaS) offers a simple database solution, where the cloud provider handles all of the infrastructure work regarding the database.
We still own the schema, query tuning, access control and data security.

*Pros:*

- We can deploy the database without managing the hardware, the OS, or the DB engine install and patches.
- HA, backups and security patches are handled by the cloud provider.
- Scaling is simplified, especially for storage and read replicas.

*Cons:*

- We are restricted to the cloud provider settings: no superuser access, limited extensions, forced maintenance windows.
- We pay a premium for the managed service, so costs get higher than DIY at large scale.
- Vendor lock-in, especially with provider-specific engines or features.

#### DIY
We install, configure and handle the database ourselves.

*Pros:*

- Full ownership, no restriction on configuration, optimization or plugins.
- Lower cost at large scale.
- More portable between cloud providers or on-premise.

*Cons:*

- Ops costs are much higher: backups, downtime, monitoring and security patches are on us.
- We need database expertise in the team.
- Higher risk of misconfiguration, since we write all of the configuration ourselves.

#### Initialization

- **DBaaS:** a few clicks or a few lines of Terraform, ready in minutes. No upfront cost, we pay for what we use.
- **DIY:** we provision the VM, install and configure the engine, the network, the users and TLS. It takes hours to days, so the initial cost is mostly engineer time.

#### Maintenance

- **DBaaS:** patches and minor upgrades are applied by the provider. Monitoring is built in.
- **DIY:** we apply patches, plan upgrades and downtime, run our own monitoring, and someone has to be on call.

#### Backups

- **DBaaS:** automatic snapshots and point-in-time recovery (PITR) out of the box. A restore is one click.
- **DIY:** we set up the tooling ourselves
