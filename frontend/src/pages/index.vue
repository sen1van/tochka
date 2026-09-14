<script setup lang="ts">
import { computed, ref } from 'vue'
import { useHead } from '@unhead/vue'


type Sensor = { sensorId: number; name: string; lastSeen: string }
type Telemetry = { sensorId: number; tryNumber: number; value: number; timestamp: string }
type ApiList<T> = { total: number; data: T[] }
type RawRecord = Record<string, unknown>
type Screen = 'login' | 'connecting' | 'dashboard'


const token = ref('')
const screen = ref<Screen>('login')
const error = ref('')
const sensors = ref<Sensor[]>([])
const telemetry = ref<Telemetry[]>([])
const showTelemetryModal = ref(false)
const selectedSensor = ref<Sensor | null>(null)
const selectedTelemetry = ref<Telemetry[]>([])

useHead({
  title: 'Tochka',
  meta: [{ name: 'theme-color', content: '#ffffff' }]
})

const sensorColumns = [
  { accessorKey: 'name', header: 'Name' },
  { accessorKey: 'sensorId', header: 'ID' },
  { accessorKey: 'lastSeen', header: 'Last seen' }
]
const telemetryColumns = [
  { accessorKey: 'sensorId', header: 'Sensor' },
  { accessorKey: 'value', header: 'Value' },
  { accessorKey: 'tryNumber', header: 'Try' },
  { accessorKey: 'timestamp', header: 'Timestamp' }
]

const formattedSensors = computed(() => sensors.value.map(sensor => ({ ...sensor, lastSeen: formatDate(sensor.lastSeen) })))
const formattedTelemetry = computed(() => telemetry.value.map(reading => ({ ...reading, timestamp: formatDate(reading.timestamp) })))

function formatDate(value: string) {
  return value ? new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value)) : '—'
}

async function request<T>(path: string, options: RequestInit = {}) {
  const headers = new Headers(options.headers)
  headers.set('Accept', 'application/json')
  headers.set('Content-Type', 'application/json')
  if (token.value !== 'faketoken') headers.set('X-API-Token', token.value)
  const url = new URL(`/api${path}`, window.location.origin).toString()
  const startedAt = performance.now()
  console.debug('[api] request', { method: options.method ?? 'GET', url })

  try {
    const response = await fetch(url, { ...options, headers })
    const duration = Math.round(performance.now() - startedAt)
    const responseText = response.status === 204 ? '' : await response.clone().text()
    let responseBody: unknown = responseText
    try { responseBody = responseText ? JSON.parse(responseText) : undefined } catch { /* non-JSON response */ }
    console.debug('[api] response', { method: options.method ?? 'GET', url, status: response.status, durationMs: duration, body: responseBody })

    if (!response.ok) {
      const body = responseBody as { message?: string } | undefined
      throw new Error(body?.message ?? `Request failed (${response.status})`)
    }
    return response.status === 204 ? undefined as T : await response.json() as T
  } catch (cause) {
    const duration = Math.round(performance.now() - startedAt)
    console.error('[api] error', { method: options.method ?? 'GET', url, durationMs: duration, error: cause })
    throw cause
  }
}

async function connect() {
  if (!token.value.trim()) {
    error.value = 'Enter an API token to continue.'
    return
  }
  screen.value = 'connecting'
  error.value = ''
  try {
    if (token.value.trim() === 'faketoken') {
      sensors.value = [
        { sensorId: 101, name: 'Warehouse temperature', lastSeen: new Date().toISOString() },
        { sensorId: 102, name: 'Office humidity', lastSeen: new Date(Date.now() - 1000 * 60 * 18).toISOString() },
        { sensorId: 103, name: 'Loading dock', lastSeen: new Date(Date.now() - 1000 * 60 * 52).toISOString() }
      ]
      telemetry.value = [
        { sensorId: 101, tryNumber: 12, value: 23, timestamp: new Date().toISOString() },
        { sensorId: 102, tryNumber: 8, value: 48, timestamp: new Date(Date.now() - 1000 * 60 * 6).toISOString() },
        { sensorId: 103, tryNumber: 4, value: 19, timestamp: new Date(Date.now() - 1000 * 60 * 14).toISOString() }
      ]
    } else {
      await request<{ message: string }>('/ping')
      await loadDashboard()
    }
    localStorage.setItem('tochka-api-token', token.value)
    screen.value = 'dashboard'
  } catch (cause) {
    screen.value = 'login'
    error.value = cause instanceof Error ? cause.message : 'Connection failed'
  }
}

