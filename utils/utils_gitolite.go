package utils

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"os/exec"
	"bytes"

	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"go-gcs/appError"
	"go-gcs/constants"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"

	"github.com/go-git/go-git/v6"
	// "github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing/object"
)

type GitoliteUtils struct {
}

var (
	sshAuth *gitssh.PublicKeys
)

func init() {
	os.MkdirAll(expandHomeDir(constants.GITOLITE_USER_CONF_DIR_PATH), 0755)
	os.MkdirAll(expandHomeDir(constants.GITOLITE_REPOSITORY_CONF_DIR_PATH), 0755)
	os.MkdirAll(expandHomeDir(constants.GITOLITE_KEY_DIR_PATH), 0755)
	mustRunSSHAgent()
	mustAddKey()
	sshAuth, _ = newPublicKeyAuth("git", constants.GIT_PRIVATE_KEY)

}

func (r *GitoliteUtils) InitUserConfig(userId int64) error {
	userFileName := fmt.Sprintf("%d.conf", userId)
	userConfPath := expandHomeDir(path.Join(constants.GITOLITE_USER_CONF_DIR_PATH, userFileName))

	if _, err := os.Stat(userConfPath); err == nil {
		zap.L().Error("file already exist when init user conf file", zap.Error(err))
		return appError.ErrorFileAlreadyExists
	}

	file, err := os.Create(userConfPath)
	if err != nil {
		zap.L().Error("create user config file failed", zap.Error(err))
		return appError.ErrorCreateFileFailed
	}
	defer file.Close()
	content := fmt.Sprintf(`@%d_public_repo = testing
@%d_private_repo = testing
@%d_ssh_key = @admin
repo @%d_private_repo
    RW+ = @%d_ssh_key
repo @%d_public_repo
    RW+ = @%d_ssh_key
`, userId, userId, userId, userId, userId, userId, userId)

	_, err = file.WriteString(content)
	if err != nil {
		zap.L().Error(err.Error())
		return appError.ErrorWriteFileFailed
	}
	return nil
}

func (r *GitoliteUtils) GetSshKeyFIlePath(sshKeyId int64) (string, error) {
	return expandHomeDir(path.Join(
		constants.GITOLITE_KEY_DIR_PATH,
		fmt.Sprintf("%d.pub", sshKeyId),
	)), nil
}

func (r *GitoliteUtils) GetUserFilePath(userId int64) (string, error) {
	return expandHomeDir(path.Join(
		constants.GITOLITE_USER_CONF_DIR_PATH,
		fmt.Sprintf("%d.conf", userId),
	)), nil
}

func (r *GitoliteUtils) GetRepoFilePath(repoId int64) (string, error) {
	return expandHomeDir(path.Join(
		constants.GITOLITE_REPOSITORY_CONF_DIR_PATH,
		fmt.Sprintf("%d.conf", repoId),
	)), nil
}

func (r *GitoliteUtils) AddSshKey(ctx context.Context, sshKeyId int64, sshKey string, userId int64) error {
	sshKeyFilePath, _ := r.GetSshKeyFIlePath(sshKeyId)
	userFilePath, err := r.GetUserFilePath(userId)

	if _, err := os.Stat(sshKeyFilePath); err == nil {
		zap.L().Error("this id of SshKey file already exist", zap.Error(err))
		return appError.ErrorFileAlreadyExists
	}

	sshKeyFile, err := os.Create(sshKeyFilePath)
	if err != nil {
		zap.L().Error("create SshKey file failed", zap.Error(err))
		return appError.ErrorCreateFileFailed
	}
	defer sshKeyFile.Close()
	sshKeyFile.WriteString(sshKey)

	userFile, err := os.ReadFile(userFilePath)
	if err != nil {
		zap.L().Error("read user config file failed", zap.Error(err))
		return err
	}
	lines := strings.Split(string(userFile), "\n")
	updated := false
	for i, line := range lines {
		prefix := fmt.Sprintf("@%d_ssh_key", userId)
		if strings.HasPrefix(line, prefix) {
			lines[i] = line + " " + strconv.FormatInt(sshKeyId, 10)
			updated = true
			break
		}
	}

	if updated {
		newContent := strings.Join(lines, "\n")
		err = os.WriteFile(userFilePath, []byte(newContent), 0644)
		if err != nil {
			zap.L().Error("when write updated user config file, error happens", zap.Error(err))
			return err
		}
	}
	commitMessage := fmt.Sprintf("Add user %d", userId)
	var commitFiles []string
	file1, _ := filepath.Rel(
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		userFilePath,
	)
	file2, _ := filepath.Rel(
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		constants.GITOLITE_KEY_DIR_PATH,
	)
	commitFiles = append(commitFiles, file1)
	commitFiles = append(commitFiles, file2)
	err = r.CommitAndPush(
		ctx,
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		commitMessage,
		commitFiles,
	)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error: %s. Error get when run CommitAndPush.", err.Error()), zap.Error(err))
	}
	return nil
}

