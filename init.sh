#!/bin/bash

cd ~/gitolite-admin || exit 1
echo "📁 Located in $(pwd)"

# 删除 user 下的 .conf 文件
user_confs=(conf/gitolite.d/user/*.conf)
if [[ -e "${user_confs[0]}" ]]; then
  rm -v "${user_confs[@]}"
else
  echo "⚠️  No .conf files found in user/"
fi

# 删除 repository 下的 .conf 文件
repo_confs=(conf/gitolite.d/repository/*.conf)
if [[ -e "${repo_confs[0]}" ]]; then
  rm -v "${repo_confs[@]}"
else
  echo "⚠️  No .conf files found in repository/"
fi

# 删除除 admin.pub 外的 .pub 文件
for f in keydir/*.pub; do
  [[ -f "$f" && "$(basename "$f")" != "admin.pub" ]] && rm -v "$f"
done

echo "pushing to remote..."
git add .
git commit -m "delete all user and repository config file"
git push --force

#!
