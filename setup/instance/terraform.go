package instance

var TFFiles = map[string][]byte{
	"variables.tf": []byte(VARIABLES_TF),
	"main.tf":      []byte(MAIN_TF),
	"outputs.tf":   []byte(OUTPUTS_TF),
}

const MAIN_TF = `
 hey
`

const VARIABLES_TF = `
variable "aws_access_key" {
  description = "AWS access key"
}

variable "aws_secret_key" {
  description = "AWS secret key"
}

variable "aws_region" {
  description = "AWS region"
}
`

const OUTPUTS_TF = `
output "test" {
    value = "todo"
}
`
