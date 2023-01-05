

  TODO
========

-------------------------------------------------------------------------------------------------------------------------------

 - migration script for existing data

 - app-store link in HTML

 - route to re-check all pro-token (for me)

 - tests (!)

 - deploy

 - diff my currently used scnsend script vs the one in the docs here

 - Pagination for ListChannels / ListSubscriptions / ListClients / ListChannelSubscriptions / ListUserSubscriptions

 - cannot open sqlite in dbbrowsr (cannot parse schema?)

- (?) use str-ids (also prevents wrong-joins) -> see psycho

 - error logging as goroutine, get sall errors via channel,
   (channel buffered - nonblocking send, second channel that gets a message when sender failed )
   (then all errors end up in _second_ sqlite table)
   due to message channel etc everything is non blocking and cant fail in main

 - request logging (log all requests with body response, exitcode, headers, uri, route, userid, ..., tx-retries, etc), (trim body/response if too big?)

 - jobs to clear requests-db and logs-db after to only keep X entries...

 -> logs and request-logging into their own sqlite files (sqlite-files are prepped)

-------------------------------------------------------------------------------------------------------------------------------

 - in my script: use (backupname || hostname) for sendername

-------------------------------------------------------------------------------------------------------------------------------

 - (?) default-priority for channels

 - (?) ack/read deliveries && return ack-count  (? or not, how to query?)

 - (?) "login" on website and list/search/filter messages

 - (?) make channels deleteable (soft-delete) (what do with messages in channel?)