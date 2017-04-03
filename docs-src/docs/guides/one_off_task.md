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
At this point, you can use the **task create** command to begin an instance of the task deployed above.

To run the task, use the following command:

`l0 task create demo-env one-off-task-tsk one-off-task-dpl:latest`

You will see the following output:
```
TASK ID       TASK NAME         ENVIRONMENT  DEPLOY              SCALE
one-off851c9  one-off-task-tsk  demo-env     one-off-task-dpl:1  0/1 (1)
```

## Part 4: Check the status of the task

To view the logs for this task, and evaluate its progress, you can use the **task logs** command:

`l0 task logs one-off-task-tsk`  

You will see the following output:
```
alpine
------
Task finished!
```
