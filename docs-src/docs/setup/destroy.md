# Destroying a Layer0 Instance

This section provides procedures for destroying (deleting) a Layer0 instance.

## Part 1: Clean Up Your Layer0 Environments
In order to destroy a Layer0 instance, you must first delete all environments in the instance.
List all environments with:
```
$ l0 environment list
```

For each environment listed in the previous step, with the exception of the environment named `api`, 
issue the following command (replacing `<environment_name>` with the name of the environment to delete):
```
l0 environment delete --wait <environment_name>
```


## Part 2: Destroy the Layer0 Instance
Once all environments have been deleted, the Layer0 instance can be deleted using the `l0-setup` tool. 
Run the following command (replacing `<instance_name>` with the name of the Layer0 instance):
```
$ l0-setup destroy <instance_name>
```

The `destroy` command is idempotent; if it fails, it is safe to re-attempt multiple times. 

!!! note
    If the  operation continues to fail, it is likely there are resources that were created outside of Layer0 that have dependencies on the resources `l0-setup` is attempting to destroy. You will need to manually remove these dependencies in order to get the `destroy` command to complete successfully. 
