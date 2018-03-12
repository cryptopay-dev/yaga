.PHONY=bump_version

Build=`git show -s --format=%ci`
Version=`git describe --tags --exact-match || git rev-parse HEAD`

bump_version:
	@echo "${Version}\n$(Build)" > VERSION