func (r *GitoliteUtils) RemoveSshKey(ctx context.Context, sshKeyId int64, userId int64) error {
	sshKeyFilePath, _ := r.GetSshKeyFIlePath(sshKeyId)
	userFilePath, _ := r.GetUserFilePath(userId)
	if _, err := os.Stat(sshKeyFilePath); err != nil {
		zap.L().Warn("Try to remove an non-existed sshKey", zap.Error(err))
		return err
	}
	err := os.Remove(sshKeyFilePath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error: %s", err.Error()), zap.Error(err))
		return err
	}
	userFile, err := os.ReadFile(userFilePath)
	if err != nil {
		zap.L().Error(err.Error(), zap.Error(err))
	}
	lines := strings.Split(string(userFile), "n")
	updated := false
	for i, line := range lines {
		if strings.HasPrefix(line, fmt.Sprintf("%d_ssh_key", sshKeyId)) {
			lines[i] = strings.Replace(lines[i], fmt.Sprintf(" %d", sshKeyId), "", 1)
			updated = true
			break
		}
	}
	if updated {
		var newContent = strings.Join(lines, "\n")
		err := os.WriteFile(userFilePath, []byte(newContent), 0644)
		if err != nil {
			zap.L().Error(err.Error(), zap.Error(err))
			return err
		}
	}

	commitMessage := fmt.Sprintf("Remove SshKey %d", sshKeyId)
	var commitFiles []string
	file1, _ := filepath.Rel(
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		userFilePath,
	)
	file2, _ := filepath.Rel(
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		constants.GITOLITE_KEY_DIR_PATH,
	)
	commitFiles = append(commitFiles, file1)
	commitFiles = append(commitFiles, file2)
	err = r.CommitAndPush(
		ctx,
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		commitMessage,
		commitFiles,
	)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error: %s. Error get when run CommitAndPush.", err.Error()), zap.Error(err))
	}
	return nil
}

// not consider race condition
// data race, critical section, mutex, synchronization, atomic operation
func (r *GitoliteUtils) UpdateSshKey(sshKeyId int64, sshKey string) error {
	sshKeyFilePath, _ := r.GetSshKeyFIlePath(sshKeyId)
	if _, err := os.Stat(sshKeyFilePath); err != nil {
		zap.L().Error("no ssh key exists", zap.Error(err))
		return err
	}
	err := os.WriteFile(sshKeyFilePath, []byte(sshKey), 0644)
	if err != nil {
		zap.L().Error("update ssh key failed", zap.Error(err))
		return err
	}
	return nil
}

