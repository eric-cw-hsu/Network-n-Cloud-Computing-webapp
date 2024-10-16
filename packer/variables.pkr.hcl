variable "aws_region" {
  type    = string
  default = "us-west-2"
}

variable "source_ami" {
  type    = string
  default = "ami-04dd23e62ed049936" # Ubuntu 24.04 LTS us-west-2
}

variable "ssh_username" {
  type    = string
  default = "ubuntu"
}

variable "subnet_id" {
  type    = string
  default = "subnet-0d4d366e276bb292a" # us-west-2a
}