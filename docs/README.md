# Layer0 Documentation - Developer Notes


## Local Development

- `make deps` to ensure you have `mkdocs` and `awscli` installed.
- Make changes as desired.
- `make build` to compile everything together into the static `site` directory.
- `mkdocs serve` to spin up a local dev server and test with a browser.
- `make deploy` to push your changes to the S3 bucket.
     - **Note:** Requires you have environment variables `AWS_ACCESS_KEY_ID`
     and `AWS_SECRET_ACCESS_KEY` set to an IAM user with access to the S3 bucket.
- `make gh-deploy` to push your changes to github pages.


## Production

### _External Access (Requires Login)_

Changes to the documentation are deployed to S3 as above. A Docker container
(hosted on AWS via Layer0) pulls the latest from S3 every twenty seconds
and the static site is served up by NginX. In normal, everyday operation, the
Docker image and the Layer0 deployment should not need to be touched.

#### _Updating the Docker Image_

- Build and push the image, replacing `tag` with the version number

```
docker build -t d.ims.io/xfra/l0-docs:tag .
docker push d.ims.io/xfra/l0-docs:tag
```

- Update `docs.Dockerrun.aws.json` to use the new image tag
     - **Note:** Make sure you configure the Dockerrun file to have environment
     variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` set to an IAM
     user with access to the S3 bucket.
- Update the Layer0 service


### _IMS Network Access_

We host this website in S3.

**Bucket Name**: `docs.xfra.ims.io`

* Create an S3 bucket via `cloudformation.json`
    * The bucket is configured to be a website, and to only be accessible from the office IPs
    * This has already been done with the bucket name `docs.xfra.ims.io`
* Configure Route53 to point to the S3 bucket
    * Route 53 > Hosted Zones > xfra.ims.io > Create Record Set >
    * **Name**: docs.xfra.ims.io
    * **Alias**: Yes
    * **Alias Target**: s3-website-us-west-2.amazonaws.com.
        * Should populate from the drop down.  Bucket Name and URL must align (docs.xfra.ims.io in this case)
* Build and deploy the website (**Note**: requires you have `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` set to an IAM user with access to this s3 bucket)

```
$ mkdocs build --clean
$ aws --region us-west-2 s3 sync site s3://docs.xfra.ims.io --delete
```
