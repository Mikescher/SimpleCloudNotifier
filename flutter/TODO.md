
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

 
 -----

 - [ ] Switch server to sq style from faby
        - [ ] switch from mattn to go-sqlite
        - [ ] Single struct for model/db/json
        - [ ] use sq.Query | sq.Update | sq.InsertAndQuery | ....
        - [ ] sq.DBOptions - enable CommentTrimmer and DefaultConverter
        - [ ] run unit-tests...