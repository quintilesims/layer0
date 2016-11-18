# Upgrade Layer0

This section provides procedures for upgrading your Layer0 installation to the latest version.

The process of upgrading Layer0 differs depending on the version of Layer0 you are currently running, and the version you are upgrading to. Before continuing, determine which of the following sections applies to your current Layer0 installation:

* [Upgrading from version 0.5.5 or earlier to the latest version of Layer0](#upgrading-from-version-055-or-earlier-to-the-latest-version-of-layer0)
* [Upgrading from version 0.6.0 or 0.6.1 to the latest version of Layer0 0.6.2](#upgrading-from-version-060-or-061-to-the-latest-version-of-layer0)
* [Upgrading ffrom version 0.7.x to 0.7.2](#upgrading-version-07x-to-072)
* [Upgrading versions not listed above](#upgrading-versions-not-listed-above)

---

## Upgrading from version 0.5.5 or earlier to the latest version of Layer0

The architecture of Layer0 version 0.5.5 and earlier is not compatible with the architecture of more recent versions. For this reason, Layer0 versions 0.5.5 and earlier cannot be upgraded.

If you are using Layer0 0.5.5 or earlier and want to upgrade to a more recent version, we recommend that you create a new Layer0 using the current version of Layer0, and then migrate your services onto it by deploying your Docker task definitions into the new instance.

If you have questions about migrating from an older version of Layer0, visit the [\#xfra Slack channel](https://ims-dev.slack.com/messages/xfra/).

---

## Upgrading from version 0.6.0 or 0.6.1 to the latest version of Layer0
In order to upgrade from version 0.6.0 or 0.6.1 to the current version of Layer0, you must first complete some steps that are specific to these versions of Layer0.

This upgrade will move your Layer0 API server from Elastic Beanstalk to the EC2 Container Server (ECS), which will change the URL of your Layer0 API endpoint. After completing the procedures in this section, you will need to update any references to the API endpoint in your projects.

**To update to the latest version:**

!!! note
  Do not delete the **l0** or **l0-setup** files from version 0.6.0 or 0.6.1 until step 3&mdash;you will need these files to complete this upgrade process.

1. In the [Downloads section of the home page](/index.html#download), select the appropriate installation file for your operating system. Extract the zip file to a directory on your computer, but do not remove or replace the files from version 0.6.0 or 0.6.1.
2. Type the following command to destroy the ElasticBeanstalk resources, replacing _prefix_ with the name of your Layer0 prefix:
<ul>
  <li class="command">**l0-setup terraform** _prefix_ **destroy -target aws\_elastic\_beanstalk\_application.api**</li>
</ul>
3. Type the following command to backup the state files for your Layer0:
<ul>
  <li class="command">**l0-setup backup** _prefix_</li>
</ul>
4. From the [Downloads section of the home page](/index.html#download), select the appropriate installation file for your operating system. Extract the resulting zip file to your computer, and then replace the **l0** and **l0-setup** files that are already in your system path with the versions that you just downloaded.
5. Type the following command to determine which version of the **l0-setup** application you are running:
<ul>
  <li class="command">**l0-setup --version**</li>
</ul>
The output of this command should equal the version number of the current version of Layer0.<br style="line-height:3em;" />
If you see the correct version number, proceed to the next step. Otherwise, ensure that you have replaced the **l0** and **l0-setup** files in your system path with the new versions of these files.

6. Type the following command to restore the state files:
<ul>
  <li class="command">**l0-setup restore** _prefix_</li>
</ul>
7. Type the following command to update the API Docker Image tag to the appropriate version:
<ul>
  <li class="command">**l0-setup plan** _prefix_ **-var api\_docker\_image\_tag=**_vX.Y.Z_</li>
  <li class="command">Replace *vX.Y.Z* in the command above with the version number of the current version of Layer0.</li>
</ul>
8. Type the following command to apply the current version of Layer0:
<ul>
  <li class="command">**l0-setup apply** _prefix_</li>
</ul>

9. At this point, the environment variables already applied to your shell will no longer be valid. Type the following command to view the new environment variables and apply them to your shell:
  * (Windows PowerShell): **l0-setup endpoint --insecure --powershell** _prefix_ **| Out-String | Invoke-Expression**
  * (Linux/Mac): **eval "$(l0-setup endpoint --insecure** _prefix_**)"**

### Versions supported by these procedures
The procedures above are known to work when migrating between the versions listed in the following table:

| Existing version | New version |
|---|---|
| 0.6.0 | 0.6.2 |
| 0.6.0 | 0.6.3 |
| 0.6.1 | 0.6.2 |
| 0.6.1 | 0.6.3 |

---

## Upgrading version 0.7.x to 0.7.2
The commands **service logs** and **task logs** will not work for Services and Tasks created prior to version 0.7.2.
You will likely get an error: `AWS Error: the specified log group does not exist`.
You must re-create your Service in order for Layer0 to create the proper log group.

---

## Upgrading versions not listed above

!!! note
  If you have already completed the upgrade procedures listed in one of the sections above, you do not need to complete the procedures in this section.

**To upgrade to a new version of Layer0:**

1. In the [Downloads section of the home page](/index.html#download), select the appropriate installation file for your operating system. Extract the zip file to a directory on your computer, and then move the **l0** and **l0-setup** files to a folder in your system path, replacing any previous versions of these files.
2. Follow the instructions for creating an [Administrator Access Key](install/#part-2-create-an-access-key). You will use this access key when l0-setup prompts you for AWS credentials.
3. Type the following command to verify that you are working with the correct version of Layer0:
<ul>
  <li class="command">**l0-setup --version**</li>
</ul>
The output of this command should display the version number of the most recent version of Layer0. If it does, proceed to the next step; if not, ensure that you copied the latest versions of **l0** and **l0-setup** to the appropriate directories in your system path.
4. Type the following command to restore the state files for your Layer0, replacing *prefix* with the name of your Layer0 prefix:
<ul>
  <li class="command">**l0-setup restore** *prefix*</li>
</ul>

5. Type the following command to update your api image tag:
<ul>
  <li class="command">**l0-setup plan** *prefix* **-var api_docker_image_tag**=*version* </li>
</ul>

6. (**This setup is only required when upgrading to v0.7.0 and higher**) Type the following command to update your runner image tag:
<ul>
  <li class="command">**l0-setup plan** *prefix* **-var runner_version_tag**=*version* </li>
</ul>

7. Type the following command to update your AWS Access Key ID from step 2:
<ul>
  <li class="command">**l0-setup plan** *prefix* **-var api_access_key**=*access_key_id* </li>
</ul>

8. Type the following command to update your AWS Secret Access Key from step 2:
<ul>
  <li class="command">**l0-setup plan** *prefix* **-var api_secret_key**=*secret_key* </li>
</ul>

9. Type the following command to apply the upgrade:
<ul>
  <li class="command">**l0-setup apply** *prefix*</li>
</ul>

### Versions supported by these procedures
The procedures above are known to work when migrating between the versions listed in the following table:

| Existing version | | New version |
|---|:---:|---|
| 0.6.2 | &rarr; | > 0.6.2 |
| 0.6.3 | &rarr; | > 0.6.3 |
| 0.7.0 | &rarr; | 0.7.1 |
