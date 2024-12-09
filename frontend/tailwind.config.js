/** @type {import('tailwindcss').Config} */

import colors from 'tailwindcss/colors';
import tailwindcssAnimate from 'tailwindcss-animate';

module.exports = {
    darkMode: ["class"],
    content: ["./index.html", "./src/**/*.{ts,tsx,js,jsx}"],
    theme: {
        extend: {
            fontFamily: {
                'body': ["IBM Plex Sans", "sans-serif"],
            },
            colors: {
                primary: colors.zinc[100],
                background: colors.zinc[900],
                accent: colors.cyan[500],
                secondary: colors.zinc[400],
                error: colors.red[500],
                warning: colors.yellow[500],
                success: colors.green[500],
            },
            fontSize: {
                sm: '1rem',
                base: '1.2rem',
                xl: '1.45rem',
                '2xl': '1.763rem',
                '3xl': '2.153rem',
                '4xl': '2.641rem',
                '5xl': '3.252rem',
            },
            animation: {
                shimmer: "shimmer 2s linear infinite",
            },
            keyframes: {
                shimmer: {
                    from: {
                        backgroundPosition: "0 0",
                    },
                    to: {
                        backgroundPosition: "-200% 0",
                    },
                },
            },
        },
    },
    plugins: [tailwindcssAnimate],
}
