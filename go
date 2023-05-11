#!/bin/bash

set -e

here="$( cd -P "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cmd=$1
shift

# Modify locations of files/folders here. Rest of script will stay relatively the same.
BIN_LOCATION="$here/../../../../bin/role-auth"

function build_webapp()
{
  cd rc-admin
  yarn install
  yarn build
  cd -
}

function verify_node()
{
  if [ ! -d "$here/node_modules" ]; then
    yarn install
  fi
}

case "$cmd" in

auth)
  AWS_CLI_PROFILE=${1:-default}

  # 1 arg required
  if [[ $# -ne 1 ]]; then
    echo "Usage: eval \$($0 auth <AWS_CLI_PROFILE>)"
    echo "Where:"
    echo "   <AWS_CLI_PROFILE> = aws-cli profile usually in $HOME/.aws/config"
    exit 2
  fi

  here="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
  "$BIN_LOCATION" $AWS_CLI_PROFILE
  ;;

clean)

  if ! command -v sls &> /dev/null
  then
      echo "COMMAND sls could not be found"
      exit
  fi

  sls requirements clean && sls wsgi clean
  ;;

install)
  verify_node
  ;;

install-dev)
  verify_node
  ;;

deploy)
  STAGE=${1:-default}
  if [[ "$1" == "" ]]; then
    STAGE='dev'
  fi

  export GIT_HASH=$(git rev-parse --short HEAD)
  echo "Deploying $GIT_HASH to... $STAGE"

  verify_node
  build_webapp
  yarn run deploy $STAGE
  echo "Deployed to admin--$STAGE.rcdevel.com."

esac