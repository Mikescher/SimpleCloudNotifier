#!/bin/bash

    # shellcheck disable=SC2002  # disable useless-cat warning

    set -o nounset   # disallow usage of unset vars  ( set -u )
    set -o errexit   # Exit immediately if a pipeline returns non-zero.  ( set -e )
    set -o errtrace  # Allow the above trap be inherited by all functions in the script.  ( set -E )
    set -o pipefail  # Return value of a pipeline is the value of the last (rightmost) command to exit with a non-zero status
    IFS=$'\n\t'      # Set $IFS to only newline and tab.

    # shellcheck disable=SC2034
    cr=$'\n'

    function black() { if [ -t 1 ] && [ -n "$(tput colors)" ] && [ "$(tput colors)" -ge 8 ]; then echo -e "\\x1B[30m$1\\x1B[0m"; else echo "$1"; fi  }
    function red()   { if [ -t 1 ] && [ -n "$(tput colors)" ] && [ "$(tput colors)" -ge 8 ]; then echo -e "\\x1B[31m$1\\x1B[0m"; else echo "$1"; fi; }
    function green() { if [ -t 1 ] && [ -n "$(tput colors)" ] && [ "$(tput colors)" -ge 8 ]; then echo -e "\\x1B[32m$1\\x1B[0m"; else echo "$1"; fi; }
    function yellow(){ if [ -t 1 ] && [ -n "$(tput colors)" ] && [ "$(tput colors)" -ge 8 ]; then echo -e "\\x1B[33m$1\\x1B[0m"; else echo "$1"; fi; }
    function blue()  { if [ -t 1 ] && [ -n "$(tput colors)" ] && [ "$(tput colors)" -ge 8 ]; then echo -e "\\x1B[34m$1\\x1B[0m"; else echo "$1"; fi; }
    function purple(){ if [ -t 1 ] && [ -n "$(tput colors)" ] && [ "$(tput colors)" -ge 8 ]; then echo -e "\\x1B[35m$1\\x1B[0m"; else echo "$1"; fi; }
    function cyan()  { if [ -t 1 ] && [ -n "$(tput colors)" ] && [ "$(tput colors)" -ge 8 ]; then echo -e "\\x1B[36m$1\\x1B[0m"; else echo "$1"; fi; }
    function white() { if [ -t 1 ] && [ -n "$(tput colors)" ] && [ "$(tput colors)" -ge 8 ]; then echo -e "\\x1B[37m$1\\x1B[0m"; else echo "$1"; fi; }

    # cd "$(dirname "$0")" || exit 1  # (optionally) cd to directory where script is located



pid="$( pgrep -f 'flutter_tools\.[s]napshot run' || echo '' )"

if [ -z "$pid" ]; then
    red "No [flutter run] process found - exiting"
    exit 1
fi

trap 'echo "reseived SIGNAL<EXIT> - exiting"; exit 0' EXIT
trap 'echo "reseived SIGNAL<SIGINT> - exiting"; exit 0' SIGINT
trap 'echo "reseived SIGNAL<SIGTERM> - exiting"; exit 0' SIGTERM
trap 'echo "reseived SIGNAL<SIGQUIT> - exiting"; exit 0' SIGQUIT

echo ""
blue "Listening for changes in lib/ directory - sending signals to ${pid}..."
echo ""

while true; do
    find lib/ -name '*.dart' | entr -d -p sh -c "echo 'File(s) changed - Sending SIGUSR to $pid' ; kill -USR1 $pid";
	yellow 'File list changed - restart';
done