set -e

display_commit_message_error() {
  if [[ -n "${CI}" ]]; then
    echo "::error title=Commit message check failed::$2"
    echo -e "::group::Commit message\n$1\n::endgroup::"
  else
    cat <<EndOfMessage
$1

-------------------------------------------------
EndOfMessage
  fi

  cat <<EndOfMessage
The preceding commit message is invalid
it failed '$2' of the following checks

* Separate subject from body with a blank line
* Limit the subject line to 50 characters
* Capitalize the subject line
* Do not end the subject line with a period
* Wrap the body at 72 characters
EndOfMessage

  exit 1
}

lint_commit_message() {
  SECOND_LINE=$(echo "$1" | awk 'NR == 2')

  if [[ -n $SECOND_LINE ]]; then
    display_commit_message_error "$1" 'Separate subject from body with a blank line'
  fi

  if [[ "$(echo "$1" | head -n1 | awk '{print length}')" -gt 50 ]]; then
    re='^Update module [0-9a-zA-Z./]+ to v[0-9]+\.[0-9]+\.[0-9]+( \[.*\])?$'
    if [[ "$(echo "$1" | head -n1)" =~ $re ]]; then
      echo "Ignored subject line length error for module update commit"
    else
      display_commit_message_error "$1" 'Limit the subject line to 50 characters'
    fi
  fi

  if [[ ! $1 =~ ^[A-Z] ]]; then
    display_commit_message_error "$1" 'Capitalize the subject line'
  fi

  if [[ "$(echo "$1" | awk 'NR == 1 {print substr($0,length($0),1)}')" == "." ]]; then
    display_commit_message_error "$1" 'Do not end the subject line with a period'
  fi

  if [[ "$(echo "$1" | awk '{print length}' | sort -nr | head -1)" -gt 72 ]]; then
    display_commit_message_error "$1" 'Wrap the body at 72 characters'
  fi
}

# When a single argument is passed to the script ("$#" -eq 1):
# The script expects the argument to be a file containing a commit message.
if [ "$#" -eq 1 ]; then
  if [ ! -f "$1" ]; then
    echo "$0 was passed one argument, but was not a valid file"
    exit 1
  fi
  # extracts the commit message from the file using sed
  # The -n flag suppresses automatic printing of lines.
  # /pattern/q;p :
  #   q tells sed to quit processing when it encounters this pattern.
  #   p prints all lines before the pattern is matched.
  # In this case, sed extracts all lines from the file ("$1") up to but not
  # including the line containing the pattern. This is typically used to extract
  # the commit message from a Git commit template file.
  lint_commit_message "$(sed -n '/# Please enter the commit message for your changes. Lines starting/q;p' "$1")"
else
  # iterates over a range of commits in the Git repository using git rev-list.
  for COMMIT in $(git rev-list --no-merges origin/${GITHUB_BASE_REF:-master}..); do
    # For each commit, extract the commit message using git log (-n is max num of commits)
    # 1- %B means the raw commit message body.
    # 2- --no-merges excludes merge commits from the output.
    # 3- The .. syntax specifies a range of commits. In this case, it lists all commits that
    #    are in the current branch but not in the base branch (origin/${GITHUB_BASE_REF:-master}).
    # 4- finally This command extracts the full commit message (subject and body)
    #    for the given commit.
    lint_commit_message "$(git log --format="%B" -n 1 ${COMMIT})"
  done
fi
