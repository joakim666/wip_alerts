# Bolt datamodel

Buckets

## Refresh tokens
    - Key: uuid (RefreshToken:id)
    - Value:
        RefreshToken
        - id (uuid)
        - account_id (uuid) - fk Account:id
        - rights (string, comma separated - api etc
        + status (nested bucket) (key: auto-increment counter)
            - id (uuid)
            - status (string)
            - timestamp (timestamp)
        ...

## Accounts
    - Key: uuid (Account:id)
    - Value (map):
        Account
            - id (uuid)
            - created_at (timestamp)
            - devices (nested bucket)
            - renewals (nested bucket) (key: auto-increment counter)
            - access_tokens (nested bucket) (key: auto-increment counter)
                + id (uuid)
                + data (string) - the raw data
                + rights TODO: TBD
                + created_at (timestamp)
            - api_keys (nested bucket) (key: auto-increment counter)
                + id (uuid)
                + refreshToken_id (uuid) - fk: RefreshToken:id
                + description (string)
                + created_at (timestamp)
            - new_alerts (nested bucket) (key: auto-increment counter)
                Contains alerts that have not yet been returned by the API. Once an alert has been returned by the API it is moved to the 'alerts' bucket.
                + alerts (datatype Alert)
            - alerts (nested bucket) (key: auto-increment counter)
                + alerts (datatype Alert)
            - heartbeats (nested bucket) (key: Heartbeat:identifier)
                Only the latest heartbeat per api key is saved
                + api_key_id (uuid): fk: api_keys:id
                + executed_at (string) - the date time in ISOXXXX format when this heartbeat was received

## Devices - nested bucket with Account:id (uuid) as key
    - Key: uuid (Device:id)
    - Value (map):
        Device
            - id (uuid)
            - device_id (string) - uuid of device
            - device_type (string) - ios|android[etc
            - device_info (string) - device information as json
            - created_at (timestamp)

## Renewals - nested bucked with Account:id (uuid) as key
    - Key: uuid (Renewal:id)
    - Value (map):
        Renewal
            - id (uuid)
            - refreshToken_id (uuid) - fk: RefreshToken:id
            - created_at (timestamp)

## APIKeys - nested bucked with Account:id (uuid) as key
    - Key: uuid (APIKey:id)
    - Value (map):
        APIKey
            - id (uuid)
            - refreshToken_id (uuid) - fk: RefreshToken:id
            - description (string)
            - created_at (timestamp)




## Datatype Alert
    + id (uuid)
    + api_key_id (uuid) - the api_key used to report the alert
    + title (string, required)
        A name or short sentence describing the alert.
    + short_description (string, required)
        A few sentences describing the alert in some detail.
    + long_description (string, optional)
        Complete details of the alert.
    + priority: high, normal, low (enum)
    + status: active, archived (enum)
    + created_at (timestamp)
                    
                
            
        

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

    
    

    

