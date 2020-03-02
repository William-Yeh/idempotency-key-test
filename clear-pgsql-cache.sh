#!/bin/bash

echo "Clearing PostgreSQL cache and restarting..."

case "$OSTYPE" in
  darwin*)  echo "OSX"
            # @see https://stackoverflow.com/a/29277251/714426
            pg_ctl -D /usr/local/var/postgres stop
            sync
            sudo purge
            pg_ctl -D /usr/local/var/postgres start
            ;;

  linux*)   echo "LINUX"
            # @see https://stackoverflow.com/a/32661085/714426
            service postgresql stop
            sync
            echo 3 > /proc/sys/vm/drop_caches
            service postgresql start
            ;;

  msys*)    echo "WINDOWS"
            echo "Sorry, not supported."
            ;;

  *)        echo "Unknown: $OSTYPE" ;;
esac

