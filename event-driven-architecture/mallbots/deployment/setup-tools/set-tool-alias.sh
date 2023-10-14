#!/usr/bin/env sh

# source set-tool-alias.sh
# Once you have your alias, you can verify it works by running the following:
# deploytools terraform -version
#
# When you are using this option, you need to prefix the commands in the following sections with the
# deploytools command. Let’s take this command as an example:
# aws configure
# Turn it into this command:
# deploytools aws configure

root=$(realpath ${PWD}/../..)

docker build -t deploytools:latest .

alias deploytools='docker run --rm -it -v /var/run/docker.sock:/var/run/docker.sock -v ~/.aws:/root/.aws -v ~/.docker:/root/.docker -v ~/.kube:/root/.kube -v ${root}:/mallbots -v ${PWD}:/mallbots/deployment/.current -w /mallbots/deployment/.current deploytools'

echo "---"
echo
echo "Usage: deploytools <cmd [parameters]>"
