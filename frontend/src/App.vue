<script setup lang="ts">
import { useHead } from '@unhead/vue'

useHead({ title: 'Tochka' })
</script>

<template>
  <UApp>
    <UMain>
      <RouterView />
    </UMain>
  </UApp>
</template>
