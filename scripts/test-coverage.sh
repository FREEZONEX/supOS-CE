#!/usr/bin/env bash
# 企业版后端单元测试覆盖率检查（TASK-028 阶段 0 交付，04-单元测试方案 v1.1）
#
# 两个职责：
#   1. 覆盖率报告：对目标包跑 Go 原生语句覆盖率（statement coverage），
#      使用 -coverpkg 全包口径（无测试文件的包也计入分母），输出整体 total 与按包分布。
#   2. 文件级覆盖检查（--check-files）：遍历目标目录下每个非测试 .go 文件，
#      校验是否存在同名 <名>_test.go；缺口 > 0 即 exit 1（「文件级 100% 覆盖」验收工具）。
#
# 数据库：企业版测试唯一后端为 PostgreSQL（纯 PG 约定），连接串默认
#   postgres://test:test@localhost:5432/tier0_test?sslmode=disable
#   可用环境变量 TEST_DATABASE_URL 覆盖。
#
# 用法：
#   bash scripts/test-coverage.sh --check-files [目标目录...]   # 文件级覆盖检查
#   bash scripts/test-coverage.sh [包路径...]                   # 覆盖率报告（默认 ./internal/logic/... ./internal/domain/...）

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="${REPO_ROOT}/backend"
TEMP_DIR="${REPO_ROOT}/.temp"
mkdir -p "${TEMP_DIR}"
export TEST_DATABASE_URL="${TEST_DATABASE_URL:-postgres://test:test@localhost:5432/tier0_test?sslmode=disable}"

if [[ "${1:-}" == "--check-files" ]]; then
    shift
    DIRS=("$@")
    if [[ ${#DIRS[@]} -eq 0 ]]; then
        DIRS=(internal/logic/core internal/logic/openapi/v1 internal/domain)
    fi
    echo "================================================================"
    echo ">>> 文件级覆盖检查：${DIRS[*]}"
    echo "================================================================"
    missing=()
    total=0
    set +e
    while IFS= read -r -d '' f; do
        base="$(basename "$f" .go)"
        dir="$(dirname "$f")"
        # 无代码文件（仅 package 声明与注释）无可测内容，跳过
        code_count=$(grep -vE '^[[:space:]]*(//|/\*|\*|\*/|$)|^[[:space:]]*package[[:space:]]+' "$f" | grep -cE '.')
        if [[ "$code_count" -eq 0 ]]; then
            continue
        fi
        # //go:build ignore 文件不参与编译，跳过
        if grep -q '^//go:build ignore' "$f"; then
            continue
        fi
        total=$((total + 1))
        if [[ -f "${dir}/${base}_test.go" ]]; then
            continue
        fi
        missing+=("${f#${BACKEND_DIR}/}")
    done < <(find "${DIRS[@]/#/${BACKEND_DIR}/}" -name '*.go' ! -name '*_test.go' -print0)
    set -e
    if [[ ${#missing[@]} -gt 0 ]]; then
        echo "✗ 缺口 ${#missing[@]}/${total} 个文件无测试："
        printf '  %s\n' "${missing[@]}"
        exit 1
    fi
    echo "✓ 全部 ${total} 个文件均有对应测试文件"
    exit 0
fi

PKGS=("$@")
if [[ ${#PKGS[@]} -eq 0 ]]; then
    PKGS=(./internal/logic/... ./internal/domain/...)
fi
COVER_FILE="${TEMP_DIR}/coverage-enterprise.out"

echo "================================================================"
echo ">>> 覆盖率跑批：${PKGS[*]}"
echo "================================================================"
cd "${BACKEND_DIR}"
PKG_LIST="$(go list "${PKGS[@]}")"
COVERPKG="$(echo "${PKG_LIST}" | tr '\n' ',')"
go test -coverprofile="${COVER_FILE}" -coverpkg="${COVERPKG}" "${PKGS[@]}"

echo ""
echo "================================================================"
echo ">>> 覆盖率报告（语句覆盖率，-coverpkg 全包口径）"
echo "================================================================"
go tool cover -func="${COVER_FILE}" > "${TEMP_DIR}/cover-func-enterprise.txt"
grep '^total:' "${TEMP_DIR}/cover-func-enterprise.txt" || true
echo ""
echo "--- 低覆盖文件 TOP10 ---"
awk '/^total:/{next} /\.go:/{
    file = $1; sub(/:[0-9]+:$/, "", file)
    sub(/\(/, "", $3); sub(/%/, "", $3)
    sum[file] += $3; cnt[file] += 1
}
END { for (f in sum) printf "%6.1f%%  %s\n", sum[f]/cnt[f], f
}' "${TEMP_DIR}/cover-func-enterprise.txt" | sort -n | head -10
