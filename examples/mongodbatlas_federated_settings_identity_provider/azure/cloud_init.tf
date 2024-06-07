locals {
  file_stem        = "pymongo_oidc"
  mongodb_oidc_uri = "${local.mongodb_uri}&authMechanism=MONGODB-OIDC&authMechanismProperties=ENVIRONMENT:azure,TOKEN_RESOURCE:${urlencode(var.token_audience)}"
  py_oidc_connect = templatefile("${path.module}/${local.file_stem}.tmpl.sh", {
    DATABASE    = var.insert_record_database
    COLLECTION  = var.insert_record_collection
    RECORD      = jsonencode(var.insert_record_fields)
    MONGODB_URI = local.mongodb_oidc_uri
    OS_USER     = var.vm_admin_username
  })
  init_cfg = yamlencode({
    # https://cloudinit.readthedocs.io/en/latest/reference/examples.html#writing-out-arbitrary-files
    write_files = [
      {
        content = file("${path.module}/${local.file_stem}.py")
        path    = "/home/${var.vm_admin_username}/${local.file_stem}.py"
        # cannot use this since the adminuser is created in cloud-init after files are written
        # owner = "${var.vm_admin_username}:${var.vm_admin_username}"
      }
    ]
    package_update = true
  })
}
data "cloudinit_config" "this" {
  gzip          = true
  base64_encode = true

  part {
    filename     = "init.cfg"
    content_type = "text/cloud-config"
    content      = local.init_cfg
  }

  part {
    content_type = "text/x-shellscript"
    content      = local.py_oidc_connect
  }
}
