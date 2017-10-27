variable "file_path" {
  default = "test.txt"
}

data "template_file" "test" {
  template = "some file content"
}

resource "null_resource" "local" {
  triggers {
    template = "${data.template_file.test.rendered}"
  }

  provisioner "local-exec" {
    command = "echo \"${data.template_file.test.rendered}\" > ${var.file_path}"
  }
}

output "file_path" {
  value = "${var.file_path}"
}
