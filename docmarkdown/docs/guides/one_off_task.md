# Deployment guide: Guestbook one-off task

In this example, you will learn how to use layer0 to run a one-off task. In this case, it will be to run a task to restore the [guestbook application](/guides/guestbook) from a backup.

---

## Before you start
In order to complete the procedures in this section, you must install and configure Layer0 v0.8.3 or later. If you have not already configured Layer0, see the [installation guide](/setup/install). If you are running an older version of Layer0, see the [upgrade instructions](/setup/update#upgrading-older-versions-of-layer0).

This guide expands upon the [Guestbook deployment guide](/guides/guestbook) deployment guide. You must complete the procedures in that guide before you can complete the procedures listed here. After completing the procedures in the Guestbook guide, your Layer0 should contain a service named "guestbooksvc", running a deploy named "guestbook", behind a load balancer named "guestbooklb", all within an environment named "demo".

## Part 1: Prepare the task definition

1. Download the [Guestbook One-off Task Definition](https://gitlab.imshealth.com/xfra/layer0-samples/blob/master/1offtask/Dockerrun.aws.json) and save it to your computer as **GuestbookRestore.Dockerrun.aws.json**.
2. Edit the `GUESTBOOK_URL` environment variable for the `l0-guestbook-restore` container to the url of the loadbalancer running your guestbook application. This url can be obtained by looking at the output of the command

<span style="padding-left:2em">**l0 loadbalancer get guestbooklb**</span>

3. Edit the `BACKUP_FILE_URL` environment variable for the `l0-guestbook-restore` container to the url of the backup file that you wish to restore from. A sample backup file is provided at [https://gitlab.imshealth.com/xfra/layer0-samples/raw/master/1offtask/backup.txt](https://gitlab.imshealth.com/xfra/layer0-samples/raw/master/1offtask/backup.txt).

## Part 2: Create a deploy
Next, you will create a new deploy for the task.

**To create a new deploy:**

At the command prompt, run the following command:

<span style="padding-left:2em">**l0 deploy create GuestbookRestore.Dockerrun.aws.json guestbookrestore**</span>

You will see the following output:
```
DEPLOY ID           DEPLOY NAME        VERSION
guestbookrestore.1  guestbookrestore   1
```

## Part 3: Create the task
At this point, you can use the **task create** command to begin an instance of the task deployed above. This task requires two environment variables to be supplied to the `l0-guestbook-restore` container: `GUESTBOOK_URL` and `BACKUP_FILE_URL`. These were collected in Part 1 of this guide.

To run the task, use the following command:

<span style="padding-left:2em">**l0 task create demo guestbookrestore guestbookrestore**</span>

You will see the following output:
```
TASK ID       TASK NAME         ENVIRONMENT  DEPLOY              SCALE
guestbo851c9  guestbookrestore  demo         guestbookrestore:2  0/1 (1)
```

## Part 4: Wait for the task to complete

### Check the logs for the task

To see the logs for this task, and evaluate progress, use the command:

<span style="padding-left:2em">**l0 task logs guestbookrestore**</span>

Once it has completed, check your guestbook url, and note that the entries have been replaced with the contents of your backup file.
