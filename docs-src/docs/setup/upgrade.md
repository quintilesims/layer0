# Upgrade Layer0

This section provides procedures for upgrading your Layer0 installation to the latest version.

1. In the [Downloads section of the home page](/index.html#download), select the appropriate installation file for your operating system. Extract the zip file to a directory on your computer, and then move the **l0** and **l0-setup** files to a folder in your system path, replacing any previous versions of these files.

2. You will need your existing access keys when l0-setup prompts you for AWS credentials. Type the following commands to find them:

    - `l0-setup terraform [prefix] output access_key`
    - `l0-setup terraform [prefix] output secret_key`

3. Type the following command to verify that you are working with the correct version of Layer0:

    - `l0-setup --version`

    The output of this command should display the version number of the most recent version of Layer0. If it does, proceed to the next step; if not, ensure that you copied the latest versions of **l0** and **l0-setup** to the appropriate directories in your system path.

4. Type the following command to restore the state files for your Layer0, replacing `[prefix]` with the name of your Layer0 prefix:

    - `l0-setup restore [prefix]`

5. Type the following command to update your api image tag:

    - `l0-setup plan [prefix] -var api_docker_image_tag=[version]`

6. Type the following command to update your runner image tag:

    - `l0-setup plan [prefix] -var runner_version_tag=[version]`

7. Type the following command to apply the upgrade:

    - `l0-setup apply [prefix]`
