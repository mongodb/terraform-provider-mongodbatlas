# Example - MongoDB Atlas User Data Source

This project provides a straight-forward example for using the Atlas User Data Source.

Variables Required to be set:
- `public_key`: Atlas Programmatic API public key
- `private_key`: Atlas Programmatic API private key
- `user_id`: User ID of the Atlas User to return
- `username`: Username of the Atlas User to return


This example shows the two ways that you can use the data source, either providing the `user_id` or `username` attribute.

For additional documentation, see [MongoDB Atlas API - Get User By ID](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-getuser) and [MongoDB Atlas API - Get User By Username](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-getuserbyusername) respectively.
