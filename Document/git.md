# ================================================== 安装配置 ==================================================
# 查看git版本
git -v
# 设置全局变量
git config --global user.name "Jasper Yang"
git config --global user.email geekhall@gamil.com
# 凭证存储
git config --global credential.helper store
# 列出所有全局变量
git config --global --list
# 移除某个全局变量
git config --global --unset user.name

# ================================================== 基本用法 ==================================================
# 初始化本地仓库
git init

# 推送到远程仓库
git push -u origin main

# 将本地项目连接到远程仓库
git remote add origin https://github.com/dongshan2025/WorkRelated.git

# 更改默认的远程仓库
git remote set-url origin <新的远程git仓库地址>

# 列出远程仓库名称
git remote

# 列出远程仓库名称和URL
git remote -v

# 删除远程仓库引用
git remote remove <远程仓库名称>  // 远程仓库名称可以由 git remote 获取到

# 显示指定远程仓库的详细信息
git remote show <远程仓库名称>

# 基于当前分支在本地创建新分支
git branch dev

# 切换到新分支
git switch dev

# 基于当前分支在本地创建新分支并切换到新分支
git checkout -b dev

# 推送到远程
git push origin dev 或 git push --set-upstream origin dev

# 显示本地所有分支
git branch
	
# 显示远程所有分支
git branch -r

# 显示本地和远程所有分支
git branch -a

# 切换到main分支
git switch main

# 将dev分支合并到main分支
git merge dev

# 删除本地dev分支
git branch -d dev

# 删除远程dev分支
git push origin :dev 或 git push origin --delete dev

# 显示本地分支和远程分支的对应关系
git branch -vv

# ================================================== 去掉Git历史提交记录 ==================================================
# 从当前分支创建新分支，并切换到新分支
git checkout --orphan new_branch
# 在新分支上添加所有文件
git add .
# 进行第一次提交
git commit -m "init commit"
# 删除旧分支
git branch -D master
# 将新分支重命名为旧分支名
git branch -m master
# 强制覆盖远程仓库的旧分支
git push -f origin master
# ================================================== 撤销对工作目录或暂存区文件的修改 ==================================================
# 注意：新增文件是未被跟踪状态(untracked)，也就是下面命令不会对此文件生效，新增文件要先被git记录到，也就是git add filename，之后执行这些操作才会生效

# 查看暂存文件状态
git status

# 丢弃工作区单个文件的修改
git restore filename

# 丢弃工作区所有文件的修改
git restore .

# 将暂存区的一个文件重新放回工作区，不改变修改的内容
git restore --staged filename

# 将暂存区的所有文件重新放回工作区，不改变修改的内容
git restore --staged .

# 将工作区中的文件恢复到以前的状态，如果--worktree单独指定，则默认恢复源为HEAD
git restore --worktree filename

# 将暂存区中的文件恢复到以前的状态
git restore --staged --worktree filename

# 从指定提交恢复文件
git restore --source=HEAD~1 filename

# 合并冲突时，恢复为当前分支的版本(即"我们"的版本)
git restore --ours filename

# 合并冲突时，恢复为另一个分支的版本(即"他们"的版本)
git restore --theirs filename

# ================================================== Reset用法 ==================================================
git reset --soft            保留工作区的修改        保留暂存区的修改
git reset --hard            不保留工作区的修改      不保留暂存区的修改
git reset --mixed           保留工作区的修改        不保留暂存区的修改
![alt text](<git reset.png>)


# ================================================== 查看历史提交记录 ==================================================
# 查看本地仓库的历史提交记录
git log

# 查看本地仓库的简要提交记录
git log --oneline

# 查看远程仓库的历史提交记录
## 更新远程引用
git fetch origin // 这不会自动合并到你的本地分支
git log origin/main // 查看远程仓库main分支的历史提交记录
git log origin/main --oneline

