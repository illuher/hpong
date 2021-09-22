terraform {
  backend "etcdv3" {
    endpoints = ["127.0.0.1:2379"]
    lock      = true
    prefix    = "terraform-state/hpong/"
  }

  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "~> 2.13.0"
    }
  }
}

provider "docker" {}

resource "docker_image" "hpong" {
  name         = "illuher/hpong:0.0.3-scratch"
  keep_locally = false
}

resource "docker_container" "hpong_0-0-3_scratch" {
  image = docker_image.hpong.latest
  name  = "hpong"
  ports {
    internal = 8081
    external = 8081
  }
}