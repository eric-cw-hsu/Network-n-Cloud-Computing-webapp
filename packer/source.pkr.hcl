source "amazon-ebs" "csye6225-ami" {
  region          = "${var.aws_region}"
  ami_name        = "csye6225-fall2024-app-${formatdate("YYYY-MM-DD hh-mm-ss", timestamp())}"
  ami_description = "CSYE6225 Fall 2024 Application AMI"

  ami_regions = [
    "us-west-2",
  ]

  aws_polling {
    delay_seconds = 120
    max_attempts  = 50
  }

  instance_type = "${var.instance_type}"
  source_ami    = "${var.source_ami}"
  ssh_username  = "${var.ssh_username}"
  subnet_id     = "${var.subnet_id}"

  launch_block_device_mappings {
    device_name           = "/dev/sda1"
    volume_size           = var.instance_volume_size
    volume_type           = "${var.instance_volume_type}"
    delete_on_termination = true
  }

  ami_users = var.shared_user_ids
}