// there are serious concurrency problem
// 36
func (r *GitoliteUtils) CreateRepository(ctx context.Context, repoId int64, repoName string, isPrivate bool,
	userId int64, userName string) error {
	userFilePath, _ := r.GetUserFilePath(userId)
	repoFilePath, _ := r.GetRepoFilePath(repoId)
	if _, err := os.Stat(repoFilePath); err == nil {
		zap.L().Error("detected duplicate repository", zap.Error(err))
		return err
	}
	repoFile, err := os.Create(repoFilePath)
	if err != nil {
		zap.L().Error("create repository config file failed", zap.Error(err))
		return err
	}
	content := fmt.Sprintf(`
@%d_repo_collaborator = @admin
repo %s/%s
	RW+ = @%d_repo_collaborator
`, repoId, userName, repoName, repoId)
	defer repoFile.Close()
	_, err = repoFile.WriteString(content)
	if err != nil {
		zap.L().Error("write file failed when create repository", zap.Error(err))
		return err
	}
	userFile, err := os.ReadFile(userFilePath)
	if err != nil {
		zap.L().Error("get failed when read user config file", zap.Error(err))
		return err
	}
	lines := strings.Split(string(userFile), "\n")
	updated := false
	for i, line := range lines {
		accessType := "pulic"
		if isPrivate {
			accessType = "private"
		}
		if strings.HasPrefix(line, fmt.Sprintf("@%d_%s_repo", userId, accessType)) {
			lines[i] = fmt.Sprintf("%s %s/%s", line, userName, repoName)
			updated = true
			break
		}
		// serval vars hold a file at the same time
	}
	content = strings.Join(lines, "\n")
	if updated {
		os.WriteFile(userFilePath, []byte(content), 0644)
	}
	// debug
	// zap.L().Info("commit was hacked successful")
	// return nil
	commitMessage := fmt.Sprintf("Create repository %s/%s", userName, repoName)
	file1, _ := filepath.Rel(
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		userFilePath,
	)
	file2, _ := filepath.Rel(
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		repoFilePath,
	)
	var commitFiles []string
	commitFiles = append(commitFiles, file1)
	commitFiles = append(commitFiles, file2)
	err = r.CommitAndPush(ctx, constants.GITOLITE_ADMIN_REPOSITORY_PATH, commitMessage, commitFiles)
	if err != nil {
		zap.L().Error("get err when commit and push", zap.Error(err))
		return err
	}
	return nil
}

// 15
func (r *GitoliteUtils) RemoveRepository(ctx context.Context, repoId int64, repoName string, isPrivate bool,
	userId int64, userName string) error {
	userFilePath, _ := r.GetUserFilePath(userId)
	repoFilePath, _ := r.GetRepoFilePath(repoId)
	if _, err := os.Stat(repoFilePath); err != nil {
		zap.L().Warn("to be removed repository don't exist", zap.Error(err))
		return nil
	}
	err := os.Remove(repoFilePath)
	if err != nil {
		zap.L().Error("remove repository failed", zap.Error(err))
		return err
	}
	userFile, err := os.ReadFile(userFilePath)
	if err != nil {
		zap.L().Error("can't open user config file", zap.Error(err))
		return nil
	}
	lines := strings.Split(string(userFile), "\n")
	accessType := "public"
	if isPrivate {
		accessType = "private"
	}
	for i, line := range lines {
		if strings.HasPrefix(line, fmt.Sprintf("@%d_%s_repo", userId, accessType)) {
			lines[i] = strings.Replace(line, fmt.Sprintf("%s/%s", userName, repoName), "", 1)
			break
		}
	}
	content := strings.Join(lines, "\n")
	err = os.WriteFile(userFilePath, []byte(content), 0644)
	if err != nil {
		zap.L().Error("get error when writing to user config file", zap.Error(err))
		return err
	}
	return nil
}

func (r *GitoliteUtils) AddCollaborator(ctx context.Context, repoId int64, collaboratorId int64) error {
	repoFilePath, _ := r.GetRepoFilePath(repoId)
	if _, err := os.Stat(repoFilePath); os.IsNotExist(err) {
		zap.L().Error("file not exists", zap.Error(err))
		return err
	}
	repoFile, _ := os.ReadFile(repoFilePath)
	lines := strings.Split(string(repoFile), "\n")
	updated := false
	for i, line := range lines {
		if strings.HasPrefix(line, fmt.Sprintf("@%d_repo_collaborator", repoId)) {
			lines[i] = line + fmt.Sprintf(" @%d_ssh_key", collaboratorId)
			updated = true
			break
		}
	}
	content := strings.Join(lines, "\n")
	if updated {
		err := os.WriteFile(repoFilePath, []byte(content), 0644)
		if err != nil {
			zap.L().Error("write file failed", zap.Error(err))
			return err
		}
	}
	commitMessage := fmt.Sprintf("Add collaborator %d", collaboratorId)
	file1, _ := filepath.Rel(
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		repoFilePath,
	)
	err := r.CommitAndPush(ctx, repoFilePath, commitMessage, []string{file1})
	if err != nil {
		zap.L().Error("get error when commit and push", zap.Error(err))
		return err
	}
	return nil
}

