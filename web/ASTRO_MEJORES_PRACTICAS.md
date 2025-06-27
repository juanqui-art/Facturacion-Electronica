# Mejores Prácticas de Astro 5.x 2025

## Arquitectura y Configuración

### 1. **TypeScript Strict Mode**
```typescript
// astro.config.mjs
export default defineConfig({
  typescript: {
    strict: true,
    checkJs: true
  }
});
```

### 2. **Islands Architecture**
- Hidratación parcial solo donde se necesita
- Componentes estáticos por defecto
- Uso de directivas `client:*` específicas

```astro
<!-- Solo hidratar cuando sea visible -->
<InteractiveComponent client:visible />

<!-- Hidratar inmediatamente -->
<UrgentComponent client:load />

<!-- Hidratar cuando sea necesario -->
<LazyComponent client:idle />
```

## Nuevas Características Astro 5.x

### 3. **Content Layer API**
```typescript
// astro.config.mjs
import { defineConfig } from 'astro/config';
import { contentLayer } from '@astrojs/content';

export default defineConfig({
  content: {
    collections: {
      blog: contentLayer({
        type: 'directory',
        directory: './src/content/blog',
        schema: z.object({
          title: z.string(),
          publishDate: z.date()
        })
      })
    }
  }
});
```

### 4. **Server Islands**
```astro
---
// Renderizado en el servidor con hidratación parcial
export const prerender = false;
---

<div>
  <StaticContent />
  <ServerIsland client:load>
    <DynamicServerContent />
  </ServerIsland>
</div>
```

## Manejo de Scripts y Estilos

### 5. **Scripts en Astro 5.x**
```astro
<!-- Scripts ya no se elevan automáticamente -->
<script is:inline>
  // Script ejecutado inline
  console.log('Ejecutado inmediatamente');
</script>

<script>
  // Script procesado y bundleado
  import { initializeApp } from './utils';
  initializeApp();
</script>
```

### 6. **Estilos Scoped Mejorados**
```astro
<style>
  /* Estilos scoped automáticamente */
  .card {
    @apply rounded-lg border p-4;
  }
  
  /* Variables CSS con tema automático */
  .theme-aware {
    background: var(--color-background);
    color: var(--color-text);
  }
</style>
```

## Gestión de Temas

### 7. **Dark Mode con Persistencia**
```astro
---
// ThemeManager.astro
---

<button id="theme-toggle" aria-label="Cambiar tema">
  <span class="sun-icon">☀️</span>
  <span class="moon-icon">🌙</span>
</button>

<script is:inline>
  // Ejecución inmediata para evitar flash
  const getTheme = () => {
    if (typeof localStorage !== 'undefined' && localStorage.getItem('theme')) {
      return localStorage.getItem('theme');
    }
    if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
      return 'dark';
    }
    return 'light';
  };
  
  const theme = getTheme();
  document.documentElement.classList.toggle('dark', theme === 'dark');
</script>
```

### 8. **Design Tokens CSS**
```css
/* design-tokens.css */
:root {
  /* Colores base */
  --color-primary-50: #eff6ff;
  --color-primary-500: #3b82f6;
  --color-primary-900: #1e3a8a;
  
  /* Espaciado (8pt grid) */
  --space-1: 0.25rem;  /* 4px */
  --space-2: 0.5rem;   /* 8px */
  --space-4: 1rem;     /* 16px */
  
  /* Tipografía */
  --font-sans: 'Inter', system-ui, sans-serif;
  --text-xs: 0.75rem;
  --text-sm: 0.875rem;
  
  /* Elevación */
  --shadow-sm: 0 1px 2px 0 rgb(0 0 0 / 0.05);
  --shadow-xl: 0 20px 25px -5px rgb(0 0 0 / 0.1);
}

.dark {
  --color-background: #0f172a;
  --color-text: #f1f5f9;
}
```

## Componentes Modernos

### 9. **Componentes con Props Tipados**
```astro
---
export interface Props {
  variant?: 'primary' | 'secondary' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  loading?: boolean;
  class?: string;
}

const { 
  variant = 'primary',
  size = 'md',
  disabled = false,
  loading = false,
  class: className = '',
  ...rest 
} = Astro.props;

const baseClasses = `
  inline-flex items-center justify-center font-medium
  rounded-xl transition-all duration-200 ease-out
  focus:outline-none focus:ring-2 focus:ring-offset-2
  disabled:opacity-50 disabled:cursor-not-allowed
  active:scale-95 transform hover:scale-105
`.replace(/\s+/g, ' ').trim();
---

<button 
  class={`${baseClasses} ${variantClasses[variant]} ${className}`}
  disabled={disabled || loading}
  {...rest}
>
  {loading && <LoadingSpinner />}
  <slot />
</button>
```

### 10. **Micro-interacciones**
```astro
<style>
  .card {
    @apply transition-all duration-300 ease-out;
    @apply hover:-translate-y-1 hover:shadow-xl;
  }
  
  .button {
    @apply active:scale-95 transform hover:scale-105;
    @apply transition-transform duration-150 ease-out;
  }
  
  /* Animaciones de entrada */
  @keyframes slideIn {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
  
  .animate-slide-in {
    animation: slideIn 0.3s ease-out;
  }
</style>
```

## Optimización de Performance

