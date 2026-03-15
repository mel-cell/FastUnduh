import forms from '@tailwindcss/forms';
import containerQueries from '@tailwindcss/container-queries';

/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./client/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'sans-serif'],
      },
      colors: {
        brand: {
          DEFAULT: '#0d93f2',
          dark: '#0b7ecf',
        },
        neutral: {
          950: '#0a0a0a',
        }
      }
    }
  },
  plugins: [forms, containerQueries],
}
