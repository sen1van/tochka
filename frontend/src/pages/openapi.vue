<script setup lang="ts">
import { ApiReference } from '@scalar/api-reference'
import '@scalar/api-reference/style.css'

const currentOrigin = window.location.origin
const openApiUrl = new URL('/openapi.yaml', currentOrigin).toString()
</script>

<template>
  <ApiReference
    :configuration="{
      url: openApiUrl,
      hideClientButton: true,
      servers: [{ url: currentOrigin }]
    }"
  />
</template>
