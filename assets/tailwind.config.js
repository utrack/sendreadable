require('@fullhuman/postcss-purgecss');

module.exports = {
  purge: {
    enabled: true,
    content: ['./src/**/*.html'],
  },
  theme: {
  },
  variants: {},
  plugins: [
    require('@tailwindcss/forms')
  ],
}

