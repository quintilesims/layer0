provider "layer0" {
    endpoint = "${var.endpoint}"
    token = "${var.token}"
    skip_ssl_verify = true
}

resource "random_pet" "environment_names" {
    length = 1
    count = "${var.num_environments}"
}

resource "layer0_environment" "te" {
    name = "${element(random_pet.environment_names.*.id, count.index)}"
    size = "t2.micro"
    count = "${var.num_environments}"
}
