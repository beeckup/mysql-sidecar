#!/bin/sh


_now=$(date +"%s_%m_%d_%Y")

_name="$MYSQL_SQL_FILENAME_$_now"

_file_for_start="dumpdb/$_name"

_file="dumpdb/$_name.sql"

echo $_file


if [ "$MYSQL_ALL_DB" = "true" ]; then

    mysql -h "$MYSQL_DATABASE" -N -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e 'show databases' > todo.txt

    while read dbname; do
        mysqldump --single-transaction=TRUE -h "$MYSQL_DATABASE" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --complete-insert  "$dbname" > "$_file_for_start$dbname"_bck_`date +%Y%m%d`.sql;
        if [ "$ZIP_FILE" = "true" ]; then


          echo "Compress...";
          tar -cvzf "$_file_for_start$dbname"_bck_`date +%Y%m%d`.tar.gz "$_file_for_start$dbname"_bck_`date +%Y%m%d`.sql;
          rm "$_file_for_start$dbname"_bck_`date +%Y%m%d`.sql

          _file="$_file_for_start$dbname"_bck_`date +%Y%m%d`.tar.gz

          if [ "$S3_UPLOAD" = "true" and "$MYSQL_ALL_DB" = "true" ]; then



            echo "S3 upload..."

            bucket="$S3_BUCKET"

            host="$S3_HOST"
            link="$S3_PROTOCOL""://""$S3_HOST"

            echo "$link"

            s3_key="$S3_KEY"
            s3_secret="$S3_SECRET"

            resource="/${bucket}/${_file}"
            content_type="application/octet-stream"
            date=`date -R`
            _signature="PUT\n\n${content_type}\n${date}\n${resource}"
            signature=`echo -en ${_signature} | openssl sha1 -hmac ${s3_secret} -binary | base64`
            echo "Eseguo curl"
            echo "Upload di ${_file}"
            curl -v -X PUT -T "${_file}" \
                      -H "Host: $host" \
                      -H "Date: ${date}" \
                      -H "Content-Type: ${content_type}" \
                      -H "Authorization: AWS ${s3_key}:${signature}" \
                       $link${resource}

            rm $_file


          fi


        else

          _file="$_file_for_start$dbname"_bck_`date +%Y%m%d`.sql;

          if [ "$S3_UPLOAD" = "true" and "$MYSQL_ALL_DB" = "true" ]; then



            echo "S3 upload..."

            bucket="$S3_BUCKET"

            host="$S3_HOST"
            link="$S3_PROTOCOL""://""$S3_HOST"

            echo "$link"

            s3_key="$S3_KEY"
            s3_secret="$S3_SECRET"

            resource="/${bucket}/${_file}"
            content_type="application/octet-stream"
            date=`date -R`
            _signature="PUT\n\n${content_type}\n${date}\n${resource}"
            signature=`echo -en ${_signature} | openssl sha1 -hmac ${s3_secret} -binary | base64`
            echo "Eseguo curl"
            echo "Upload di ${_file}"
            curl -v -X PUT -T "${_file}" \
                      -H "Host: $host" \
                      -H "Date: ${date}" \
                      -H "Content-Type: ${content_type}" \
                      -H "Authorization: AWS ${s3_key}:${signature}" \
                       $link${resource}

            rm $_file


          fi

        fi

    done < todo.txt


fi



if [ "$MYSQL_ALL_DB" = "" ]; then

    mysqldump  "$MYSQL_DATABASE" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -h "$MYSQL_HOST" > $_file

    if [ "$ZIP_FILE" = "true" ]; then
        echo "Compress..."
        tar -cvzf dumpdb/$MYSQL_SQL_FILENAME_$_now.tar.gz $_file
        rm $_file
        _file=dumpdb/$MYSQL_SQL_FILENAME_$_now.tar.gz

    fi


    if [ "$S3_UPLOAD" = "true" and "$MYSQL_ALL_DB" = "" ]; then



        echo "S3 upload..."

        bucket="$S3_BUCKET"

        host="$S3_HOST"
        link="$S3_PROTOCOL""://""$S3_HOST"

        echo "$link"

        s3_key="$S3_KEY"
        s3_secret="$S3_SECRET"

        resource="/${bucket}/${_file}"
        content_type="application/octet-stream"
        date=`date -R`
        _signature="PUT\n\n${content_type}\n${date}\n${resource}"
        signature=`echo -en ${_signature} | openssl sha1 -hmac ${s3_secret} -binary | base64`
        echo "Eseguo curl"
        echo "Upload di ${_file}"
        curl -v -X PUT -T "${_file}" \
                  -H "Host: $host" \
                  -H "Date: ${date}" \
                  -H "Content-Type: ${content_type}" \
                  -H "Authorization: AWS ${s3_key}:${signature}" \
                   $link${resource}


        rm $_file

    fi

fi


