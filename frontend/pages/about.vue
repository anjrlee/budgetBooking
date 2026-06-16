<script setup lang="ts">
</script>

<template>
  <div class="flex flex-col items-center pt-8 text-center sm:pt-16">

    <div class="relative mt-8 max-w-sm w-full">

      <div class="parchment px-8 py-8">
        <h1 class="parchment-title font-display text-2xl font-bold mb-4">{{ $t('about.title') }}</h1>
        <p class="parchment-text leading-relaxed">{{ $t('about.description') }}</p>
      </div>

      <!-- fountain pen: upper-right → lower-left -->
      <svg class="pen" viewBox="0 0 30 200" xmlns="http://www.w3.org/2000/svg">
        <!-- clip -->
        <rect x="19" y="4" width="3.5" height="44" rx="1.75" fill="#1e3a8a"/>
        <ellipse cx="20.75" cy="47" rx="2" ry="2.5" fill="#1e3a8a"/>

        <!-- cap -->
        <path d="M8,2 Q8,0 15,0 Q22,0 22,2 L22,52 Q22,56 15,56 Q8,56 8,52Z" fill="#1d4ed8"/>
        <rect x="10" y="3" width="3" height="49" rx="1.5" fill="white" opacity="0.14"/>

        <!-- cap–barrel ring -->
        <rect x="7.5" y="54" width="15" height="5" rx="1.5" fill="#1e3a8a"/>

        <!-- barrel -->
        <rect x="8" y="59" width="14" height="95" rx="2" fill="#2563eb"/>
        <rect x="10" y="61" width="3" height="91" rx="1.5" fill="white" opacity="0.13"/>

        <!-- barrel–grip ring -->
        <rect x="7.5" y="152" width="15" height="5" rx="1.5" fill="#1e3a8a"/>

        <!-- grip section -->
        <path d="M9,157 L9,172 Q9,175 15,175 Q21,175 21,172 L21,157Z" fill="#1d4ed8"/>

        <!-- nib (gold) -->
        <path d="M10.5,173 L15,197 L19.5,173 Q17,175 15,175 Q13,175 10.5,173Z" fill="#ca9e10"/>
        <line x1="15" y1="174" x2="15" y2="195" stroke="#7a5a00" stroke-width="0.7" opacity="0.6"/>
        <path d="M11.5,175 L14,191 L14.5,175" fill="#f0c830" opacity="0.45"/>

        <!-- ink tip dot -->
        <ellipse cx="15" cy="197" rx="1.2" ry="1.5" fill="#1a1a6e" opacity="0.75"/>
      </svg>
    </div>

    <div class="mt-14">
      <NuxtLink to="/" class="text-xs text-ink-soft underline-offset-2 hover:underline">
        ← {{ $t('about.back') }}
      </NuxtLink>
    </div>
  </div>
</template>

<style scoped>
.parchment {
  position: relative;
  background:
    url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='200' height='200'%3E%3Cfilter id='n'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.75' numOctaves='4' stitchTiles='stitch'/%3E%3CfeColorMatrix type='saturate' values='0'/%3E%3C/filter%3E%3Crect width='200' height='200' filter='url(%23n)' opacity='0.06'/%3E%3C/svg%3E"),
    linear-gradient(160deg, #f7edd8 0%, #eedfc0 30%, #f2e4c8 60%, #e8d4a8 100%);
  border-radius: 2px;
  border: 1px solid #c9a96e;
  box-shadow:
    0 4px 24px rgba(100, 60, 10, 0.18),
    0 1px 4px rgba(100, 60, 10, 0.12),
    inset 0 0 40px rgba(180, 120, 40, 0.12),
    inset 0 2px 8px rgba(160, 100, 20, 0.10);
}

.parchment::before,
.parchment::after {
  content: '';
  position: absolute;
  left: 6px;
  right: 6px;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(160, 100, 20, 0.25) 20%, rgba(160, 100, 20, 0.25) 80%, transparent);
}
.parchment::before { top: 10px; }
.parchment::after  { bottom: 10px; }

.parchment-title {
  color: #3d2408;
  text-shadow: 0 1px 0 rgba(255, 240, 200, 0.6);
}

.parchment-text {
  color: #5c3d1a;
  font-size: 0.925rem;
  text-shadow: 0 1px 0 rgba(255, 240, 200, 0.6);
  font-family: Georgia, 'Times New Roman', serif;
}

/* fountain pen: drawn nib-down, rotated so cap→upper-right, nib→lower-left */
.pen {
  position: absolute;
  width: 22px;
  right: -42px;
  top: 20px;
  transform: rotate(45deg);
  transform-origin: 50% 50%;
  filter: drop-shadow(1px 3px 4px rgba(10, 20, 80, 0.25));
}
</style>
