# MongoDB Atlas Provider -- Atlas Project with custom limits
This example creates a Project with tags and defines custom values for certain limits.
It also shows how you can ignore a specific tag key.

Variables Required to be set:
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `org_id`: Organization ID where project will be created

## How to ignore tags managed outside of terraform

### Single Tag Key(s) are managed outside terraform
Consider this code:

```hcl
resource "mongodbatlas_project" "tags_example" {
  name   = "tags-readme-example1"
  org_id = var.org_id

  tags = {
    Owner       = "Terraform"
    Environment = "Example"
    Team        = "tf-experts"
    CurrentDRI  = "unset"
  }
  lifecycle {
    ignore_changes = [
      tags["CostCenter"]
      tags["CurrentDRI"]
    ]
  }
}
```

- Notice how `CostCenter` does **not** exist in the tags defined above.
  - If someone adds the tag outside terraform and we run `terraform apply`, the output will be `No Changes`
  - However, if `tag["CostCenter]` is excluded from the `ignore_changes` section, `terraform apply` would detect a plan drift and remove the `CostCenter` tag
- Notice how `CurrentDRI` exists in the tags defined above
  - This can be useful when you want an initial empty value, but someone will update the value outside of Terraform

### All Tag Keys are managed outside terraform
You can ignore the `tags` field:

```hcl
resource "mongodbatlas_project" "tags_example" {
  name   = "tags-readme-example2"
  org_id = var.org_id

  lifecycle {
    ignore_changes = [
      tags
    ]
  }
} 
  ```
