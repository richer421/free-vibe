package scaffold

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func loadRegistry(projectRoot string) (Registry, error) {
	path := filepath.Join(projectRoot, RegistryFileName)
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Registry{Version: 1, Modules: []Module{}}, nil
		}
		return Registry{}, fmt.Errorf("read registry: %w", err)
	}
	var reg Registry
	if err := yaml.Unmarshal(b, &reg); err != nil {
		return Registry{}, fmt.Errorf("unmarshal registry: %w", err)
	}
	if reg.Version == 0 {
		reg.Version = 1
	}
	if reg.Modules == nil {
		reg.Modules = []Module{}
	}
	return reg, nil
}

func saveRegistry(projectRoot string, reg Registry) error {
	b, err := yaml.Marshal(reg)
	if err != nil {
		return fmt.Errorf("marshal registry: %w", err)
	}
	path := filepath.Join(projectRoot, RegistryFileName)
	if err := os.WriteFile(path, b, 0o644); err != nil {
		return fmt.Errorf("write registry: %w", err)
	}
	return nil
}

func ensureProjectReadme(projectRoot, projectName string) error {
	path := filepath.Join(projectRoot, "README.md")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	content := fmt.Sprintf("# %s\n\nManaged by FreeVibe CLI.\n\n- List templates: `freevibe template ls`\n- Add module: `freevibe add --template <name> --repo <url>`\n- Remove module: `freevibe remove <module>`\n- Sync modules: `git submodule update --init --recursive`\n", projectName)
	return os.WriteFile(path, []byte(content), 0o644)
}

func ensureProjectKnowledge(projectRoot string) error {
	dir := filepath.Join(projectRoot, "knowledge")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create knowledge dir: %w", err)
	}

	path := filepath.Join(dir, "README.md")
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	content := `# Knowledge

这是一个基于 submodule 形式管理的微服务 AI vibe coding 项目。
根级 knowledge 用于沉淀父项目层面的业务背景、服务边界、跨模块协作约束和公共术语。

建议优先在这里维护：
1. 项目级业务背景与目标
2. 服务边界、模块职责与 owner 约定
3. 跨模块接口、联调顺序与依赖约束
4. 需要多个子模块共享的业务词汇与工作定义

使用规则：
1. 涉及跨模块协作或父项目层面的业务语义时，先阅读根级 knowledge。
2. 进入具体子模块实现后，再结合子模块自己的 ` + "`knowledge/`、`AGENTS.md` 与 `.codex`" + ` 继续处理。
3. 如果根级 knowledge 与子模块知识不一致，先指出冲突，再给推荐修正方案。
`

	return os.WriteFile(path, []byte(content), 0o644)
}

func ensureProjectCodex(projectRoot string) error {
	codexDir := filepath.Join(projectRoot, ".codex")
	skillsDir := filepath.Join(codexDir, "skills", "codex-submodule-worktree-best-practices")
	if err := os.MkdirAll(skillsDir, 0o755); err != nil {
		return fmt.Errorf("create codex skills dir: %w", err)
	}

	readmePath := filepath.Join(codexDir, "README.md")
	if _, err := os.Stat(readmePath); err != nil {
		readmeContent := `# Codex Local Notes

Place project-local Codex skills, prompts, and workflow notes in this directory.

Current local skills:

- ` + "`.codex/skills/codex-submodule-worktree-best-practices/SKILL.md`" + `
`
		if err := os.WriteFile(readmePath, []byte(readmeContent), 0o644); err != nil {
			return fmt.Errorf("write codex README: %w", err)
		}
	}

	skillPath := filepath.Join(skillsDir, "SKILL.md")
	if _, err := os.Stat(skillPath); err == nil {
		return nil
	}

	skillContent := `---
name: codex-submodule-worktree-best-practices
description: Use when working in this project or a scaffolded parent repo with Codex, the repository uses git submodules, or any task involves worktree branches, branch switching, commit/push, merging main, or the project keyword “接收”.
---

# Codex Submodule Worktree Best Practices

## 优先级

这是本项目的高优先级流程 skill。

只要任务涉及以下任一场景，就必须优先使用本 skill，而不是按通用 git / worktree 习惯自行判断：

1. 在本项目中进行分支创建、切换、提交、推送、合并。
2. 在 worktree 中开发，或需要切回主项目目录执行 main 合并。
3. 用户明确说“接收”。
4. 处理子模块指针更新、子模块 main 合并、主项目回填。

项目关键词约定：

- **接收**：表示“提交当前工作分支改动、推送当前分支、按本 skill 切换到主项目目录合并 ` + "`main`" + `、推送 ` + "`origin/main`" + `，并完成最终状态校验”。
- 以后用户只要说“接收”，默认按上面的完整流程执行；除非用户明确缩小范围。

执行要求：

1. 命中上述场景时，必须先说明“正在按 ` + "`codex-submodule-worktree-best-practices`" + ` 执行”。
2. 必须先说明当前所在的是“worktree 项目”还是“主项目”。
3. 若要合并 ` + "`main`" + `，必须明确说明为什么要切到主项目目录执行。
4. 若本 skill 与通用习惯冲突，以本 skill 为准；除非用户明确要求覆盖。

## 基本概念

1. 主项目：真实开发项目（非 worktree 常驻目录）。
2. worktree 项目：从主项目切出的工作副本，用于某个需求的隔离开发。
3. 子模块：主项目中通过 git submodule 管理的模块集合，以主仓 ` + "`.gitmodules`" + ` 为准，不硬编码模块名单。
4. 非子模块目录：主仓中的普通目录，不走 submodule 指针提交流程。

## 规则

1. 分支统一：
主仓库和所有参与开发的子模块，必须使用同一个需求分支名，例如 ` + "`feat/xxx需求-20250320`" + `。
2. 分支来源统一：
无论主仓库还是子模块，创建需求分支前都必须先同步并基于各自 ` + "`origin/main`" + ` 切出。
3. 新增模块同规则：
开发中途新增子模块时，也必须先从该模块 ` + "`origin/main`" + ` 切出同名需求分支。
4. 提交/推送与合并分离：
用户只说“提交、推送”时，只提交并推送当前需求分支，不主动合并 ` + "`main`" + `。
5. 合并由用户触发：
只有用户明确说“合并 main 分支”或“接收”时，才执行主仓库和相关子模块的 ` + "`main`" + ` 合并流程。
6. 合并后必须双更新：
合并完成后，不仅更新主仓库子模块指针，还要更新主项目中的子模块工作副本到最新 ` + "`main`" + `。
7. 执行前置说明必做：
凡是命中本 skill 的任务，开始执行前都必须先说明当前目录角色（主项目 / worktree 项目）、目标分支、以及是否会切换到主项目目录处理 ` + "`main`" + `。

## 流程

1. 初始化需求分支：
主仓库和目标子模块先 ` + "`fetch/pull origin main`" + `，再切同名需求分支。
2. 需求开发：
只在需求分支上开发；中途新增模块时按同样规则补齐同名分支。
3. 提交与推送：
按用户指令提交并推送需求分支，不自动执行 main 合并。
4. 合并 main / 接收（仅用户明确要求时）：
先合并并推送各子模块 ` + "`main`" + `，再回主仓库更新并提交子模块指针到主仓 ` + "`main`" + `。
5. 主项目回填：
回到非 worktree 主项目目录，执行 ` + "`pull + submodule sync/update`" + `，确保主项目与远端一致。
6. 结束校验：
以主项目（非 worktree）` + "`status` 与 `submodule status`" + ` 作为最终完成依据。
`

	if err := os.WriteFile(skillPath, []byte(skillContent), 0o644); err != nil {
		return fmt.Errorf("write codex skill: %w", err)
	}

	return nil
}

