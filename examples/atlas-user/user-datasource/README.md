# Example - MongoDB Atlas User Data Source

This project aims to provide a very straight-forward example for using the Atlas User Data Source.

Variables Required to be set:
- `public_key`: Atlas public key
- `private_key`: Atlas  private key
- `user_id`: User ID of the Atlas User that will be fetched
- `username`: Username of the Atlas User that will be fetched


Example shows the two possible way the data source can be used, either providing the `user_id` or `username` attribute.

For additional documentation, you can reference to [MongoDB Atlas API - Get User By ID](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/MongoDB-Cloud-Users/operation/getUser) and [MongoDB Atlas API - Get User By Username](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/#tag/MongoDB-Cloud-Users/operation/getUserByUsername) respectively.
