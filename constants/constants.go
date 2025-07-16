package constants

import (
	"path"
	"time"
)

var (
	VerificationCodeTTL = 5 * time.Minute
)

var (
	ALL_API_PREFIX     = "api/v1"
	REPO_API_PREFIX    = "/repository"
	USER_API_PREFIX    = "/users"
	LOGIN_API_PREFIX   = "/login"
	SSH_KEY_API_PROFIX = "/ssh"
)

// model
var (
	SSH_KEY_TABLE_NAME    = "t_ssh_key"
	USER_TABLE_NAME       = "t_users"
	REPOSITORY_TABLE_NAME = "t_repository"
)

var (
	GIT_PRIVATE_KEY             = "~/gitolite-admin"
	GIT_SERVER_USERNAME         = "git"
	GIT_SERVER_HOME             = "/home/git"
	GIT_SERVER_DOMAIN           = ""
	GIT_SERVER_PORT             = ""
	GIT_SERVER_ADMIN_REPOSITORY = ""

	GITOLITE_ADMIN_REPOSITORY_PATH    = "~/gitolite-admin"
	GITOLITE_CONF_DIR_PATH            = path.Join(GITOLITE_ADMIN_REPOSITORY_PATH, "conf")
	GITOLITE_CONF_FILE_PATH           = path.Join(GITOLITE_CONF_DIR_PATH, "gitolite.conf")
	GITOLITE_USER_CONF_DIR_PATH       = path.Join(GITOLITE_CONF_DIR_PATH, "gitolite.d", "user")
	GITOLITE_REPOSITORY_CONF_DIR_PATH = path.Join(GITOLITE_CONF_DIR_PATH, "gitolite.d", "repository")
	GITOLITE_KEY_DIR_PATH             = path.Join(GITOLITE_ADMIN_REPOSITORY_PATH, "keydir")
)
