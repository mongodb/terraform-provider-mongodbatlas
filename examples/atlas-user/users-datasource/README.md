# Example - MongoDB Atlas Users Data Source

This project provides a straight-forward example for using the Atlas Users Data Source.

Variables Required to be set:
- `public_key`: Atlas Programmatic API public key
- `private_key`: Atlas Programmatic API private key
- `org_id`: Org ID that identifies the organization whose users you want to return
- `project_id`: Project ID that identifies the project whose users you want to return
- `team_id`: Team ID that identifies the team whose users you want to return

The example demonstrates the three ways you can use this data source to obtain users from an Atlas organization, project, or team.

For additional documentation, see:

- For obtaining users of an Organization: [MongoDB Atlas API - List Organization Users](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Organizations/operation/listOrganizationUsers) 
- For obtaining users of a Project: [MongoDB Atlas API - List Project Users](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Projects/operation/listProjectUsers)
- For obtaining users of a Team: [MongoDB Atlas API - List Team Users](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Teams/operation/listTeamUsers)

