touch .edit; MISSING=$(find . -type f -name \*.go|grep -v -f .edit); [ -z "${MISSING}" ] || echo "${MISSING}" >> .edit

execute(){
	rm -f golib
	go mod tidy
	BDEBUG=1 go test -v
}
updatelibs(){
	go get -u github.com/rodolfoap/gx@$(cat ~/git/gx/VERSION)
}
tagversion(){
	# Always increase VERSION
	NEWVERS=$(cat VERSION|awk -F. '{print $1"."$2"."$3+1}')
	echo Current version is: $(cat VERSION)
	echo New version will be: ${NEWVERS}
	read -p "Tag message: " TAGMESSAGE
	echo ${NEWVERS}>VERSION

	# Always commit
	git add .;
	git commit -m "${TAGMESSAGE}"
	git push

	# Tag
	git tag $(cat VERSION)
	git push origin $(cat VERSION)
}

case "$1" in
t)	tagversion;
	;;
e) 	vi -p $(grep -v '^#' .edit) .edit
	ls *.go|xargs -n1 goformat
	execute;
	;;
u)	updatelibs
	;;
"")	execute
	;;
esac
