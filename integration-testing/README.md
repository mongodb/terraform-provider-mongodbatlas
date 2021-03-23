### Integration tests

The integration tests required extra credentials, like aws, azure. 
In order to execute the complete terra cycle (init, apply, destroy)

For all the testing it needs the common environment variables 
```
    MONGODB_ATLAS_PROJECT_ID
    MONGODB_ATLAS_PUBLIC_KEY
    MONGODB_ATLAS_PRIVATE_KEY
```

For especific aws related interactions 
```
    AWS_ACCESS_KEY_ID
    AWS_SECRET_ACCESS_KEY
    AWS_REGION

    AWS_CUSTOMER_MASTER_KEY_ID (only cloud at rest)

```