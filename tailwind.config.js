/** @type {import('tailwindcss').Config} */
module.exports = {
    content: ["./internal/**/*.{html,go,js}"],
    safelist: [
      {
        //for tabs
        pattern: /grid-cols-.+/,
      }
    ],
    theme: {
      extend: {
        colors: {
          background: "hsl(var(--background))",
          foreground: "hsl(var(--foreground))",
          primary: {
            DEFAULT: "hsl(var(--primary))",
            foreground: "hsl(var(--primary-foreground))",
          },
          secondary: {
            DEFAULT: "hsl(var(--secondary))",
            foreground: "hsl(var(--secondary-foreground))",
          },
          muted: {
            DEFAULT: "hsl(var(--muted))",
            foreground: "hsl(var(--muted-foreground))",
          },
          danger: {
            DEFAULT: "hsl(var(--danger))",
            foreground: "hsl(var(--danger-foreground))",
          }
        }
      },
    },
    plugins: [],
  }
  
  