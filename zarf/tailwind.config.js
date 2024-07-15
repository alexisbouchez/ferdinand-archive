/** @type {import('tailwindcss').Config} */
module.exports = {
  theme: {
    extend: {
      colors: {
        ferdinand: {
          50: "#fbf8f1",
          100: "#f5eedf",
          200: "#e7d4b5",
          300: "#dcbf95",
          400: "#cd9f6a",
          500: "#c3874c",
          600: "#b57341",
          700: "#975b37",
          800: "#7a4a32",
          900: "#633e2b",
          950: "#351f15",
        },
      },
      fontFamily: {
        sans: [
          "-apple-system",
          "BlinkMacSystemFont",
          "Segoe UI",
          "Roboto",
          "Helvetica",
          "Arial",
          "sans-serif",
          "Apple Color Emoji",
          "Segoe UI Emoji",
          "Segoe UI Symbol",
        ],
      },
      opacity: {
        1: "0.01",
        2.5: "0.025",
        5: "0.05",
        7.5: "0.075",
        15: "0.15",
      },
      width: {
        3.5: "0.875rem",
      },
    },
  },
  content: ["./views/**/*.templ"],
  plugins: [require("tailwindcss-animate"), require("@tailwindcss/forms")],
};
