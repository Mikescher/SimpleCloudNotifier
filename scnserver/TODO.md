

  TODO
========


#### DO DO DO

 - app-store link in HTML

 - ios purchase verification

 - use goext.ginWrapper
  
 - use goext.exerr

#### UNSURE

 - (?) default-priority for channels

 - (?) "login" on website and list/search/filter messages

 - (?) make channels deleteable (soft-delete) (what do with messages in channel?)

 - (?) desktop client for notifications

 - (?) add querylog (similar to requestlog/errorlog) - only for main-db

 - (?) specify 'type' of message (debug, info, warn, error, fatal)  ->  distinct from priority 

#### LATER

 - do i need bool2db()? it seems to work for keytokens without them?

 - We no longer have a route to reshuffle all keys (previously in updateUser), add a /user/:uid/keys/reset ?
   Would delete all existing keys and create 3 new ones?

 - error logging as goroutine, gets all errors via channel,
   (channel buffered - nonblocking send, second channel that gets a message when sender failed )
   (then all errors end up in _second_ sqlite table)
   due to message channel etc everything is non blocking and cant fail in main
 
 - => implement proper error logging in goext, kinda combines zerolog and wrapped-errors
   copy basic code from bringman, but remove all bm specific stuff and make it abstract
   Register(ErrType) methods, errtypes then as structs
   log.xxx package with same interface as zerolog

 - jobs to clear error-db to only keep X entries... (requests-db already exists)

 - route to re-check all pro-token (for me)

 - /send endpoint should be compatible with the [ webhook ] notifier of uptime-kuma
   (or add another /kuma endpoint)
   -> https://webhook.site/


````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````



{"heartbeat":{"monitorID":89,"status":0,"time":"2023-07-31 18:56:15.374","msg":"timeout of 16000ms exceeded","important":true,"duration":36,"timezone":"Europe/Berlin","timezoneOffset":"+02:00","localDateTime":"2023-07-31 20:56:15"},"monitor":{"id":89,"name":"test","description":null,"pathName":"test","parent":null,"childrenIDs":[],"url":"https://exampleXYZ.com","method":"GET","hostname":null,"port":null,"maxretries":1,"weight":2000,"active":true,"forceInactive":false,"type":"http","interval":20,"retryInterval":20,"resendInterval":0,"keyword":null,"expiryNotification":false,"ignoreTls":false,"upsideDown":false,"packetSize":56,"maxredirects":10,"accepted_statuscodes":["200-299"],"dns_resolve_type":"A","dns_resolve_server":"1.1.1.1","dns_last_result":null,"docker_container":"","docker_host":null,"proxyId":null,"notificationIDList":{"2":true},"tags":[],"maintenance":false,"mqttTopic":"","mqttSuccessMessage":"","databaseQuery":null,"authMethod":null,"grpcUrl":null,"grpcProtobuf":null,"grpcMethod":null,"grpcServiceName":null,"grpcEnableTls":false,"radiusCalledStationId":null,"radiusCallingStationId":null,"game":null,"httpBodyEncoding":"json","includeSensitiveData":false},"msg":"[test] [🔴 Down] timeout of 16000ms exceeded"}


=====================================================================================================================================================================================================


{
  "heartbeat": {
    "monitorID": 89,
    "status": 1,
    "time": "2023-07-31 18:56:57.151",
    "msg": "200 - OK",
    "ping": 55,
    "important": true,
    "duration": 41,
    "timezone": "Europe/Berlin",
    "timezoneOffset": "+02:00",
    "localDateTime": "2023-07-31 20:56:57"
  },
  "monitor": {
    "id": 89,
    "name": "test",
    "description": null,
    "pathName": "test",
    "parent": null,
    "childrenIDs": [],
    "url": "https://example.com",
    "method": "GET",
    "hostname": null,
    "port": null,
    "maxretries": 1,
    "weight": 2000,
    "active": true,
    "forceInactive": false,
    "type": "http",
    "interval": 20,
    "retryInterval": 20,
    "resendInterval": 0,
    "keyword": null,
    "expiryNotification": false,
    "ignoreTls": false,
    "upsideDown": false,
    "packetSize": 56,
    "maxredirects": 10,
    "accepted_statuscodes": [
      "200-299"
    ],
    "dns_resolve_type": "A",
    "dns_resolve_server": "1.1.1.1",
    "dns_last_result": null,
    "docker_container": "",
    "docker_host": null,
    "proxyId": null,
    "notificationIDList": {
      "2": true
    },
    "tags": [],
    "maintenance": false,
    "mqttTopic": "",
    "mqttSuccessMessage": "",
    "databaseQuery": null,
    "authMethod": null,
    "grpcUrl": null,
    "grpcProtobuf": null,
    "grpcMethod": null,
    "grpcServiceName": null,
    "grpcEnableTls": false,
    "radiusCalledStationId": null,
    "radiusCallingStationId": null,
    "game": null,
    "httpBodyEncoding": "json",
    "includeSensitiveData": false
  },
  "msg": "[test] [✅ Up] 200 - OK"
}


=====================================================================================================================================================================================================


{
  "heartbeat": {
    "monitorID": 89,
    "status": 0,
    "time": "2023-07-31 18:57:44.037",
    "msg": "getaddrinfo ENOTFOUND exampleasdsda.com",
    "important": true,
    "duration": 20,
    "timezone": "Europe/Berlin",
    "timezoneOffset": "+02:00",
    "localDateTime": "2023-07-31 20:57:44"
  },
  "monitor": {
    "id": 89,
    "name": "test",
    "description": null,
    "pathName": "test",
    "parent": null,
    "childrenIDs": [],
    "url": "https://exampleasdsda.com",
    "method": "GET",
    "hostname": null,
    "port": null,
    "maxretries": 1,
    "weight": 2000,
    "active": true,
    "forceInactive": false,
    "type": "http",
    "interval": 20,
    "retryInterval": 20,
    "resendInterval": 0,
    "keyword": null,
    "expiryNotification": false,
    "ignoreTls": false,
    "upsideDown": false,
    "packetSize": 56,
    "maxredirects": 10,
    "accepted_statuscodes": [
      "200-299"
    ],
    "dns_resolve_type": "A",
    "dns_resolve_server": "1.1.1.1",
    "dns_last_result": null,
    "docker_container": "",
    "docker_host": null,
    "proxyId": null,
    "notificationIDList": {
      "2": true
    },
    "tags": [],
    "maintenance": false,
    "mqttTopic": "",
    "mqttSuccessMessage": "",
    "databaseQuery": null,
    "authMethod": null,
    "grpcUrl": null,
    "grpcProtobuf": null,
    "grpcMethod": null,
    "grpcServiceName": null,
    "grpcEnableTls": false,
    "radiusCalledStationId": null,
    "radiusCallingStationId": null,
    "game": null,
    "httpBodyEncoding": "json",
    "includeSensitiveData": false
  },
  "msg": "[test] [🔴 Down] getaddrinfo ENOTFOUND exampleasdsda.com"
}

````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````````


 - endpoint to list all servernames of user (distinct select)

 - weblogin, webapp, ...

 - Pagination for ListChannels / ListSubscriptions / ListClients / ListChannelSubscriptions / ListUserSubscriptions

 - Use only single struct for DB|Model|JSON
     * needs sq.Converter implementation
     * needs to handle joined data
     * rfctime.Time...

 - use job superclass (copy from isi/bnet/?), reduce duplicate code

 - admin panel (especially errors and requests)

 - cli app (?)

#### FUTURE

 - Remove compat, especially do not create compat id for every new message...
