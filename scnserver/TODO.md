

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

 - diff my currently used scnsend script vs the one in the docs here

- (?) use str-ids (hide counts and prevents wrong-joins) -> see psycho
  -> ensre that all queries that return multiple are properly ordered
  -> how does it work with existing data? 
  -> do i care, there are only 2 active users... (are there?)

 - convert existing user-ids on compat /send endpoint

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

 - return channel as "[..] asdf" in compat methods (mark clients as compat and send compat FB to them...)
   (then we can replace the old server without switching phone clients)
   (still needs switching of the send-script)

 - move to KeyToken model
     * [X] User can have multiple keys with different permissions
     * [X] compat simply uses default-keys
     * [X] CRUD routes for keys
     * [X] KeyToken.messagecounter
     * [ ] update old-data migration to create token-keys
     * [ ] unit tests

 - We no longer have a route to reshuffle all keys (previously in updateUser), add a /user/:uid/keys/reset ?
   Would delete all existing keys and create 3 new ones?

#### PERSONAL

 - in my script: use `srvname` for sendername

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

 - cannot open sqlite in dbbrowsr (cannot parse schema?)
   -> https://github.com/sqlitebrowser/sqlitebrowser/issues/292 -> https://github.com/sqlitebrowser/sqlitebrowser/issues/29266

#### FUTURE

 - Remove compat, especially do not create compat id for every new message...
