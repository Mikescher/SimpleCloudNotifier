#!/usr/bin/env bash

#
# Wrapper around SCN ( https://simplecloudnotifier.de/ )
# ======================================================
#
# ./scn_send [@channel] title [content] [priority]
#
#
# Call with   scn_send              "${title}"
#        or   scn_send              "${title}" ${content}"
#        or   scn_send              "${title}" ${content}" "${priority:0|1|2}"
#        or   scn_send "@${channel} "${title}"
#        or   scn_send "@${channel} "${title}" ${content}"
#        or   scn_send "@${channel} "${title}" ${content}" "${priority:0|1|2}"
#
# content can be of format "--scnsend-read-body-from-file={path}" to read body from file
# (this circumvents max commandline length)
#

################################################################################

usage() {
    echo "Usage: "
    echo "  scn_send [@channel] title [content] [priority]"
    echo ""
}

function cfgcol { [ -t 1 ] && [ -n "$(tput colors)" ] && [ "$(tput colors)" -ge 8 ]; }

function rederr() { if cfgcol; then >&2 echo -e "\x1B[31m$1\x1B[0m"; else >&2 echo "$1"; fi; }
function green()  { if cfgcol; then     echo -e "\x1B[32m$1\x1B[0m"; else     echo "$1"; fi; }

################################################################################

#
# Get env 'SCN_UID' and 'SCN_KEY' from conf file
# 
# shellcheck source=/dev/null
. "/etc/scn.conf"
SCN_UID=${SCN_UID:-}
SCN_KEY=${SCN_KEY:-}

[ -z "${SCN_UID}" ] && { rederr "Missing config value 'SCN_UID' in /etc/scn.conf"; exit 1; }
[ -z "${SCN_KEY}" ] && { rederr "Missing config value 'SCN_KEY' in /etc/scn.conf"; exit 1; }

################################################################################

args=( "$@" )

title=""
content=""
channel=""
priority=""
usr_msg_id="$(head /dev/urandom | tr -dc A-Za-z0-9 | head -c 32)"
sendtime="$(date +%s)"
sender="$(hostname)"

if command -v srvname &> /dev/null; then
  sender="$( srvname )"
fi

if [[ "${args[0]}" = "--" ]]; then
    # only positional args form here on (currently not handled)
    args=("${args[@]:1}")
fi

if [ ${#args[@]} -lt 1 ]; then
    rederr "[ERROR]: no title supplied via parameter"
    usage
    exit 1
fi

if [[ "${args[0]}" =~ ^@.* ]]; then
    channel="${args[0]}"
    args=("${args[@]:1}")
    channel="${channel:1}"
fi

if [ ${#args[@]} -lt 1 ]; then
    rederr "[ERROR]: no title supplied via parameter"
    usage
    exit 1
fi

title="${args[0]}"
args=("${args[@]:1}")

content=""

if [ ${#args[@]} -gt 0 ]; then
    content="${args[0]}"
    args=("${args[@]:1}")
fi

if [ ${#args[@]} -gt 0 ]; then
    priority="${args[0]}"
    args=("${args[@]:1}")
fi

if [ ${#args[@]} -gt 0 ]; then
    rederr "Too many arguments to scn_send"
    usage
    exit 1
fi

if [[ "$content" == --scnsend-read-body-from-file=* ]]; then
  path="$( awk '{ print substr($0, 31) }' <<< "$content" )"
  content="$( cat "$path" )"
fi

curlparams=()

curlparams+=( "--data-urlencode" "user_id=${SCN_UID}"  )
curlparams+=( "--data-urlencode" "key=${SCN_KEY}"      )
curlparams+=( "--data-urlencode" "title=$title"        )
curlparams+=( "--data-urlencode" "timestamp=$sendtime" )
curlparams+=( "--data-urlencode" "msg_id=$usr_msg_id"  )

if [[ -n "$content" ]]; then
    curlparams+=("--data-urlencode" "content=$content")
fi

if [[ -n "$priority" ]]; then
    curlparams+=("--data-urlencode" "priority=$priority")
fi

if [[ -n "$channel" ]]; then
    curlparams+=("--data-urlencode" "channel=$channel")
fi

if [[ -n "$sender" ]]; then
    curlparams+=("--data-urlencode" "sender_name=$sender")
fi

while true ; do

    outf="$(mktemp)"

    curlresp=$(curl --silent                             \
                    --output "${outf}"                   \
                    --write-out "%{http_code}"           \
                    "${curlparams[@]}"                   \
                    "https://simplecloudnotifier.de/"    )

    curlout="$(cat "$outf")"
    rm "$outf"

    if [ "$curlresp" == 200 ] ; then
        green "Successfully send"
        exit 0
    fi

    if [ "$curlresp" == 400 ] ; then
        rederr "Bad request - something went wrong"
        echo "$curlout"
        echo ""
        exit 1
    fi

    if [ "$curlresp" == 401 ] ; then
        rederr "Unauthorized - wrong userid/userkey"
        exit 1
    fi

    if [ "$curlresp" == 403 ] ; then
        rederr "Quota exceeded - wait 5 min before re-try"
        sleep 300
    fi

    if [ "$curlresp" == 412 ] ; then
        rederr "Precondition Failed - No device linked"
        exit 1
    fi

    if [ "$curlresp" == 500 ] ; then
        rederr "Internal server error - waiting for better times"
        sleep 60
    fi

    # if none of the above matched we probably have no network ...
    rederr "Send failed (response code $curlresp) ... try again in 5s"
    sleep 5
done
