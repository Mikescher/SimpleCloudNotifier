
# TODO

 - [ ] Message List
        * [ ] CRUD
 - [ ] Message Big-View
 - [ ] Search/Filter Messages
 - [ ] Channel List
        * [ ] Show subs
        * [ ] CRUD
        * [ ] what about unsubbed foreign channels? - thex should still be visible (or should they, do i still get the messages?)
 - [ ] Sub List
        * [ ] Sub/Unsub/Accept/Deny
 - [ ] Debug List (Show logs, requests)
 - [ ] Key List
        * [ ] CRUD
 - [ ] Auto R-only key for admin, use for QR+link+send
 - [ ] settings
 - [ ] notifications
 - [ ] push navigation stack
 - [ ] read + migrate old SharedPrefs (or not? - who uses SCN even??)
 - [ ] Account-Page
 - [ ] Logout
 - [ ] Send-page

 - [ ] Still @ERROR on scn-init, but no logs? - better persist error (write in SharedPrefs at error_$date=txt ?), also perhaps print first error line in scn-init notification?

 -----

# TODO iOS specific

 - [ ] payment / pro
 - [ ] show notifiactions (foreground/background/etc)
 - [ ] handle click-on-notifications should open message
 - [ ] share message
 - [ ] scan QR

 -----

# TODO Server

 - [ ] Switch server to sq style from faby
        - [ ] switch from mattn to go-sqlite
        - [ ] Single struct for model/db/json
        - [ ] use ginext
        - [ ] use sq.Query | sq.Update | sq.InsertAndQuery | ....
        - [ ] sq.DBOptions - enable CommentTrimmer and DefaultConverter
        - [ ] run unit-tests...
        - [ ] Copy db.Migrate code

 - [ ] Disable compat | remove code 
        - [x] compat message title
        - [ ] ...
 - [ ] RWLock directly in go - prevent/reduce db-locked exception
