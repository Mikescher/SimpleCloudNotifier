

  TODO
========


#### BEFORE RELEASE

 - finish tests (!)

 - migration script for existing data
   apply local deletion in (my) app
   delete excessive dockerwatch messages (directly in db?)

 - app-store link in HTML

 - route to re-check all pro-token (for me)

 - deploy

 - error logging as goroutine, gets all errors via channel,
   (channel buffered - nonblocking send, second channel that gets a message when sender failed )
   (then all errors end up in _second_ sqlite table)
   due to message channel etc everything is non blocking and cant fail in main

 - => implement proper error logging in goext, kinda combines zerolog and wrapped-errors
      copy basic code from bringman, but remove all bm specific stuff and make it abstract
      Register(ErrType) methods, errtypes then as structs
      log.xxx package with same interface as zerolog
      
 - jobs to clear requests-db and logs-db after to only keep X entries...

 - /send endpoint should be compatible with the [ webhook ] notifier of uptime-kuma
   (or add another /kuma endpoint)
   -> https://webhook.site/

 - endpoint to list all servernames of user (distinct select)

 - ios purchase verification

 - move to KeyToken model
     * [X] User can have multiple keys with different permissions
     * [X] compat simply uses default-keys
     * [X] CRUD routes for keys
     * [X] KeyToken.messagecounter
     * [x] update old-data migration to create token-keys
     * [x] unit tests

 - We no longer have a route to reshuffle all keys (previously in updateUser), add a /user/:uid/keys/reset ?
   Would delete all existing keys and create 3 new ones?

 - TODO-comments

 - why do some tests take 5 seconds (= duration of context timeout??)

#### PERSONAL

 - in my script: use `srvname` for sendername

 - switch send script everywhere (we can use the new server, but we need to send correct channels)

 - do i need bool2db()? it seems to work for keytokens without them?

#### UNSURE

 - (?) default-priority for channels

 - (?) ack/read deliveries && return ack-count  (? or not, how to query?)

 - (?) "login" on website and list/search/filter messages

 - (?) make channels deleteable (soft-delete) (what do with messages in channel?)

 - (?) desktop client for notifications

- (?) add querylog (similar to requestlog/errorlog) - only for main-db

#### LATER

 - weblogin, webapp, ...

 - Pagination for ListChannels / ListSubscriptions / ListClients / ListChannelSubscriptions / ListUserSubscriptions

 - Use only single struct for DB|Model|JSON
     * needs sq.Converter implementation
     * needs to handle joined data
     * rfctime.Time...

 - cannot open sqlite in dbbrowsr (cannot parse schema?)
   -> https://github.com/sqlitebrowser/sqlitebrowser/issues/292 -> https://github.com/sqlitebrowser/sqlitebrowser/issues/29266

#### FUTURE

 - Remove compat, especially do not create compat id for every new message...
