# Gorm Updates cannot updates virtual attributes

## Quick Started 

    # Get code
    $ go get -u github.com/fengjh/gorm-with-rows

    # Set environment variables
    $ source .env

    # Setup postgres database
    $ postgres=# CREATE USER gorm_cannot_updates_virtual_attributes;
    $ postgres=# CREATE DATABASE gorm_cannot_updates_virtual_attributes OWNER gorm_cannot_updates_virtual_attributes;

    # Run Application
    $ go test -run TestUpdateUserPassword
