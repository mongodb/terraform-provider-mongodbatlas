resource "mongodbatlas_resource_policy" "project_ip_access_list" {
  org_id = var.org_id
  name   = "forbid-access-from-anywhere"
  description = "Forbids access from anywhere"

  policies = [
    {
      body = <<EOF
        forbid (
                principal,
                action == ResourcePolicy::Action::"project.ipAccessList.modify",
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
  description = "Forbids AWS and Azure for clusters"
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
      type = " ResourcePolicy::Action"
      id   = "cluster.modify"
    }
    any_resource = true
    when {
      text = "context.cluster.regions.contains(ResourcePolicy::Region::\"gcp:us-east1\")"
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
  value = { for policy in data.mongodbatlas_resource_policies.this.results : policy.name => policy.id }
}
