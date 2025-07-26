#!/bin/sh
###
 # @Author: lidi10@staff.weibo.com
 # @Date: 2025-07-26 20:02:30
 # @LastEditTime: 2025-07-26 22:11:02
 # @LastEditors: lidi10@staff.weibo.com
 # @Description: 
 # Copyright (c) 2023 by Weibo, All Rights Reserved. 
### 

# 数据库初始化脚本
# 检查数据库是否已经初始化，如果没有则进行初始化

echo "开始检查数据库初始化状态..."

# 等待数据库服务启动
echo "等待数据库服务启动..."
sleep 30

# 检查数据库连接
echo "检查数据库连接..."

# 尝试连接数据库并检查是否已有数据
# 这里我们检查 users 表是否存在且有数据
MYSQL_CMD="mysql -h${MYSQL_HOST:-mysql-dev} -P${MYSQL_PORT:-3306} -u${MYSQL_USER:-mcprapi} -p${MYSQL_PASSWORD:-devpassword} ${MYSQL_DATABASE:-api_auth_dev}"

# 检查数据库是否已初始化（检查是否存在 admin 用户）
ADMIN_EXISTS=$(echo "SELECT COUNT(*) FROM users WHERE username='admin';" | $MYSQL_CMD -s 2>/dev/null || echo "0")

if [ "$ADMIN_EXISTS" = "0" ] || [ -z "$ADMIN_EXISTS" ]; then
    echo "数据库未初始化，开始执行初始化..."
    
    # 运行初始化程序
    /app/init_admin --config /app/configs/dev.yaml
    
    if [ $? -eq 0 ]; then
        echo "✅ 数据库初始化完成！"
    else
        echo "❌ 数据库初始化失败！"
        exit 1
    fi
else
    echo "✅ 数据库已经初始化，跳过初始化步骤"
fi

echo "数据库初始化检查完成"