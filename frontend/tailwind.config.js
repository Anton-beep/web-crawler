/** @type {import('tailwindcss').Config} */

const colors = require('tailwindcss/colors')

module.exports = {
    darkMode: ["class"],
    content: ["./index.html", "./src/**/*.{ts,tsx,js,jsx}"],
    theme: {
        fontFamily: {
            'body': ["Nunito", "sans-serif"],
        },
        colors: {
            primary: colors.zinc[100],
            background: colors.zinc[900],
        }
    },
    plugins: [require("tailwindcss-animate")],
}
