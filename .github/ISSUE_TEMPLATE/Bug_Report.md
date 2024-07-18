---
name: Bug report [Deprecated]
about: Something unexpected happen? Report an error, crash, or an incorrect behavior here.

---

Hello!

Thank you for opening an issue. **These GitHub issues** are only for community reporting and assistance; as such, we don't have a guaranteed SLA.

**If your issue relates to Terraform itself**, please open it in the Terraform repository https://github.com/hashicorp/terraform/issues.

**If you have an active MongoDB Atlas Support contract**, the best way for us to assist you with the Terraform MongoDB Atlas Provider is through a [support ticket](https://support.mongodb.com/).

**Please note:** In order for us to provide the best experience in supporting our customers, we kindly ask to make sure that all the following sections are correctly filled with all the required information. Our support will **prioritise** issues that contain **all the required** information that follows the **"one-click reproducible issues" principle** (see below).

**Please note:** In order for us to provide the best experience in supporting our customers, we kindly ask to make sure that all the following sections are correctly filled with all the required information. 
Our support will prioritise issues that contain all the required information that follows the [one-click reproducible issues principle](../../README.md#one-click-reproducible-issues-principle).


### Terraform CLI and Terraform MongoDB Atlas Provider Version
Please ensure your issue is reproducible on a supported Terraform version. You may review our [Terraform version compatibility matrix](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/docs/index.md#hashicorp-terraform-version-compatibility-matrix) to know more.
<!---
Run `terraform version` to show the version, and paste the result for Terraform and the MongoDB Atlas Provider between the ``` marks below.

If you are not running the latest version of Terraform and the MongoDB Atlas Provider, please try to reproduce the issue in a more recent version in case it has been fixed already. 
-->

```hcl
# Copy-paste your version.tf and provider.tf (or equivalent) here
```

### Terraform Configuration File 
<!--
Paste here all the required parts of your Terraform MongoDB Atlas Provider configuration file below. You can alternatively consider creating a [**GitHub Gist**](https://gist.github.com/) with the information and share it here.

Share your configuration by **using variables**. Create a `variables.tf` file and share its content. Please be sure to redact any sensitive information; common examples include API Keys, passwords, other secrets, or any items your project/business may consider sensitive such as hostnames or usernames.
-->

```hcl
# Copy-paste all your configuration (e.g. main.tf and variable.tf) info here
```

### Steps to Reproduce

Please be **as detailed as possible.**
* If the issue **involves version changes** (e.g. `run plan with version 1.k.x then re-run with version 1.k.y`), please make sure you report the detailed sequence on commands to execute.
* If the issue requires **changes from the UI**, please document them **preferably with screenshots**.

```hcl
# Write here the detailed list of required steps.
```

### Expected Behavior
<!--
What should have happened?
-->

### Actual Behavior
<!--
What actually happened?
-->

### Debug Output 
<!--
Note: Debug output can be incredibly helpful in narrowing down an issue.

Full debug output can be obtained by running Terraform with the environment variable `TF_LOG=trace`. Please create either a GitHub Gist or attach a file containing the debug output. Please do _not_ paste the debug output in the issue, since debug output can be very long.

Debug output may contain sensitive information. Please review it before posting publicly, and if you are concerned feel free to redact it.
-->

### Crash Output
<!--
If the console output indicates that Terraform crashed, please either share a link to a GitHub Gist containing the output of the `crash.log` file or attach the file.
-->

### Additional Context
<!--
Are there anything atypical about your situation that we should know? 
-->

### References
<!--
Are there any other related GitHub issues (open or closed) or Pull Requests that should be linked here? 
-->
