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

variable "shared_user_ids" {
  type    = list(string)
  default = ["761018880006"]
}

variable "instance_type" {
  type    = string
  default = "t2.micro"
}

variable "instance_volume_size" {
  type    = number
  default = 8
}

variable "instance_volume_type" {
  type    = string
  default = "gp2"
}