---
page_title: "Migration Guide: Atlas User Management"
sidebar_current: "docs-mongodbatlas-guides-200"
---

# Migration Guide: Atlas User Management

## Overview

With MongoDB Atlas Terraform Provider `2.0.0`, several attributes and resources were deprecated in favor of new, assignment-based resources.  
These changes improve **clarity, separation of concerns, and alignment with Atlas APIs**.  
This guide covers migrating to the new resources/attributes for Atlas user management in context of **organization, teams, and projects**:

## Quick Finder: What changed

- **Org membership:** The `mongodbatlas_org_invitation` resource is deprecated. Use `mongodbatlas_cloud_user_org_assignment`.  
  → See [Org Invitation to Cloud User Org Assignment](#migr-org-invitation)

- **Team membership:** The `usernames` attribute on `mongodbatlas_team` is deprecated. Use `mongodbatlas_cloud_user_team_assignment`.  
  → See [Team Usernames to Cloud User Team Assignment](#migr-team-usernames)

- **Project team roles:** The `teams` block inside `mongodbatlas_project` is deprecated. Use `mongodbatlas_team_project_assignment`.  
  → See [Project Teams to Team Project Assignment](#migr-project-teams)

- **Project membership:** The `mongodbatlas_project_invitation` resource is deprecated. Use `mongodbatlas_cloud_user_project_assignment`.  
  → See [Project Invitation to Cloud User Project Assignment](#migr-project-invitation)
  
- **User details:** The `mongodbatlas_atlas_user` and `mongodbatlas_atlas_users` data sources are deprecated.  
  Use `mongodbatlas_cloud_user_org_assignment` for a single user in an org, and the `users` attributes on `mongodbatlas_organization`, `mongodbatlas_project`, or `mongodbatlas_team` for listings.  
  → See [Atlas User/Users Data Sources](#migr-atlas-user-users)

These updates ensure that **organization membership, team membership, and project assignments** are modeled as explicit and independent resources — giving you more flexible control over Atlas access management.

---

## Before You Begin
- Backup your Terraform state file: https://developer.hashicorp.com/terraform/cli/commands/state  
- Use MongoDB Atlas Terraform Provider **v2.0.0+**.  
- Terraform version requirements:
  - **v1.5+** for **import blocks**: https://developer.hashicorp.com/terraform/language/import  
  - **v1.1+** for **moved blocks** (works in modules): https://developer.hashicorp.com/terraform/language/moved  
  - **v1.7+** for **removed blocks**: https://developer.hashicorp.com/terraform/language/resources/syntax#removing-resources  
- Import blocks **cannot** live inside modules (root-only).

---


<a id="migr-org-invitation"></a>
<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Org Invitation to Cloud User Org Assignment</span></summary>


</details>

<a id="migr-team-usernames"></a>
<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Team Usernames to Cloud User Team Assignment</span></summary>


</details>

<a id="migr-project-teams"></a>
<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Project Teams to Team Project Assignment</span></summary>


</details>

<a id="migr-project-invitation"></a>
<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Project Invitation to Cloud User Project Assignment</span></summary>


</details>


<a id="migr-atlas-user-users"></a>
<details>
  <summary><span style="font-size:1.4em; font-weight:bold;">Atlas User/Users Data Sources</span></summary>


</details>
