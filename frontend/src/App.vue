<script setup>
import { ref, onMounted } from "vue";
// Wails automatically generates standard JS bindings here when you build/dev
import {
  GetPorts,
  ConnectSerial,
  SendCommand,
  SetAutoMode,
  HideWindow,
} from "../wailsjs/go/main/App";

const ports = ref([]);
const selectedPort = ref("");
const connectionStatus = ref("Disconnected");
const logPath = ref("");
const isAutoMode = ref(false);

const refreshPorts = async () => {
  ports.value = await GetPorts();
};

const connect = async () => {
  if (!selectedPort.value) return;
  connectionStatus.value = await ConnectSerial(selectedPort.value);
};

const sendManualCommand = async (color) => {
  if (isAutoMode.value) return; // Prevent manual override if auto is on
  await SendCommand(color);
};

const toggleAutoMode = async () => {
  const result = await SetAutoMode(isAutoMode.value, logPath.value);
  console.log(result);
};

const minimizeToTray = () => {
  HideWindow();
};

onMounted(() => {
  refreshPorts();
});
</script>

<template>
  <main class="p-6 font-sans text-gray-800">
    <h1 class="text-2xl font-bold mb-4">Teams Busy Light</h1>

    <!-- Port Selection -->
    <div class="mb-6 border p-4 rounded bg-gray-50">
      <h2 class="font-semibold mb-2">Hardware Setup</h2>
      <select v-model="selectedPort" class="border p-1 mr-2">
        <option disabled value="">Select COM Port</option>
        <option v-for="port in ports" :key="port" :value="port">
          {{ port }}
        </option>
      </select>
      <button @click="refreshPorts" class="bg-blue-200 px-3 py-1 mr-2 rounded">
        Refresh
      </button>
      <button @click="connect" class="bg-blue-500 text-white px-3 py-1 rounded">
        Connect
      </button>
      <span class="ml-3 text-sm">{{ connectionStatus }}</span>
    </div>

    <!-- Mode Toggle -->
    <div class="mb-6 border p-4 rounded">
      <h2 class="font-semibold mb-2">Operating Mode</h2>
      <label class="flex items-center space-x-2">
        <input type="checkbox" v-model="isAutoMode" @change="toggleAutoMode" />
        <span>Enable Automatic Mode (Teams Sync)</span>
      </label>

      <div v-if="isAutoMode" class="mt-4">
        <label class="block text-sm mb-1"
          >Teams Log Path (Leave blank for default):</label
        >
        <input
          v-model="logPath"
          type="text"
          placeholder="%localappdata%\Packages\MSTeams_8wekyb3d8bbwe\..."
          class="w-full border p-2 text-sm"
        />
        <button
          @click="toggleAutoMode"
          class="mt-2 bg-green-500 text-white px-3 py-1 rounded text-sm"
        >
          Update Path
        </button>
      </div>
    </div>

    <!-- Manual Controls -->
    <div
      class="mb-6 border p-4 rounded"
      :class="{ 'opacity-50 pointer-events-none': isAutoMode }"
    >
      <h2 class="font-semibold mb-2">Manual Override</h2>
      <div class="flex space-x-2">
        <button
          @click="sendManualCommand('Green')"
          class="bg-green-500 text-white px-4 py-2 rounded"
        >
          Available
        </button>
        <button
          @click="sendManualCommand('Red')"
          class="bg-red-500 text-white px-4 py-2 rounded"
        >
          Busy
        </button>
        <button
          @click="sendManualCommand('Yellow')"
          class="bg-yellow-500 text-white px-4 py-2 rounded"
        >
          Away
        </button>
        <button
          @click="sendManualCommand('BlinkRed')"
          class="bg-red-700 text-white px-4 py-2 rounded"
        >
          Calling
        </button>
      </div>
    </div>

    <!-- Background running -->
    <button
      @click="minimizeToTray"
      class="text-sm underline text-gray-500 hover:text-gray-700"
    >
      Minimize to System Tray
    </button>
  </main>
</template>
