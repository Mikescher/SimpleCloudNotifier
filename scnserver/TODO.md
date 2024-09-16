

  TODO
========


#### DO DO DO

 - app-store link in HTML

 - ios purchase verification

 - exerr.New | exerr.Wrap

 - Properly handle UNREGISTERED firebase error (remove token from client?)
   WRN logic/application.go:284 > FCM Delivery failed error="FCM-Request returned 404: 
   {  \"error\": {\n    \"code\": 404,\n    \"message\": \"Requested entity was not found.\",\n    \"status\": \"NOT_FOUND\",\n    \"details\": [\n      {\n        \"@type\": \"type.googleapis.com/google.firebase.fcm.v1.FcmError\",\n        \"errorCode\": \"UNREGISTERED\"\n      }\n    ]\n  }\n}\n" 
   ClientID=CLNGOSVIaCnm5cQmCI0pC5kR MessageID=MSG8w7NvVRm0OtJERnJlEe3C

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

 - endpoint to list all servernames of user (distinct select)

 - weblogin, webapp, ...

 - Pagination for ListChannels / ListSubscriptions / ListClients / ListChannelSubscriptions / ListUserSubscriptions

 - use job superclass (copy from isi/bnet/?), reduce duplicate code

 - admin panel (especially errors and requests)

 - cli app (?)

#### FUTURE

 - Remove compat, especially do not create compat id for every new message...