function warnSchema(path: string, message: string, details?: unknown) {
  console.warn(`[openapi] ${path}: ${message}`, details ?? '')
}

function checkObjectKeys(path: string, value: RawRecord, expected: string[]) {
  const actual = Object.keys(value)
  const missing = expected.filter(key => !(key in value))
  const extra = actual.filter(key => !expected.includes(key))
  const wrongCase = actual.filter(key => expected.some(expectedKey => key.toLowerCase() === expectedKey.toLowerCase() && key !== expectedKey))

  if (missing.length) warnSchema(path, `missing fields: ${missing.join(', ')}`, value)
  if (extra.length) warnSchema(path, `unexpected fields: ${extra.join(', ')}`, value)
  if (wrongCase.length) warnSchema(path, `field capitalization differs: ${wrongCase.join(', ')}`, value)
}

function checkListResponse(path: string, value: unknown, itemFields: string[]) {
  if (Array.isArray(value)) {
    warnSchema(path, 'expected an object with total and data, received an array', value)
    value.forEach((item, index) => {
      if (item && typeof item === 'object') checkObjectKeys(`${path}[${index}]`, item as RawRecord, itemFields)
      else warnSchema(`${path}[${index}]`, 'expected an object', item)
    })
    return
  }

  if (!value || typeof value !== 'object') {
    warnSchema(path, 'expected an object response', value)
    return
  }

  const response = value as RawRecord
  checkObjectKeys(path, response, ['total', 'data'])
  if (!Array.isArray(response.data)) {
    warnSchema(`${path}.data`, 'expected an array', response.data)
    return
  }
  response.data.forEach((item, index) => {
    if (item && typeof item === 'object') checkObjectKeys(`${path}.data[${index}]`, item as RawRecord, itemFields)
    else warnSchema(`${path}.data[${index}]`, 'expected an object', item)
  })
}

function unwrapList(value: unknown): RawRecord[] {
  if (Array.isArray(value)) return value as RawRecord[]
  if (value && typeof value === 'object' && Array.isArray((value as ApiList<RawRecord>).data)) return (value as ApiList<RawRecord>).data
  return []
}

function normalizeSensor(raw: RawRecord): Sensor {
  const lastTelemetry = (raw.lastTelemetry ?? raw.LastTelemetry) as RawRecord | null | undefined
  return {
    sensorId: Number(raw.sensorId ?? raw.ID ?? raw.id),
    name: String(raw.name ?? raw.Name ?? ''),
    lastSeen: String(raw.lastSeen ?? raw.LastSeen ?? lastTelemetry?.timestamp ?? lastTelemetry?.Timestamp ?? '')
  }
}

function normalizeTelemetry(raw: RawRecord): Telemetry {
  return {
    sensorId: Number(raw.sensorId ?? raw.SensorID ?? raw.sensor_id ?? raw.ID ?? raw.id),
    tryNumber: Number(raw.tryNumber ?? raw.TryNumber ?? raw.try_number ?? raw.Try ?? 0),
    value: Number(raw.value ?? raw.Value ?? 0),
    timestamp: String(raw.timestamp ?? raw.Timestamp ?? '')
  }
}

async function loadDashboard() {
  const [sensorResponse, telemetryResponse] = await Promise.all([
    request<unknown>('/sensors'),
    request<unknown>('/telemetry?limit=20&offset=0')
  ])
  checkListResponse('/sensors', sensorResponse, ['sensorId', 'name', 'lastTelemetry'])
  checkListResponse('/telemetry', telemetryResponse, ['sensorId', 'tryNumber', 'value', 'timestamp'])
  sensors.value = unwrapList(sensorResponse).map(normalizeSensor)
  telemetry.value = unwrapList(telemetryResponse).map(normalizeTelemetry)
}

async function refresh() {
  error.value = ''
  try { await loadDashboard() } catch (cause) { error.value = cause instanceof Error ? cause.message : 'Unable to load dashboard' }
}

