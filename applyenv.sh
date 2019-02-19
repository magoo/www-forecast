# exports the vars in .env into your shell
# `. ./applyenv.sh` to affect current shell session's env vars
export $(egrep -v '^#' .env | xargs)
