#/bin/sh
	for os in linux darwin windows; do \
        cd build/$$os/bin || exit ; \
		if [ ! -f terraform.zip ]; then
			wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_$$os_amd64.zip \
				-O build/$$os/bin/terraform.zip
		fi
		unzip terraform.zip ; \
		find . -type f \
			-not -name 'terraform' \
			-not -name 'terraform.exe' \
			-not -name 'terraform-provider-aws*' \
			-not -name 'terraform-provider-template*' \
		| xargs rm ; \
		cd ../../.. || exit ; \
    done
