#!/bin/bash
if [ -d "$PROJECT_ROOT_DIR"/vendor ]; then
  chmod -R 755 "$PROJECT_ROOT_DIR"/vendor
  # chown -R $USER:$USER "$PROJECT_ROOT_DIR"/vendor
fi

chmod -R +x "$PROJECT_ROOT_DIR"/scripts
chmod -R 755 "$PROJECT_ROOT_DIR"/scripts

#if [ -d "$PROJECT_ROOT_DIR"/sandbox ]; then
#  chmod -R 755 "$PROJECT_ROOT_DIR"/sandbox
#fi
#
#if [ -d "$PROJECT_ROOT_DIR"/sandbox/majadbdata ]; then
#  chmod -R 775 "$PROJECT_ROOT_DIR"/sandbox/majadbdata
#fi


chmod -R 664 "$PROJECT_ROOT_DIR"/.gitconfig
