const { defineConfig } = require('@vue/cli-service')

module.exports = defineConfig({
  transpileDependencies: true,
  lintOnSave: false, // 禁用ESLint
  devServer: {
    port: 8082,
    proxy: {
      '/api': {
        target: 'http://localhost:8081',
        changeOrigin: true,
        pathRewrite: {
          '^/api': '/api'
        }
      }
    }
  },
  productionSourceMap: false,
  configureWebpack: {
    resolve: {
      fallback: {
        "path": require.resolve("path-browserify")
      }
    }
  },
  // 如果你想要在生产环境下使用 babel-plugin-dynamic-import-node
  // 需要将这个选项设置为 false
  // 它可以显著提高构建速度
  chainWebpack: config => {
    config.plugin('html')
      .tap(args => {
        args[0].title = 'MCP RAPI - API权限管理系统'
        return args
      })
  }
})