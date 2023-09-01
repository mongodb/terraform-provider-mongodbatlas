# Example - MongoDB Atlas Users Data Source

This project aims to provide a very straight-forward example for using the Atlas Users Data Source.

Variables Required to be set:
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `org_id`: Org ID that identifies the organization whose users you want to return
- `project_id`: Project ID that identifies the project whose users you want to return
- `team_id`: Team ID that identifies the team whose users you want to return

The example demonstrates the three ways this data source can be used to obtain users from an organization, project, or team.

For additional documentation, you can reference the following documentation for each use case:

- For obtaining users of an Organization: [MongoDB Atlas API - List Organization Users](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Organizations/operation/listOrganizationUsers) 
- For obtaining users of a Project: [MongoDB Atlas API - List Project Users](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Projects/operation/listProjectUsers)
- For obtaining users of a Team: [MongoDB Atlas API - List Team Users](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/Teams/operation/listTeamUsers)

