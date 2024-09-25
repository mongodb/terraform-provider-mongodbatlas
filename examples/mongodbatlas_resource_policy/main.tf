resource "mongodbatlas_resource_policy" "project_ip_access_list" {
  org_id = var.org_id
  name   = "forbid-access-from-anywhere"

  policies = [
    {
      body = <<EOF
        forbid (
                principal,
                action == cloud::Action::"project.edit",
                resource
        )
                when {
                context.project.ipAccessList.contains(ip("0.0.0.0/0"))
        };
EOF
    },
  ]
}

resource "mongodbatlas_resource_policy" "cloud_provider" {
  org_id = var.org_id
  name   = "forbid-cloud-provider"
  policies = [
    {
      body = templatefile("${path.module}/cloud-provider.cedar", {
        CLOUD_PROVIDER = "azure"
      })
    },
    {
      body = templatefile("${path.module}/cloud-provider.cedar", {
        CLOUD_PROVIDER = "aws"
      })
    },
  ]
}

data "cedar_policyset" "cloud_region" {
  policy {
    any_principal = true
    effect        = "forbid"
    action = {
      type = " cloud::Action"
      id   = "cluster.createEdit"
    }
    any_resource = true
    when {
      text = "context.cluster.regions.contains(cloud::region::\"gcp:us-east1\")"
    }
  }
}

resource "mongodbatlas_resource_policy" "cloud_region" {
  org_id = var.org_id
  name   = "forbid-cloud-region"
  policies = [
    {
      body = data.cedar_policyset.cloud_region.text
    },
  ]
}


data "mongodbatlas_resource_policy" "project_ip_access_list" {
  org_id = mongodbatlas_resource_policy.project_ip_access_list.org_id
  id     = mongodbatlas_resource_policy.project_ip_access_list.id
}

data "mongodbatlas_resource_policies" "this" {
  org_id = data.mongodbatlas_resource_policy.project_ip_access_list.org_id

  depends_on = [mongodbatlas_resource_policy.project_ip_access_list, mongodbatlas_resource_policy.cloud_provider, mongodbatlas_resource_policy.cloud_region]
}


output "policy_ids" {
  value = { for policy in data.mongodbatlas_resource_policies.this.resource_policies : policy.name => policy.id }
}
