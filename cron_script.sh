#!/bin/sh

echo "Process starting..."

_now=$(date +"%s_%m_%d_%Y")

_name=$MYSQL_SQL_FILENAME"_"$_now"_"

_file_for_start="dumpdb/$_name"

_file="dumpdb/$_name.sql"

echo $_file


if [ "$MYSQL_ALL_DB" = "true" ]; then

    echo "Dumping all dbs"

    mysql -h "$MYSQL_HOST" -N -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -e 'show databases' > todo.txt

    while read dbname; do
        mysqldump --skip_add_locks --skip-lock-tables -h "$MYSQL_HOST" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" --max_allowed_packet=3000M --complete-insert  "$dbname" > "$_file_for_start$dbname"_bck_`date +%Y%m%d`.sql;

        echo "dumping $dbname database..."

        if [ "$ZIP_FILE" = "true" ]; then


          echo "Compress $dbname database...";
          tar -cvzf "$_file_for_start$dbname"_bck_`date +%Y%m%d`.tar.gz "$_file_for_start$dbname"_bck_`date +%Y%m%d`.sql;
          rm "$_file_for_start$dbname"_bck_`date +%Y%m%d`.sql

          _file="$_file_for_start$dbname"_bck_`date +%Y%m%d`.tar.gz

          if [ "$S3_UPLOAD" = "true" ]; then



            echo "S3 upload $dbname database ($_file)..."

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

            echo "Removing temp file $dbname $_file"

            rm $_file


          fi


        else

          _file="$_file_for_start$dbname"_bck_`date +%Y%m%d`.sql;

          if [ "$S3_UPLOAD" = "true" ]; then



            echo "S3 upload $dbname database ($_file)..."

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

            echo "Removing temp file $dbname ($_file)"

            rm $_file


          fi

        fi

    done < todo.txt


fi



if [ "$MYSQL_ALL_DB" = "" ]; then

    echo "Dumping $MYSQL_DATABASE mysql database..."

    mysqldump --skip_add_locks --skip-lock-tables  "$MYSQL_DATABASE" -u"$MYSQL_USER" -p"$MYSQL_PASSWORD" -h "$MYSQL_HOST" --max_allowed_packet=3000M --complete-insert > $_file

    if [ "$ZIP_FILE" = "true" ]; then
        echo "Compress $MYSQL_DATABASE mysql database..."
        tar -cvzf "dumpdb/"$_name".tar.gz" $_file
        rm $_file
        _file="dumpdb/"$_name".tar.gz"

    fi


    if [ "$S3_UPLOAD" = "true" ] && [ "$MYSQL_ALL_DB" = "" ]; then



        echo "S3 upload $MYSQL_DATABASE database ($_file) ..."

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
        echo "Removing temp file of $MYSQL_DATABASE database ($_file)"

        rm $_file

    fi

fi


