# .devcontainer/Dockerfile
FROM mcr.microsoft.com/devcontainers/go:1.25

# 设置环境变量
ENV GO111MODULE=on
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

# 安装系统依赖
USER root
RUN apt-get update && apt-get install -y \
    default-mysql-client \
    git \
    curl \
    wget \
    protobuf-compiler \
    && rm -rf /var/lib/apt/lists/*

# 安装 go-zero 的 goctl 工具（固定版本，避免 latest 漂移导致构建失败）
RUN go install github.com/zeromicro/go-zero/tools/goctl@v1.8.2

# 安装 protolint（Proto 文件检查工具）
RUN go install github.com/yoheimuta/protolint/cmd/protolint@latest

# 安装 protoc-gen-go（Proto Go 插件）
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 安装 protoc-gen-go-grpc（可选，如果用 gRPC）
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 安装 golangci-lint（Go 代码检查）
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.55.2

# 切换回非 root 用户
USER vscode

# 验证安装
RUN goctl --help >/dev/null && protolint version && golangci-lint version

WORKDIR /workspaces/SyncNote