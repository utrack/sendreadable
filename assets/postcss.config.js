// postcss.config.js
module.exports = {
  plugins: [
    require('postcss-import'),
    require('tailwindcss'),
    require('autoprefixer'),
    // require('@fullhuman/postcss-purgecss')({
    //   content: ['src/**/*.html']
    // }),
  ]
}
