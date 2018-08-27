#!/bin/bash

# Usage: ./cleaner "bucketname" "30 days" "directory"

aws configure set aws_access_key_id $S3_KEY
aws configure set aws_secret_access_key $S3_SECRET

aws --endpoint-url $S3_PROTOCOL://$S3_HOST/ s3 ls s3://$1/$3/ |  while read -r line;
  do
    createDate=`echo $line|awk {'print $1" "$2'}`
    #echo  $line
    #echo $createdDate
    createDate=`date -d"$createDate" +%s`
    olderThan=`date -d"-$2" +%s`
    if [[ $createDate -lt $olderThan ]]
      then
        fileName=`echo $line|awk {'print $4'}`
        if [[ $fileName != "" ]]
          then
            printf 'Deleting "%s"\n' $fileName
            aws  --endpoint-url $S3_PROTOCOL://$S3_HOST/ s3 rm s3://$1/$3/"$fileName"
        fi
    fi
  done;