async function showSensorTelemetry(sensor: Sensor) {
  if (!sensor?.sensorId) {
    console.warn('[ui] unable to identify selected sensor', sensor)
    error.value = 'Unable to identify the selected sensor'
    return
  }

  selectedSensor.value = sensor
  error.value = ''
  try {
    if (token.value === 'faketoken') {
      selectedTelemetry.value = telemetry.value.filter(reading => reading.sensorId === sensor.sensorId)
    } else {
      const response = await request<unknown>(`/sensors/${sensor.sensorId}/telemetry?limit=20&offset=0`)
      checkListResponse(`/sensors/${sensor.sensorId}/telemetry`, response, ['sensorId', 'tryNumber', 'value', 'timestamp'])
      selectedTelemetry.value = unwrapList(response).map(normalizeTelemetry)
    }
    showTelemetryModal.value = true
  } catch (cause) {
    error.value = cause instanceof Error ? cause.message : 'Unable to load sensor telemetry'
  }
}

function signOut() {
  token.value = ''
  screen.value = 'login'
  sensors.value = []
  telemetry.value = []
  selectedSensor.value = null
  selectedTelemetry.value = []
  localStorage.removeItem('tochka-api-token')
}
</script>

<template>
  <UApp>
    <UContainer class="py-8">
      <section v-if="screen !== 'dashboard'" class="mx-auto max-w-md py-16">
        <div class="mb-4 flex justify-end"><UButton to="/openapi" label="API reference" icon="i-lucide-book-open" color="neutral" variant="outline" /></div>
        <UPageHeader title="Вход в панель" description="Введите токен." class="mb-8" />
        <UAlert v-if="error" :description="error" color="error" variant="subtle" class="mb-4" />
        <UCard>
          <form class="space-y-4" @submit.prevent="connect">
            <UFormField label="API token" hint="Sent as X-API-Token"><UInput v-model="token" class="w-full" type="password" placeholder="Enter token" autofocus /></UFormField>
            <UButton type="submit" block color="neutral" :loading="screen === 'connecting'" :label="screen === 'connecting' ? 'Checking connection…' : 'Вход'" />
          </form>
        </UCard>
      </section>

      <template v-else>
        <div class="mb-8 flex items-end justify-between gap-4"><UPageHeader title="Панель" /><div class="flex gap-2"><UButton to="/openapi" label="API reference" icon="i-lucide-book-open" color="neutral" variant="outline" /><UButton icon="i-lucide-refresh-cw" color="neutral" variant="ghost" aria-label="Refresh" @click="refresh" /><UButton label="Выход" color="neutral" variant="outline" @click="signOut" /></div></div>
        <UAlert v-if="error" :description="error" color="error" variant="subtle" class="mb-6" />
        <div class="grid gap-4 md:grid-cols-3"><UCard><p class="text-sm text-muted">Датчики</p><p class="text-3xl font-semibold">{{ sensors.length }}</p></UCard><UCard><p class="text-sm text-muted">Показаний</p><p class="text-3xl font-semibold">{{ telemetry.length }}</p></UCard><UCard><p class="text-sm text-muted">Последний показатель</p><p class="text-3xl font-semibold">{{ telemetry[0]?.value ?? '—' }}</p></UCard></div>
        <div class="mt-6 grid gap-6 lg:grid-cols-2"><UCard><template #header><div><h2 class="font-semibold">Датчики</h2></div></template><div v-if="sensors.length" class="divide-y divide-default"><button v-for="sensor in sensors" :key="sensor.sensorId" type="button" class="flex w-full items-center justify-between gap-4 px-4 py-4 text-left hover:bg-elevated" @click="showSensorTelemetry(sensor)"><span><span class="block font-medium">{{ sensor.name }}</span><span class="block text-sm text-muted">ID {{ sensor.sensorId }} · {{ formatDate(sensor.lastSeen) }}</span></span><UIcon name="i-lucide-chevron-right" class="text-muted" /></button></div><p v-else class="py-8 text-center text-sm text-muted">Датчики не найдены</p></UCard><UCard><template #header><h2 class="font-semibold">Показания</h2></template><UTable :data="formattedTelemetry" :columns="telemetryColumns" /><p v-if="!telemetry.length" class="py-8 text-center text-sm text-muted">No telemetry found.</p></UCard></div>
      </template>

      <UModal v-model:open="showTelemetryModal" :title="selectedSensor ? `${selectedSensor.name} telemetry` : 'Sensor telemetry'"><template #body><UTable :data="selectedTelemetry.map(reading => ({ ...reading, timestamp: formatDate(reading.timestamp) }))" :columns="telemetryColumns" /><p v-if="!selectedTelemetry.length" class="py-8 text-center text-sm text-muted">Телеметрия не найдена</p></template></UModal>
    </UContainer>
  </UApp>
</template>