### 11. **Lazy Loading y Code Splitting**
```astro
---
// Importación dinámica para code splitting
const HeavyComponent = lazy(() => import('../components/HeavyComponent.astro'));
---

<!-- Carga diferida con intersección -->
<HeavyComponent client:visible />

<!-- Preload crítico -->
<link rel="preload" href="/critical.css" as="style" />
<link rel="preload" href="/hero-image.webp" as="image" />
```

### 12. **Optimización de Imágenes**
```astro
---
import { Image } from 'astro:assets';
import heroImage from '../assets/hero.jpg';
---

<Image 
  src={heroImage}
  alt="Descripción"
  width={800}
  height={600}
  format="webp"
  quality={80}
  loading="lazy"
  decoding="async"
/>
```

## Estructura de Proyecto

### 13. **Organización de Archivos**
```
src/
├── components/           # Componentes reutilizables
│   ├── ui/              # Componentes base (Button, Input, Card)
│   ├── layout/          # Componentes de layout
│   └── features/        # Componentes específicos de funcionalidad
├── layouts/             # Layouts de página
├── pages/               # Rutas de la aplicación
├── styles/              # Estilos globales
│   ├── design-tokens.css
│   └── global.css
├── utils/               # Utilidades y helpers
└── content/             # Contenido estructurado
```

### 14. **API Routes Optimizadas**
```typescript
// src/pages/api/data.json.ts
export async function GET({ params, request }) {
  const data = await fetchData();
  
  return new Response(JSON.stringify(data), {
    status: 200,
    headers: {
      'Content-Type': 'application/json',
      'Cache-Control': 'public, max-age=300', // 5 minutos
    }
  });
}
```

## Accesibilidad y SEO

### 15. **Accesibilidad Mejorada**
```astro
<!-- Navegación semántica -->
<nav aria-label="Navegación principal">
  <ul role="list">
    <li><a href="/" aria-current="page">Inicio</a></li>
    <li><a href="/about">Acerca de</a></li>
  </ul>
</nav>

<!-- Estados de loading accesibles -->
<div aria-live="polite" aria-busy={loading}>
  {loading ? 'Cargando...' : 'Contenido cargado'}
</div>

<!-- Focus management -->
<button 
  class="focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
  aria-describedby="button-help"
>
  Acción
</button>
```

### 16. **SEO Optimizado**
```astro
---
// Layout.astro
export interface Props {
  title: string;
  description?: string;
  canonical?: string;
  ogImage?: string;
}

const { title, description, canonical, ogImage } = Astro.props;
const fullTitle = `${title} | Mi Aplicación`;
---

<head>
  <title>{fullTitle}</title>
  <meta name="description" content={description} />
  <link rel="canonical" href={canonical || Astro.url.href} />
  
  <!-- Open Graph -->
  <meta property="og:title" content={fullTitle} />
  <meta property="og:description" content={description} />
  <meta property="og:image" content={ogImage} />
  <meta property="og:url" content={Astro.url.href} />
  
  <!-- JSON-LD -->
  <script type="application/ld+json" set:html={JSON.stringify(structuredData)} />
</head>
```

## Testing y Desarrollo

### 17. **Testing con Vitest**
```typescript
// vitest.config.ts
import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    environment: 'jsdom',
    setupFiles: './src/test/setup.ts'
  }
});
```

### 18. **Desarrollo con Hot Reload**
```javascript
// astro.config.mjs
export default defineConfig({
  server: {
    port: 4321,
    host: true, // Permite conexiones externas
  },
  devToolbar: {
    enabled: true
  }
});
```

## Integración con Frameworks

### 19. **React/Vue en Islas**
```astro
---
// Solo donde necesites interactividad
import ReactCounter from '../components/ReactCounter.jsx';
import VueForm from '../components/VueForm.vue';
---

<div>
  <h1>Contenido estático de Astro</h1>
  
  <!-- Solo esta parte se hidrata -->
  <ReactCounter client:visible initialCount={0} />
  <VueForm client:idle />
</div>
```

### 20. **Deployment y CI/CD**
```yaml
# .github/workflows/deploy.yml
name: Deploy to Production
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
      
      - run: npm ci
      - run: npm run build
      - run: npm run test
      
      # Deploy estático optimizado
      - uses: withastro/action@v2
        with:
          path: ./dist
```

## Mejores Prácticas Generales

### ✅ **Hacer:**
- Usar TypeScript strict mode
- Implementar lazy loading para componentes pesados
- Utilizar design tokens para consistencia
- Precargar recursos críticos
- Implementar dark mode desde el inicio
- Usar Server Islands para contenido dinámico
- Optimizar imágenes con el componente Image
- Mantener componentes pequeños y enfocados

### ❌ **Evitar:**
- Hidratación innecesaria (`client:load` en todo)
- Estilos inline sin design system
- JavaScript pesado en el cliente
- Falta de fallbacks para estados de carga
- Ignorar accesibilidad
- Bundle sizes grandes sin code splitting
- CSS no optimizado para dark mode

## Conclusión

Astro 5.x proporciona herramientas poderosas para crear aplicaciones web modernas, rápidas y accesibles. Las nuevas características como Content Layer y Server Islands permiten arquitecturas más flexibles, mientras que el enfoque en performance y developer experience hace que sea una excelente opción para proyectos de cualquier escala.

Estas prácticas han sido implementadas en nuestro sistema de facturación SRI, demostrando su efectividad en aplicaciones de producción reales.