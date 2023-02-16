locals {

  ip_address_list = [
    {
      ip_address = "47.225.213.178"
      comment    = "IP Address 1"
    },

    {
      ip_address = "47.225.214.179"
      comment    = "IP Address 2"
    },
  ]

  cidr_block_list = [
    {
      cidr_block = "10.1.0.0/16"
      comment    = "CIDR Block 1"
    },
    {
      cidr_block = "12.2.0.0/16"
      comment    = "CIDR Block 2"
    },
  ]
}


resource "mongodbatlas_project_ip_access_list" "ip" {
  for_each = {
    for index, ip in local.ip_address_list :
    ip.comment => ip
  }
  project_id = var.project_id
  ip_address = each.value.ip_address
  comment    = each.value.comment
}


resource "mongodbatlas_project_ip_access_list" "cidr" {

  for_each = {
    for index, cidr in local.cidr_block_list :
    cidr.comment => cidr
  }
  project_id = var.project_id
  cidr_block = each.value.cidr_block
  comment    = each.value.comment
}