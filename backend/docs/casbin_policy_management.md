# Casbin 策略管理最佳实践

## 问题背景

在使用 Casbin 进行权限管理时，`SavePolicy()` 方法存在一个潜在的风险：它会先清空数据库中的所有策略，然后将内存中的策略重新写入数据库。如果内存中的策略不完整或出现问题，就会导致数据库中的策略被意外清空。

## 解决方案

我们提供了以下几种安全的策略管理方法：

### 1. SafeSavePolicy() - 安全保存策略

```go
// 使用前会检查内存中的策略是否为空
err := enforcer.SafeSavePolicy()
```

**特点：**
- 在保存前检查策略是否为空
- 防止意外清空数据库
- 适用于需要全量保存策略的场景

### 2. SyncPolicyToMemory() - 同步策略到内存（推荐）

```go
// 从数据库重新加载策略到内存
err := enforcer.SyncPolicyToMemory()
```

**特点：**
- 从数据库加载最新策略到内存
- 不会修改数据库内容
- 确保内存和数据库的一致性
- **推荐在大多数场景中使用**

### 3. AddPolicyAndSync() - 添加策略并同步

```go
// 添加策略并自动同步到内存
added, err := enforcer.AddPolicyAndSync(sub, obj, act, eft)
```

**特点：**
- 添加策略后自动同步内存
- 确保操作的原子性
- 适用于单个策略添加场景

### 4. AddPolicyWithDeptAndSync() - 添加部门策略并同步

```go
// 添加带部门维度的策略并自动同步
added, err := enforcer.AddPolicyWithDeptAndSync(sub, obj, act, dept, eft)
```

**特点：**
- 支持部门维度的权限控制
- 添加后自动同步内存
- 适用于多租户场景

## 使用建议

### 1. 日常权限操作
- ✅ 使用 `SyncPolicyToMemory()` 进行策略同步
- ✅ 使用 `AddPolicyAndSync()` 系列方法添加策略
- ❌ 避免直接使用 `SavePolicy()`

### 2. 批量操作
```go
// 批量添加权限后统一同步
for _, perm := range permissions {
    enforcer.AddPolicyWithDept(role, perm.Path, perm.Method, dept, "allow")
}
// 最后统一同步
err := enforcer.SyncPolicyToMemory()
```

### 3. 紧急情况
如果必须使用 `SavePolicy()`，请先使用 `SafeSavePolicy()`：
```go
// 安全保存策略
err := enforcer.SafeSavePolicy()
if err != nil {
    log.Error("策略保存失败: %v", err)
    // 可以尝试重新加载策略
    enforcer.SyncPolicyToMemory()
}
```

## 当前系统中的应用

以下服务方法已经更新为使用安全的策略管理方式：

### RoleService
- `UpdateRolePermissions()` - 使用 `SyncPolicyToMemory()`
- `UpdateRolePermissionsByAPIIDs()` - 使用 `SyncPolicyToMemory()`
- `AssignRolePermissions()` - 使用 `SyncPolicyToMemory()`
- `addDefaultPermissionsToRoleWithoutSave()` - 使用 `AddPolicyWithDeptAndSync()`

### DeptPermissionService
- `GrantDeptAdmin()` - 使用 `SyncPolicyToMemory()`
- `RevokeDeptAdmin()` - 使用 `SyncPolicyToMemory()`
- `InitSystemAdmin()` - 使用 `SyncPolicyToMemory()`

## 监控和调试

### 1. 策略数量监控
```go
policies := enforcer.GetPolicy()
groupingPolicies := enforcer.GetGroupingPolicy()
log.Info("当前策略数量: %d, 分组策略数量: %d", len(policies), len(groupingPolicies))
```

### 2. 策略一致性检查
定期检查内存和数据库中的策略是否一致：
```go
// 重新加载策略并比较
oldPolicies := enforcer.GetPolicy()
enforcer.SyncPolicyToMemory()
newPolicies := enforcer.GetPolicy()
// 比较 oldPolicies 和 newPolicies
```

## 总结

通过使用这些安全的策略管理方法，我们可以：
1. 避免意外清空策略数据库
2. 确保内存和数据库的策略一致性
3. 提供更可靠的权限管理功能
4. 支持部门级别的权限控制

记住：**优先使用 `SyncPolicyToMemory()` 而不是 `SavePolicy()`**！