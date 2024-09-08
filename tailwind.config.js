/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["**/*.templ"],
  theme: {
    extend: {
        width: {
            "3ch": "3ch"
        },
    },
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
  daisyui: {
    themes: ["winter"],
  },
};
