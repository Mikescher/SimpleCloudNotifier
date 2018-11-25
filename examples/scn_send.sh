#!/usr/bin/env bash

#
# Call with   `scn_send title`
#        or   `scn_send title content`
#        or   `scn_send title content priority`
#
#

if [ "$#" -lt 1 ]; then
    echo "no title supplied via parameter"
    exit 1
fi

################################################################################
# INSERT YOUR DATA HERE                                                        #
################################################################################
user_id=999
user_key="????????????????????????????????????????????????????????????????"
################################################################################

title=$1
content=""
sendtime=$(date +%s)

if [ "$#" -gt 1 ]; then
    content=$2
fi

priority=1

if [ "$#" -gt 2 ]; then
    priority=$3
fi

usr_msg_id=$(uuidgen)

while true ; do

    curlresp=$(curl -s -o /dev/null -w "%{http_code}" \
                    -d "user_id=$user_id" -d "user_key=$user_key" -d "title=$title" -d "timestamp=$sendtime" \
                    -d "content=$content" -d "priority=$priority" -d "msg_id=$usr_msg_id" \
                    https://scn.blackforestbytes.com/send.php)
    
    if [ "$curlresp" == 200 ] ; then
        echo "Successfully send"
        exit 0
    fi

    if [ "$curlresp" == 400 ] ; then
        echo "Bad request - something went wrong"
        exit 1
    fi

    if [ "$curlresp" == 401 ] ; then
        echo "Unauthorized - wrong userid/userkey"
        exit 1
    fi

    if [ "$curlresp" == 403 ] ; then
        echo "Quota exceeded - wait one hour before re-try"
        sleep 3600
    fi

    if [ "$curlresp" == 412 ] ; then
        echo "Precondition Failed - No device linked"
        exit 1
    fi

    if [ "$curlresp" == 500 ] ; then
        echo "Internal server error - waiting for better times"
        sleep 60
    fi

    # if none of the above matched we probably hav no network ...
    echo "Send failed (response code $curlresp) ... try again in 5s"
    sleep 5
done