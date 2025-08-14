# UI Setup Documentation

## Tailwind CSS 4 + shadcn/ui Configuration

This document outlines the UI setup for the CS2 Log Generator project.

### Components Available

#### Button Component
```tsx
import { Button } from '@/components/ui/button'

// Standard variants
<Button variant="default">Default</Button>
<Button variant="destructive">Destructive</Button>
<Button variant="outline">Outline</Button>
<Button variant="secondary">Secondary</Button>
<Button variant="ghost">Ghost</Button>
<Button variant="link">Link</Button>

// CS2 specific variants
<Button variant="ct">Counter-Terrorist</Button>
<Button variant="t">Terrorist</Button>
```

#### Card Component
```tsx
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

<Card>
  <CardHeader>
    <CardTitle>Title</CardTitle>
    <CardDescription>Description</CardDescription>
  </CardHeader>
  <CardContent>
    Content goes here
  </CardContent>
</Card>
```

### Theme Configuration

#### CSS Variables
The theme uses CSS variables for dynamic color switching:
- Light/dark mode support
- CS2-specific colors: CT blue (#5E98D9), T orange (#DE9B35)

#### Custom Utility Classes
```css
.text-ct    /* Counter-Terrorist text color */
.text-t     /* Terrorist text color */
.bg-ct      /* Counter-Terrorist background */
.bg-t       /* Terrorist background */
.border-ct  /* Counter-Terrorist border */
.border-t   /* Terrorist border */
```

### Development Setup

#### Adding New Components
1. Use shadcn/ui MCP server: `mcp__shadcn-ui__get_component`
2. Save to `/src/components/ui/[component].tsx`
3. Export from `/src/components/ui/index.ts`

#### Theme Toggle
```tsx
const toggleTheme = () => {
  document.documentElement.classList.toggle('dark')
}
```

### File Structure
```
/src
  /components
    /ui
      button.tsx      # Button component with CS2 variants
      card.tsx        # Card component suite
      index.ts        # UI component exports
  /lib
    utils.ts          # cn() utility function
  index.css           # Tailwind directives + CSS variables
```

### Next Steps for UI Development
1. Add more shadcn/ui components as needed (form, input, select, etc.)
2. Create custom CS2-specific components
3. Implement responsive layout components
4. Add animations and transitions

### Testing the Setup
Run `npm run dev` and visit http://localhost:5173 to see the theme demo with:
- Theme switching (light/dark)
- CS2 color variants
- Component examples
- Status indicators