func (r *GitoliteUtils) RemoveCollaborator(ctx context.Context, repoId int64, collaboratorId int64) error {
	repoFilePath, _ := r.GetRepoFilePath(repoId)
	if _, err := os.Stat(repoFilePath); os.IsNotExist(err) {
		zap.L().Error("repository config file not exists when Remove Collaborator", zap.Error(err))
		return err
	}
	repoFile, _ := os.ReadFile(repoFilePath)
	lines := strings.Split(string(repoFile), "\n")
	updated := false
	for i, line := range lines {
		if strings.HasPrefix(line, fmt.Sprintf("@%d_repo_collaborator", repoId)) {
			lines[i] = strings.Replace(line, fmt.Sprintf(" @%d_ssh_key", collaboratorId), "", 1)
			updated = true
			break
		}
	}
	if updated {
		content := strings.Join(lines, "\n")
		err := os.WriteFile(repoFilePath, []byte(content), 0644)
		if err != nil {
			zap.L().Error("write file failed when remove collaborator", zap.Error(err))
			return err
		}
		file1, _ := filepath.Rel(
			constants.GITOLITE_ADMIN_REPOSITORY_PATH,
			repoFilePath,
		)
		commitMessage := fmt.Sprintf("Remove collaborator %d", collaboratorId)
		err = r.CommitAndPush(ctx, repoFilePath, commitMessage, []string{file1})
		if err != nil {
			zap.L().Error("get error when commit and push in remving collaborator", zap.Error(err))
			return err
		}

	}
	return nil
}

func (r *GitoliteUtils) CommitAndPush(ctx context.Context, repoPath string, message string, files []string) error {
	g, err := git.PlainOpen(repoPath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error: %s. Error get when Open Git Repository", err.Error()), zap.Error(err))
		return err
	}
	workdir, err := g.Worktree()
	if err != nil {
		zap.L().Error(
			fmt.Sprintf("Error: %s. Error get when get Repos workdir",
				err.Error(),
			),
			zap.Error(err),
		)
		return err
	}
	for _, file := range files {
		_, err = workdir.Add(file)
		if err != nil {
			zap.L().Error(
				fmt.Sprintf("Error: %s. Error get when add file",
					err.Error(),
				),
				zap.Error(err),
			)
			return err
		}
	}
	commit, err := workdir.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "admin",
			Email: "admin@localhost.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		zap.L().Error(
			fmt.Sprintf("Error: %s. Error get when commit",
				err.Error(),
			),
			zap.Error(err),
		)
		return err
	}
	obj, err := g.CommitObject(commit)
	if err == nil {
		zap.L().Info(
			fmt.Sprintf("Committed: %s", obj.Hash),
		)
	}

	err = g.Push(&git.PushOptions{})

	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			zap.L().Info("Already up-to-data")
			return nil
		}
		zap.L().Error(
			fmt.Sprintf("Error: %s. Error get when push",
				err.Error(),
			),
			zap.Error(err),
		)
		return err
	}
	return nil
}

func expandHomeDir(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func newPublicKeyAuth(user, privateKeyPath string) (*gitssh.PublicKeys, error) {
	key, err := os.ReadFile(expandHomeDir(privateKeyPath))
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}
	auth := &gitssh.PublicKeys{
		User:   user,
		Signer: signer,
	}
	return auth, nil
}

func mustRunSSHAgent() {
	cmd := exec.Command("ssh-agent", "-s")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "启动 ssh-agent 失败: %v\n", err)
		os.Exit(1)
	}

	// 解析 ssh-agent 输出，设置环境变量
	for _, line := range strings.Split(out.String(), "\n") {
		if strings.Contains(line, "SSH_AUTH_SOCK") || strings.Contains(line, "SSH_AGENT_PID") {
			parts := strings.SplitN(line, ";", 2)
			if kv := strings.SplitN(parts[0], "=", 2); len(kv) == 2 {
				os.Setenv(kv[0], kv[1])
			}
		}
	}
}

func mustAddKey() {
	cmd := exec.Command("ssh-add", os.ExpandEnv("$HOME/.ssh/id_rsa"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "添加私钥失败: %v\n", err)
		os.Exit(1)
	}
}
