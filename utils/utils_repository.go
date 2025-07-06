package utils

import "fmt"

func GenerateHttpURL(host string, port int, username string,
	repoName string) string {
	return fmt.Sprintf("https://%s:%d/%s/%s.git", host, port, username, repoName)
}

func GenerateSshURL(host string, port int, username string,
	repoName string) string {
	return fmt.Sprintf("git@%s:%d/%s/%s.git", host, port, username, repoName)
}
