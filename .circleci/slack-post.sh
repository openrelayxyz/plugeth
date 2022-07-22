#!/usr/bin/env bash

# Usage: slackpost -w <webhook_url> -c <channel> -u <username> -m <message> [-a <alert_type>]
# Usage: echo <message> | slackpost -w <webhook_url> -c <channel> -u <username> [-a <alert_type>]

# exit immediately if a command exits with a non-zero status
set -e

# error if variable referenced before being set
set -u

# produce failure return code if any command fails in pipe
set -o pipefail

# accepted values: good, warning, danger
alert_type=""
channel=""
message=""
username=""
webhook_url=""

# colon after var means it has a value rather than it being a bool flag
while getopts 'a:c:m:u:w:' OPTION; do
  case "$OPTION" in
    a)
      alert_type="$OPTARG"
      ;;
    c)
      channel="$OPTARG"
      ;;
    m)
      message="$OPTARG"
      ;;
    u)
      username="$OPTARG"
      ;;
    w)
      webhook_url="$OPTARG"
      ;;
    ?)
      echo "script usage: $(basename $0) {-c channel} {-m message} {-u username} {-w webhook} [-a alert_type]" >&2
      exit 1
      ;;
  esac
done
shift "$(($OPTIND -1))"

# # exit if channel not provided
# if [[ -z "$channel" ]]
# then
#   echo "No channel specified"
#   exit 1
# fi

# read piped data as message if message argument is not provided
if [[ -z "$message" ]]
then
  message=$*

  while IFS= read -r line; do
    message="$message$line\n"
  done
fi

# # exit if username not provided
# if [[ -z "$username" ]]
# then
#   echo "No username specified"
#   exit 1
# fi

# exit if webhook not provided
if [[ -z "$webhook_url" ]]
then
  echo "No webhook_url specified"
  exit 1
fi

# escape message text
escapedText=$(echo $message | sed 's/"/\"/g' | sed "s/'/\'/g")

# create JSON payload
# json="{\"channel\": \"$channel\", \"username\":\"$username\", \"icon_emoji\":\"ghost\", \"attachments\":[{\"color\":\"$alert_type\" , \"text\": \"$escapedText\"}]}"
json="{\"attachments\":[{\"color\":\"$alert_type\" , \"text\": \"$escapedText\"}]}"

# fire off slack message post
curl -s -d "payload=$json" "$webhook_url"
