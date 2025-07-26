const getters = {
  sidebar: state => state.app.sidebar,
  size: state => state.app.size,
  device: state => state.app.device,
  token: state => state.user.token,
  avatar: state => state.user.avatar,
  name: state => state.user.name,
  userId: state => state.user.userId,
  roles: state => state.user.roles,
  permissions: state => state.user.permissions,
  deptId: state => state.user.deptId,
  routes: state => state.permission.routes,
  addRoutes: state => state.permission.addRoutes
}

export default getters