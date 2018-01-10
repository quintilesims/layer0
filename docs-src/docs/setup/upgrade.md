# Upgrade a Layer0 Instance

This section provides procedures for upgrading your Layer0 installation to the latest version.
This assumes you are using Layer0 version `v0.10.0` or later. 

!!! note 
    Layer0 does not support updating MAJOR or MINOR versions in place unless explicitly stated otherwise.
    Users will either need to create a new Layer0 instance and migrate to it or destroy and re-create their Layer0 instance in these circumstances.

Run the **upgrade** command, replacing `<instance_name>` and `<version>` with the name of the Layer0 instance and new version, respectively:
```
$ l0-setup upgrade <instance_name> <version>
```

This will prompt you about the updated `source` and `version` inputs changing. 
If you are not satisfied with the changes, exit the application during the prompts. 
For full control on changing inputs, please use the **set** command. 

**Example Usage**
```
$ l0-setup upgrade mylayer0 v0.10.1
This will update the 'version' input
        From: [v0.10.0]
        To:   [v0.10.1]

        Press 'enter' to accept this change:
This will update the 'source' input
        From: [github.com/quintilesims/layer0//setup/module?ref=v0.10.0]
        To:   [github.com/quintilesims/layer0//setup/module?ref=v0.10.1]

        Press 'enter' to accept this change:
        ...
        
Everything looks good! You are now ready to run 'l0-setup apply mylayer0'
```

As stated by the command output, run the **apply** command to apply the changes to the Layer0 instance.
If any errors occur, please contact the Layer0 team. 
