#!/bin/bash

# Used in https://github.com/zaunerc/go-release-scripts/blob/master/createTarGzPackage.sh.

function assembleArtifact {
	mkdir -p "$ARCHIVE_ROOT/usr/local/bin"
	mkdir -p "$ARCHIVE_ROOT/usr/local/etc/cntrbrowserd"
	mkdir -p "$ARCHIVE_ROOT/usr/local/share/cntrbrowserd"

	cp ../cntrbrowserd "$ARCHIVE_ROOT/usr/local/bin"
	cp -r ../static_data/* "$ARCHIVE_ROOT/usr/local/share/cntrbrowserd"
}

