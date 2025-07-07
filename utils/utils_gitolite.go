package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"go-gcs/appError"
	"go-gcs/constants"
	"go.uber.org/zap"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GitoliteUtils struct {
}

func (r *GitoliteUtils) InitUserConfig(userId int64) error {
	userFileName := fmt.Sprintf("%d.conf", userId)
	userConfPath := constants.GITOLITE_CONF_DIR_PATH + userFileName

	if _, err := os.Stat(userConfPath); err == nil {
		zap.L().Error("file already exist when init user conf file")
		return appError.ErrorFileAlreadyExists
	}

	file, err := os.Create(userConfPath)
	if err != nil {
		zap.L().Error("create user config file failed")
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
	return path.Join(
		constants.GITOLITE_KEY_DIR_PATH,
		fmt.Sprintf("%d.pub", sshKeyId),
	), nil
}

func (r *GitoliteUtils) GetUserFilePath(userId int64) (string, error) {
	return path.Join(
		constants.GITOLITE_USER_CONF_DIR_PATH,
		fmt.Sprintf("%d.conf", userId),
	), nil
}

func (r *GitoliteUtils) GetRepoFilePath(repoId int64) (string, error) {
	return path.Join(
		constants.GITOLITE_REPOSITORY_CONF_DIR_PATH,
		fmt.Sprintf("%d.conf", repoId),
	), nil
}

func (r *GitoliteUtils) AddSshKey(sshKeyId int64, sshKey string, userId int64) error {
	sshKeyFilePath, _ := r.GetSshKeyFIlePath(sshKeyId)
	userFilePath, err := r.GetUserFilePath(userId)

	if _, err := os.Stat(sshKeyFilePath); err == nil {
		zap.L().Error("this id of SshKey file already exist")
		return appError.ErrorFileAlreadyExists
	}

	sshKeyFile, err := os.Create(sshKeyFilePath)
	if err != nil {
		zap.L().Error("create SshKey file failed")
		return appError.ErrorCreateFileFailed
	}
	defer sshKeyFile.Close()
	sshKeyFile.WriteString(sshKey)

	userFile, err := os.ReadFile(userFilePath)
	if err != nil {
		zap.L().Error("read user config file failed")
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
			zap.L().Error("when write updated user config file, error happens")
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
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		commitMessage,
		commitFiles,
	)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error: %s. Error get when run CommitAndPush.", err.Error()))
	}
	return nil
}

func (r *GitoliteUtils) RemoveSshKey(sshKeyId int64, userId int64) error {
	sshKeyFilePath, _ := r.GetSshKeyFIlePath(sshKeyId)
	userFilePath, _ := r.GetUserFilePath(userId)
	if _, err := os.Stat(sshKeyFilePath); err != nil {
		zap.L().Warn("Try to remove an non-existed sshKey")
		return err
	}
	err := os.Remove(sshKeyFilePath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error: %s", err.Error()))
		return err
	}
	userFile, err := os.ReadFile(userFilePath)
	if err != nil {
		zap.L().Error(err.Error())
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
			zap.L().Error(err.Error())
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
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		commitMessage,
		commitFiles,
	)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error: %s. Error get when run CommitAndPush.", err.Error()))
	}
	return nil
}

// not consider race condition
// data race, critical section, mutex, synchronization, atomic operation
func (r *GitoliteUtils) UpdateSshKey(sshKeyId int64, sshKey string) error {
	sshKeyFilePath, _ := r.GetSshKeyFIlePath(sshKeyId)
	if _, err := os.Stat(sshKeyFilePath); err != nil {
		zap.L().Error("no ssh key exists")
		return err
	}
	err := os.WriteFile(sshKeyFilePath, []byte(sshKey), 0644)
	if err != nil {
		zap.L().Error("update ssh key failed")
		return err
	}
	return nil
}

// there are serious concurrency problem
// 36
func (r *GitoliteUtils) CreateRepository(repoId int64, repoName string, isPrivate bool,
	userId int64, userName string) error {
	userFilePath, _ := r.GetUserFilePath(userId)
	repoFilePath, _ := r.GetRepoFilePath(repoId)
	if _, err := os.Stat(repoFilePath); err == nil {
		zap.L().Error("detected duplicate repository")
		return err
	}
	repoFile, err := os.Create(repoFilePath)
	if err != nil {
		zap.L().Error("create repository config file failed")
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
		zap.L().Error("write file failed when create repository")
		return err
	}
	userFile, err := os.ReadFile(userFilePath)
	if err != nil {
		zap.L().Error("get failed when read user config file")
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
	err = r.CommitAndPush(constants.GITOLITE_ADMIN_REPOSITORY_PATH, commitMessage, commitFiles)
	if err != nil {
		zap.L().Error("get err when commit and push")
		return err
	}
	return nil
}

// 15
func (r *GitoliteUtils) RemoveRepository(repoId int64, repoName string, isPrivate bool,
	userId int64, userName string) error {
	userFilePath, _ := r.GetUserFilePath(userId)
	repoFilePath, _ := r.GetRepoFilePath(repoId)
	if _, err := os.Stat(repoFilePath); err != nil {
		zap.L().Warn("to be removed repository don't exist")
		return nil
	}
	err := os.Remove(repoFilePath)
	if err != nil {
		zap.L().Error("remove repository failed")
		return err
	}
	userFile, err := os.ReadFile(userFilePath)
	if err != nil {
		zap.L().Error("can't open user config file")
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
		zap.L().Error("get error when writing to user config file")
		return err
	}
	return nil
}

func (r *GitoliteUtils) AddCollaborator(repoId int64, collaboratorId int64) error {
	repoFilePath, _ := r.GetRepoFilePath(repoId)
	if _, err := os.Stat(repoFilePath); os.IsNotExist(err) {
		zap.L().Error("file not exists")
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
			zap.L().Error("write file failed")
			return err
		}
	}
	commitMessage := fmt.Sprintf("Add collaborator %d", collaboratorId)
	file1, _ := filepath.Rel(
		constants.GITOLITE_ADMIN_REPOSITORY_PATH,
		repoFilePath,
	)
	err := r.CommitAndPush(repoFilePath, commitMessage, []string{file1})
	if err != nil {
		zap.L().Error("get error when commit and push")
		return err
	}
	return nil
}

func (r *GitoliteUtils) RemoveCollaborator(repoId int64, collaboratorId int64) error {
	repoFilePath, _ := r.GetRepoFilePath(repoId)
	if _, err := os.Stat(repoFilePath); os.IsNotExist(err) {
		zap.L().Error("repository config file not exists when Remove Collaborator")
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
			zap.L().Error("write file failed when remove collaborator")
			return err
		}
		file1, _ := filepath.Rel(
			constants.GITOLITE_ADMIN_REPOSITORY_PATH,
			repoFilePath,
		)
		commitMessage := fmt.Sprintf("Remove collaborator %d", collaboratorId)
		err = r.CommitAndPush(repoFilePath, commitMessage, []string{file1})
		if err != nil {
			zap.L().Error("get error when commit and push in remving collaborator")
			return err
		}

	}
	return nil
}

func (r *GitoliteUtils) CommitAndPush(repoPath string, message string, files []string) error {
	g, err := git.PlainOpen(repoPath)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Error: %s. Error get when Open Git Repository", err.Error()))
		return err
	}
	workdir, err := g.Worktree()
	if err != nil {
		zap.L().Error(
			fmt.Sprintf("Error: %s. Error get when get Repos workdir",
				err.Error(),
			),
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
		)
		return err
	}
	return nil
}
