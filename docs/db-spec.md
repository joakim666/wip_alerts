# Database specification

## Assumptions
    Marking a field with * means it's immutable. Changes should probably be appened as a new row in the table instead

    Assume all field are not-null unless explicitely stated that they are nullable.


## accounts [append only table]
    Contains all accounts, one per unique device.

    - id* (string) - uuid
    - created_at* (timestamp)

## devices [append only table]
    Contains all device information connected to one account.

    TODO: how to handle changing to another device?

    - id* (string) - uuid
    - device_id* (string) - uuid of device
    - device_type* (string) - ios|android[etc
    - device_info* (string) - device information as json
    - created_at* (timestamp)
    - account_id* (string) - fk accounts:id

## renewals [append only table]
    Contains all renewals connected to an account and the token used to make the renewal.

    - id* (string) - uuid
    - account_id
    - token_id
    - created_at* (timestamp)

## tokens [append only table]
    Contains all created tokens both refresh tokens and access tokens.

    - id* (string) - uuid
    - type*: access_token, refresh_token (enum)
    - data* (string)
    - created_at* (timestamp)
    - account_id* - fk accounts:id

//## account_tokens [append only table]
//    Maps accounts to tokens. I.e. all tokens created for an account.
//
//    - account_id* - fk accounts:id
//    - token_id* - fk tokens:id

## token_status [append only table]
    Holds the status of the token. The current status of a token is the latest entry (the one with the highest id) for the given token_id. 

    - id* (counter) - autoincremented
    - token_id* (string) - fk to tokens:id
    - status*: active, expired, deactivated (enum)
    - created_at* (timestamp) - timestamp of row creation

## api keys [append only table]
    Holds all issued api keys
    - id* (string) - uuid
    - token_id* - fk tokens:id
    - created_at* (timestamp)
    - account_id* - fk accounts:id
    
    
    

