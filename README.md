# CERTMAN

SSL certificate validity checker. Built with speed and concurrency in mind.

## USAGE
Usage: certman -h [FILE] -s [FILE] -m [BOOL]  
-h - Path to hosts.json file. Defaults to ./hosts.json. This file is required.  
-s - Path to settings.json file. Defaults to ./settings.json. Optional.  
-m - Flag specifying if mail is to be sent. Defaults to false. Optional.  

## STRUCTURE OF CONFIG FILES
### EXAMPLE HOSTS FILE
Unless specified as shown below, default port is 443.
```json
{
    "server.example.com":[
        "example.com",
        "mail.example.com",
        "mail.example.com:993"
    ],
    "0.0.0.0":[
        "www.example.com"
    ]
}
```
### EXAMPLE SETTINGS FILE
Available options in settings file:
```json
{
    "emailAddr":"john@example.com",
    "emailPass": "pa$$word",
    "emailServer": "mail.example.com",
    "emailPort": "587",
    "emailDest": [
        "admin@example.com",
        "john@example.com"
    ],
    "concurrencyLimit": 100,
    "days": [
        7,
        14,
        30,
        60
    ]
}
```