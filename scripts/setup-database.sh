#!/bin/bash

# 数据库选择脚本
# 用于配置项目使用的数据库类型

set -e

echo "======================================"
echo "UYou Go API Starter - 数据库配置"
echo "======================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 显示菜单
show_menu() {
    echo "请选择要使用的数据库："
    echo "  1) PostgreSQL (默认)"
    echo "  2) MongoDB"
    echo "  3) PostgreSQL + MongoDB (同时使用)"
    echo "  4) PostgreSQL + Redis"
    echo "  5) MongoDB + Redis"
    echo "  6) 全部 (PostgreSQL + MongoDB + Redis)"
    echo ""
}

# 更新配置文件
update_config() {
    local use_postgres=$1
    local use_mongodb=$2
    local use_redis=$3
    
    echo ""
    echo -e "${GREEN}正在更新配置文件...${NC}"
    
    # 更新 config.yaml
    if [ "$use_mongodb" = "true" ]; then
        sed -i 's/mongodb:/mongodb:\n  enabled: true/' configs/config.yaml 2>/dev/null || true
    fi
    
    if [ "$use_redis" = "true" ]; then
        sed -i 's/redis:/redis:\n  enabled: true/' configs/config.yaml 2>/dev/null || true
    fi
    
    echo -e "${GREEN}配置文件已更新${NC}"
}

# 更新 docker-compose.yml
update_docker_compose() {
    local use_postgres=$1
    local use_mongodb=$2
    local use_redis=$3
    
    echo -e "${GREEN}正在更新 docker-compose.yml...${NC}"
    
    # 备份原文件
    cp docker-compose.yml docker-compose.yml.bak
    
    # 创建新的 docker-compose.yml
    cat > docker-compose.yml << 'EOF'
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
      - "9090:9090"  # gRPC
      - "9091:9091"  # Prometheus metrics
    environment:
      - APP_ENVIRONMENT=development
      - DATABASE_HOST=${DATABASE_HOST:-db}
      - DATABASE_PORT=${DATABASE_PORT:-5432}
      - DATABASE_USER=${DATABASE_USER:-postgres}
      - DATABASE_PASSWORD=${DATABASE_PASSWORD:-postgres}
      - DATABASE_NAME=${DATABASE_NAME:-uyou_api}
      - JWT_SECRET=${JWT_SECRET:-your-secret-key-change-this-in-production}
EOF

    if [ "$use_mongodb" = "true" ]; then
        cat >> docker-compose.yml << 'EOF'
      - MONGODB_ENABLED=true
      - MONGODB_URI=mongodb://mongodb:27017
      - MONGODB_DATABASE=uyou_api
EOF
    fi
    
    if [ "$use_redis" = "true" ]; then
        cat >> docker-compose.yml << 'EOF'
      - REDIS_ENABLED=true
      - REDIS_HOST=redis
      - REDIS_PORT=6379
EOF
    fi
    
    cat >> docker-compose.yml << 'EOF'
    depends_on:
EOF
    
    if [ "$use_postgres" = "true" ]; then
        cat >> docker-compose.yml << 'EOF'
      - db
EOF
    fi
    
    if [ "$use_mongodb" = "true" ]; then
        cat >> docker-compose.yml << 'EOF'
      - mongodb
EOF
    fi
    
    if [ "$use_redis" = "true" ]; then
        cat >> docker-compose.yml << 'EOF'
      - redis
EOF
    fi
    
    if [ "$use_postgres" = "true" ]; then
        cat >> docker-compose.yml << 'EOF'

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=${DATABASE_USER:-postgres}
      - POSTGRES_PASSWORD=${DATABASE_PASSWORD:-postgres}
      - POSTGRES_DB=${DATABASE_NAME:-uyou_api}
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
EOF
    fi
    
    if [ "$use_mongodb" = "true" ]; then
        cat >> docker-compose.yml << 'EOF'

  mongodb:
    image: mongo:7
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
EOF
    fi
    
    if [ "$use_redis" = "true" ]; then
        cat >> docker-compose.yml << 'EOF'

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
EOF
    fi
    
    cat >> docker-compose.yml << 'EOF'

volumes:
EOF
    
    if [ "$use_postgres" = "true" ]; then
        echo "  postgres_data:" >> docker-compose.yml
    fi
    
    if [ "$use_mongodb" = "true" ]; then
        echo "  mongodb_data:" >> docker-compose.yml
    fi
    
    if [ "$use_redis" = "true" ]; then
        echo "  redis_data:" >> docker-compose.yml
    fi
    
    echo -e "${GREEN}docker-compose.yml 已更新${NC}"
}

# 显示配置摘要
show_summary() {
    local use_postgres=$1
    local use_mongodb=$2
    local use_redis=$3
    
    echo ""
    echo "======================================"
    echo "配置摘要"
    echo "======================================"
    echo -e "PostgreSQL: ${use_postgres}"
    echo -e "MongoDB:    ${use_mongodb}"
    echo -e "Redis:      ${use_redis}"
    echo "======================================"
    echo ""
}

# 主逻辑
main() {
    show_menu
    
    read -p "请输入选项 [1-6]: " choice
    
    case $choice in
        1)
            echo -e "${GREEN}选择: PostgreSQL${NC}"
            update_config "true" "false" "false"
            update_docker_compose "true" "false" "false"
            show_summary "true" "false" "false"
            ;;
        2)
            echo -e "${GREEN}选择: MongoDB${NC}"
            update_config "false" "true" "false"
            update_docker_compose "false" "true" "false"
            show_summary "false" "true" "false"
            ;;
        3)
            echo -e "${GREEN}选择: PostgreSQL + MongoDB${NC}"
            update_config "true" "true" "false"
            update_docker_compose "true" "true" "false"
            show_summary "true" "true" "false"
            ;;
        4)
            echo -e "${GREEN}选择: PostgreSQL + Redis${NC}"
            update_config "true" "false" "true"
            update_docker_compose "true" "false" "true"
            show_summary "true" "false" "true"
            ;;
        5)
            echo -e "${GREEN}选择: MongoDB + Redis${NC}"
            update_config "false" "true" "true"
            update_docker_compose "false" "true" "true"
            show_summary "false" "true" "true"
            ;;
        6)
            echo -e "${GREEN}选择: 全部数据库${NC}"
            update_config "true" "true" "true"
            update_docker_compose "true" "true" "true"
            show_summary "true" "true" "true"
            ;;
        *)
            echo -e "${RED}无效选项，使用默认配置 (PostgreSQL)${NC}"
            update_config "true" "false" "false"
            update_docker_compose "true" "false" "false"
            show_summary "true" "false" "false"
            ;;
    esac
    
    echo ""
    echo -e "${GREEN}配置完成！${NC}"
    echo ""
    echo "下一步："
    echo "  1. 运行 'make quick-start' 启动服务"
    echo "  2. 或运行 'docker-compose up -d' 手动启动"
    echo ""
}

# 运行主函数
main
