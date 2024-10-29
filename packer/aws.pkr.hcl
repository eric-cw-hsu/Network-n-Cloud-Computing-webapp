packer {
  required_plugins {
    amazon = {
      version = "~> 1.0"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

build {
  sources = [
    "source.amazon-ebs.csye6225-ami"
  ]

  provisioner "file" {
    source      = "./app"
    destination = "/tmp/app"
  }

  provisioner "file" {
    source      = "./migrations"
    destination = "/tmp/migrations"
  }

  provisioner "file" {
    source      = "./packer/app.service"
    destination = "/tmp/app.service"
  }

  provisioner "file" {
    source      = "./packer/nginx.conf"
    destination = "/tmp/nginx.conf"
  }

  provisioner "shell" {
    script = "./packer/scripts/os-init.sh"
  }

  provisioner "shell" {
    script = "./packer/scripts/webapp-deploy.sh"
  }

  provisioner "shell" {
    script = "./packer/scripts/amazon-cloudwatch-agent-setup.sh"
  }
}