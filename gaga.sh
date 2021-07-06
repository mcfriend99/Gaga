#!/bin/bash


build() {
  if [[ $1 == "" ]]
  then
    NAME="gaga"
  else
    NAME="$1"
  fi

  rm -rf "build"
  mkdir -p "build"
  go get .
  go build -o "build/$NAME"

  echo "Successfully built $NAME into build/$NAME"
}

serve() {
  build "$1"
  "$(pwd)/build/$NAME"
}

clean() {
  if [[ $1 == "cache" ]]
  then
    rm -rf data/logs/*.log
  elif [[ $1 == "logs" ]]
  then
    rm -rf data/logs/*.log
  else
    rm -rf data/logs/*.log
  fi

  echo "Cleaning successful!"
}

help() {
  echo "Usage: gaga [action]"
  echo "  When action is not specified, action=help"
  echo "actions:"
  echo "  - build:  Builds the application"
  echo "            You may pass the name of the output executable as an argument."
  echo "            [default=gaga]"
  echo "  - serve:  Starts the server"
  echo "            You may pass the name of the output executable of the build"
  echo "            process as an argument."
  echo "            [default=gaga]"
  echo "  - clean:  Clean the gaga cache and log files."
  echo "            You may specify which item to clean as below:"
  echo "                > logs: clean logs only"
  echo "                > cache: clean cache only"
  echo "             e.g. gaga clean logs"
  echo "  - help:   Show this help message"
}

ACTION="$1"

case $ACTION in

  build)
    build "$2"
    ;;

  serve)
    serve "$2"
    ;;

  clean)
    clean "$2"
    ;;

  *)
    help
    ;;
esac