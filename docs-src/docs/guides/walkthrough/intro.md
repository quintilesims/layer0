# An Iterative Walkthrough

This guide aims to take you through two increasingly-complex deployment examples using Layer0.
Successive sections build upon the previous ones, and each deployment can be completed either through the Layer0 CLI directly, or through Terraform using our custom [Layer0 Terraform Provider](/reference/terraform-plugin).

We assume that you're using Layer0 v0.9.0 or later.
If you have not already installed and configured Layer0, see the [installation guide](/setup/install).
If you are running an older version of Layer0, you may need to [upgrade](/setup/upgrade#upgrading-older-versions-of-layer0).

If you intend to deploy services using the Layer0 Terraform Provider, you'll want to make sure that you've [installed](/reference/terraform-plugin/#install) the provider correctly.

Regardless of the deployment method you choose, we maintain a [guides repository](https://github.com/quintilesims/guides/) that you should clone/download.
It contains all the files you will need to progress through this walkthrough.
As you do so, we will assume that your working directory matches the part of the guide that you're following (for example, Deployment 1 of this guide will assume that your working directory is `.../walkthrough/deployment-1/`).

**Table of Contents**:

- [Deployment 1](deployment-1): Deploying a web service (Guestbook)
- [Deployment 2](deployment-2): Deploying Guestbook and a data store service (Redis)


---

