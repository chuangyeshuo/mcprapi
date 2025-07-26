module.exports = {
  title: 'MCP RAPI',

  /**
   * @type {boolean} true | false
   * @description 是否显示设置右面板
   */
  showSettings: true,

  /**
   * @type {boolean} true | false
   * @description 是否固定头部
   */
  fixedHeader: true,

  /**
   * @type {boolean} true | false
   * @description 是否在侧边栏中显示logo
   */
  sidebarLogo: true,

  /**
   * @type {boolean} true | false
   * @description 是否显示标签视图
   */
  tagsView: true,

  /**
   * @type {string | array} 'production' | ['production', 'development']
   * @description 需要显示错误日志组件的环境
   * 默认仅在生产环境显示
   */
  errorLog: 'production'
}