func ensureProjectAgents(projectRoot string) error {
	path := filepath.Join(projectRoot, "AGENTS.md")
	if _, err := os.Stat(path); err == nil {
		return nil
	}

	content := `# AGENTS

## 仓库定位

这是一个基于 submodule 形式管理的微服务 AI vibe coding 项目。
这是父项目根仓库，不承载具体业务模块代码。
它负责管理 git submodule、模块注册表、根级协作约束以及跨模块变更。

## Skill 使用要求

1. 只要任务涉及 git submodule、worktree、分支创建/切换、提交、推送、合并，或用户明确说“接收”，必须先使用本地 skill：` + "`.codex/skills/codex-submodule-worktree-best-practices/SKILL.md`" + `
2. 若该 skill 与通用 git / worktree 习惯冲突，以该 skill 为准；若与用户最新指令冲突，以用户最新指令为准。

## 工作入口规则

1. 开始处理任务前，先查看 ` + "`freevibe.modules.yaml` 和 `.gitmodules`" + `，确认模块清单、路径、仓库来源和影响范围。
2. 先判断当前任务属于哪一层：
   - 父仓库编排：模块增删、submodule 同步、根目录文档、CI、发布、协作规范。
   - 子模块实现：某个业务模块、前端模块或服务模块内部代码与配置。
   - 跨模块联动：涉及多个子模块的接口、依赖、交付顺序或联调约束。
3. 如果是子模块实现，进入对应子模块目录，并遵循该子模块自己的 ` + "`AGENTS.md` 与本地 `.codex` 约束" + `。
4. 如果是跨模块任务，先明确受影响模块、边界和依赖顺序，再分别在对应模块内落地，不要把实现混写在父仓库根目录。
5. 不要在父仓库根目录直接承载本该属于子模块的业务逻辑、页面、服务实现或数据库变更。
6. 涉及项目级业务背景、服务边界、跨模块术语或协作规则时，优先查看 ` + "`knowledge/`" + `。

## 父仓库允许承载的内容

1. ` + "`freevibe.modules.yaml`、`.gitmodules`、根 `Makefile`" + `
2. 根级 README、协作约定、开发规范、CI/发布脚本
3. ` + "`knowledge/`" + ` 下的项目级知识文档
4. ` + "`.codex/`" + ` 下的父项目本地技能与协作提示
5. 管理 submodule 与跨模块协作所需的根级配置

## 执行要求

1. 优先做最小且正确的改动，避免把模块内逻辑错误地上提到父仓库。
2. 涉及具体业务语义、页面、接口、数据库、服务实现时，应下钻到对应子模块处理。
3. 所有结论必须可验证：能运行、能编译、能测试、能复现。
4. 遇到不确定性时，先给 1-2 个可行方案并明确推荐，不把选择题原样甩给用户。
`

	return os.WriteFile(path, []byte(content), 0o644)
}

func generateRootMakefile(projectRoot string) error {
	content := `.PHONY: modules status pull

modules:
	@cat freevibe.modules.yaml

status:
	@git submodule status

pull:
	@git submodule update --init --recursive
`
	path := filepath.Join(projectRoot, "Makefile")
	return os.WriteFile(path, []byte(content), 0o644)
}
