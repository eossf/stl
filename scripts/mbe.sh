# markdown bash execute
# 1. add in your: .bashrc
# 2. in a bash console run: mbe readme.md
# 3. or launch directly
function mbe() {
  if [[ -f "$1" ]]; then
    cat $1 | sed -n '/````bash/,/````/p' | sed 's/````bash//g' | sed 's/````//g' | sed '/^$/d' | sed 's/####//g' | sed 's/###//g' | sed 's/##//g' | sed 's/#//g' | /usr/bin/env bash ;
  else
    echo "${1} is not valid" ;
  fi
}

filetoread=$1
if [[ $filetoread == "" ]] ; then
	exit;
else
  mbe $filetoread
fi
