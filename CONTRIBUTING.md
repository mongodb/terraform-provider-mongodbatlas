Contributing
---------------------------

# Contributing

## Workflow

MongoDB welcomes community contributions! If you’re interested in making a contribution to  Terraform MongoDB Atlas provider, please follow the steps below before you start writing any code:

1. Sign the [contributor's agreement](http://www.mongodb.com/contributor). This will allow us to review and accept contributions.
1. Read the [Terraform contribution guidelines](https://www.terraform.io/docs/extend/community/contributing.html) and the [README](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/README.md) in this repo
1. Reach out by filing an [issue](https://github.com/mongodb/terraform-provider-mongodbatlas/issues) to discuss your proposed contribution, be it a bug fix or feature/other improvements.  

After the above 3 steps are completed and we've agreed on a path forward:
1. Fork the repository on GitHub
1. Create a branch with a name that briefly describes your submission
1. Implement your feature, improvement or bug fix, ensuring it adheres to the [Terraform Plugin Best Practices](https://www.terraform.io/docs/extend/best-practices/index.html)
1. Ensure you follow the [Terraform Plugin Testing requirements](https://www.terraform.io/docs/extend/testing/index.html).
1. Add comments around your new code that explain what's happening
1. Commit and push your changes to your branch then submit a pull request against the current release branch, not master.  The naming scheme of the branch is `release-staging-v#.#.#`. Note: There will only be one release branch at a time.  
1. A repo maintainer will review the your pull request, and may either request additional changes or merge the pull request.

## Documentation Best Practises

1. In our documentation, when a resource field allows a maximum of only one item, we do not format that field as an array. Instead, we create a subsection specifically for this field. Within this new subsection, we enumerate all the attributes of the field. Let's illustrate Example: [cloud_backup_schedule.html.markdown](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/website/docs/r/cloud_backup_schedule.html.markdown?plain=1#L207)
