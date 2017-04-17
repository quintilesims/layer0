# Deployment guide: Guestbook one-off task

In this example, you will learn how to use layer0 to run a one-off task. A task is used to run a single instance of your Task Definition and is typically a short running job that will be stopped once finished.

---

## Before you start
In order to complete the procedures in this section, you must install and configure Layer0 v0.8.4 or later. If you have not already configured Layer0, see the [installation guide](/setup/install). If you are running an older version of Layer0, see the [upgrade instructions](/setup/upgrade#upgrading-older-versions-of-layer0).

## Part 1: Prepare the task definition

1. Download the [Guestbook One-off Task Definition](https://github.com/quintilesims/layer0-examples/blob/master/one-off-task/Dockerrun.aws.json) and save it to your computer as **Dockerrun.aws.json**.

## Part 2: Create a deploy
Next, you will create a new deploy for the task using the **deploy create** command. At the command prompt, run the following command:

`l0 deploy create Dockerrun.aws.json one-off-task-dpl`

You will see the following output:
```
DEPLOY ID           DEPLOY NAME        VERSION
one-off-task-dpl.1  one-off-task-dpl   1
```

## Part 3: Create the task
At this point, you can use the **task create** command to run a copy of the task.

To run the task, use the following command:

`l0 task create demo-env echo-tsk one-off-task-dpl:latest --wait`

You will see the following output:
```
TASK ID       TASK NAME         ENVIRONMENT  DEPLOY              SCALE
one-off851c9  echo-tsk          demo-env     one-off-task-dpl:1  0/1 (1)
```

The `SCALE` column shows the running, desired and pending counts. A value of `0/1 (1)` indicates that running = 0, desired = 1 and (1) for 1 pending task that is about to transition to running state. After your task has finished running, note that the desired count will remain 1 and pending value will no longer be shown, so the value will be `0/1` for a finished task.

## Part 4: Check the status of the task

To view the logs for this task, and evaluate its progress, you can use the **task logs** command:

`l0 task logs one-off-task-tsk`  

You will see the following output:
```
alpine
------
Task finished!
```

You can also use the following command for more information in the task.

`l0 -o json task get echo-tsk`

Outputs:

```
[
    {
        "copies": [
            {
                "details": [],
                "reason": "Waiting for cluster capacity to run",
                "task_copy_id": ""
            }
        ],
        "deploy_id": "one-off-task-dpl.2",
        "deploy_name": "one-off-task-dpl",
        "deploy_version": "2",
        "desired_count": 1,
        "environment_id": "demoenv669e4",
        "environment_name": "demo-env",
        "pending_count": 1,
        "running_count": 0,
        "task_id": "echotsk1facd",
        "task_name": "echo-tsk"
    }
]
```

After the task has finished, running `l0 -o json task get echo-tsk` will show a pending_count of 0.

Outputs:

```
...
"copies": [
    {
        "details": [
            {
                "container_name": "alpine",
                "exit_code": 0,
                "last_status": "STOPPED",
                "reason": ""
            }
        ],
        "reason": "Essential container in task exited",
        "task_copy_id": "arn:aws:ecs:us-west-2:856306994068:task/0e723c3e-9cd1-4914-8393-b59abd40eb89"
    }
],
...
"pending_count": 0,
"running_count": 0,
...
```
