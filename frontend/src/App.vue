<script setup>
import { ref, onMounted } from "vue";
import {
  GetPorts,
  ConnectSerial,
  DisconnectSerial,
  SendCommand,
  SetAutoMode,
  HideWindow,
} from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";

const ports = ref([]);
const selectedPort = ref("");
const connectionStatus = ref("Disconnected");
const logPath = ref("");
const isAutoMode = ref(false);
const activeTab = ref("hardware"); // Default to hardware to see the monitor easily

// Serial Monitor State
const showSerialMonitor = ref(false);
const serialLogs = ref([]);

const refreshPorts = async () => {
  ports.value = await GetPorts();
};

const connect = async () => {
  if (!selectedPort.value) return;
  connectionStatus.value = await ConnectSerial(selectedPort.value);
};

const disconnect = async () => {
  connectionStatus.value = await DisconnectSerial();
};

const sendManualCommand = async (color) => {
  if (isAutoMode.value) return;
  await SendCommand(color);
};

const toggleAutoMode = async () => {
  const result = await SetAutoMode(isAutoMode.value, logPath.value);
  console.log(result);
};

const minimizeToTray = () => {
  HideWindow();
};

const clearLogs = () => {
  serialLogs.value = [];
};

onMounted(() => {
  refreshPorts();

  // Listen for incoming serial data from the Go backend
  EventsOn("serial-data", (data) => {
    serialLogs.value.push(data);

    // Auto-scroll logic (keep maximum 100 lines to prevent memory issues)
    if (serialLogs.value.length > 100) {
      serialLogs.value.shift();
    }
  });
});
</script>

<template>
  <div
    class="flex h-screen w-screen bg-gray-50 font-sans text-gray-800 overflow-hidden"
  >
    <!-- LEFT SIDEBAR -->
    <aside
      class="w-56 bg-[#2d3238] flex flex-col justify-between z-10 shadow-lg"
    >
      <div>
        <div class="h-16 flex items-center px-6 bg-[#25292e]">
          <div class="flex space-x-1 items-center">
            <div class="w-3 h-3 rounded-full bg-green-500"></div>
            <div class="w-3 h-3 rounded-full bg-red-500"></div>
            <h1 class="text-white font-semibold text-lg ml-2 tracking-wide">
              BUSY LIGHT
            </h1>
          </div>
        </div>

        <nav class="mt-4 flex flex-col">
          <button
            @click="activeTab = 'manual'"
            :class="[
              'text-left px-6 py-3 transition-colors border-l-4',
              activeTab === 'manual'
                ? 'bg-[#3b4148] text-white border-green-500'
                : 'text-gray-400 border-transparent hover:bg-[#363b41] hover:text-gray-200',
            ]"
          >
            Solid Color
          </button>

          <button
            @click="activeTab = 'auto'"
            :class="[
              'text-left px-6 py-3 transition-colors border-l-4',
              activeTab === 'auto'
                ? 'bg-[#3b4148] text-white border-green-500'
                : 'text-gray-400 border-transparent hover:bg-[#363b41] hover:text-gray-200',
            ]"
          >
            Teams Sync
          </button>

          <button
            @click="activeTab = 'hardware'"
            :class="[
              'text-left px-6 py-3 transition-colors border-l-4',
              activeTab === 'hardware'
                ? 'bg-[#3b4148] text-white border-green-500'
                : 'text-gray-400 border-transparent hover:bg-[#363b41] hover:text-gray-200',
            ]"
          >
            Settings
          </button>
        </nav>
      </div>

      <button
        @click="minimizeToTray"
        class="px-6 py-4 text-sm text-gray-400 hover:text-white hover:bg-[#363b41] text-left transition-colors flex items-center"
      >
        <svg
          class="w-4 h-4 mr-2"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 9l-7 7-7-7"
          ></path>
        </svg>
        Minimize to Tray
      </button>
    </aside>

    <!-- MAIN CONTENT AREA -->
    <main class="flex-1 flex flex-col relative bg-white">
      <header
        class="h-16 border-b flex items-center justify-between px-8 bg-white"
      >
        <h2 class="text-xl font-light text-gray-600 capitalize">
          {{
            activeTab === "manual"
              ? "Select Color"
              : activeTab === "auto"
                ? "Teams Synchronization"
                : "Device Settings"
          }}
        </h2>
      </header>

      <div class="p-8 flex-1 overflow-y-auto pb-16">
        <!-- 1. MANUAL CONTROL -->
        <div
          v-show="activeTab === 'manual'"
          class="flex flex-col items-center justify-center h-full pt-10"
          :class="{ 'opacity-40 pointer-events-none': isAutoMode }"
        >
          <div
            v-if="isAutoMode"
            class="absolute top-20 bg-yellow-100 text-yellow-800 px-4 py-2 rounded shadow-sm text-sm"
          >
            Manual control is disabled while Teams Sync is active.
          </div>

          <div class="flex flex-wrap justify-center gap-8 max-w-2xl">
            <button
              @click="sendManualCommand('Off')"
              class="group flex flex-col items-center focus:outline-none"
            >
              <div
                class="w-20 h-20 rounded-full bg-gray-200 shadow-inner flex items-center justify-center transition-all group-hover:scale-105 group-hover:bg-gray-300"
              >
                <div
                  class="w-12 h-12 rounded-full border-4 border-white opacity-50"
                ></div>
              </div>
              <span class="mt-3 text-sm font-medium text-gray-500">Off</span>
            </button>
            <button
              @click="sendManualCommand('Red')"
              class="group flex flex-col items-center focus:outline-none"
            >
              <div
                class="w-20 h-20 rounded-full bg-[#f44336] shadow-[0_0_15px_rgba(244,67,54,0.4)] transition-all group-hover:scale-105"
              ></div>
              <span class="mt-3 text-sm font-medium text-gray-500">Busy</span>
            </button>
            <button
              @click="sendManualCommand('Green')"
              class="group flex flex-col items-center focus:outline-none"
            >
              <div
                class="w-20 h-20 rounded-full bg-[#8bc34a] shadow-[0_0_15px_rgba(139,195,74,0.4)] transition-all group-hover:scale-105"
              ></div>
              <span class="mt-3 text-sm font-medium text-gray-500"
                >Available</span
              >
            </button>
            <button
              @click="sendManualCommand('Yellow')"
              class="group flex flex-col items-center focus:outline-none"
            >
              <div
                class="w-20 h-20 rounded-full bg-[#ffeb3b] shadow-[0_0_15px_rgba(255,235,59,0.4)] transition-all group-hover:scale-105"
              ></div>
              <span class="mt-3 text-sm font-medium text-gray-500">Away</span>
            </button>
            <button
              @click="sendManualCommand('BlinkRed')"
              class="group flex flex-col items-center focus:outline-none"
            >
              <div
                class="w-20 h-20 rounded-full bg-[#d32f2f] shadow-[0_0_20px_rgba(211,47,47,0.6)] animate-pulse transition-all group-hover:scale-105"
              ></div>
              <span class="mt-3 text-sm font-medium text-gray-500"
                >Calling</span
              >
            </button>
          </div>
        </div>

        <!-- 2. AUTO MODE (Teams Sync) -->
        <div v-show="activeTab === 'auto'" class="max-w-xl mx-auto mt-6">
          <div class="bg-gray-50 p-6 rounded-lg border">
            <label class="flex items-center space-x-3 cursor-pointer mb-6">
              <input
                type="checkbox"
                v-model="isAutoMode"
                @change="toggleAutoMode"
                class="w-5 h-5 text-green-500 rounded border-gray-300 focus:ring-green-400"
              />
              <span class="text-gray-700 font-medium text-lg"
                >Enable MS Teams Sync</span
              >
            </label>

            <div
              class="transition-opacity"
              :class="
                isAutoMode ? 'opacity-100' : 'opacity-40 pointer-events-none'
              "
            >
              <p class="text-sm text-gray-500 mb-2">
                Override Default Teams Log Path (Optional):
              </p>
              <input
                v-model="logPath"
                type="text"
                placeholder="%localappdata%\Packages\MSTeams_8wekyb3d8bbwe\LocalCache\Microsoft\MSTeams\Logs"
                class="w-full border-gray-300 rounded shadow-sm p-3 text-sm focus:border-green-500 focus:ring focus:ring-green-200"
              />
              <button
                @click="toggleAutoMode"
                class="mt-4 bg-[#8bc34a] hover:bg-[#7cb342] text-white px-6 py-2 rounded shadow transition-colors text-sm font-medium"
              >
                Apply Path Update
              </button>
            </div>
          </div>
        </div>

        <!-- 3. HARDWARE SETTINGS -->
        <div v-show="activeTab === 'hardware'" class="max-w-2xl mx-auto mt-6">
          <div class="bg-gray-50 p-6 rounded-lg border mb-6">
            <h3 class="text-gray-700 font-medium mb-4">Serial Connection</h3>

            <div class="flex space-x-3 mb-2">
              <select
                v-model="selectedPort"
                :disabled="connectionStatus.includes('Connected')"
                class="flex-1 border-gray-300 rounded shadow-sm p-2 text-sm focus:border-green-500 focus:ring focus:ring-green-200 disabled:bg-gray-200 disabled:text-gray-400"
              >
                <option disabled value="">Select COM Port</option>
                <option v-for="port in ports" :key="port" :value="port">
                  {{ port }}
                </option>
              </select>

              <button
                @click="refreshPorts"
                :disabled="connectionStatus.includes('Connected')"
                class="bg-gray-200 hover:bg-gray-300 text-gray-700 px-4 py-2 rounded shadow-sm transition-colors text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Refresh
              </button>

              <!-- Show Connect if disconnected -->
              <button
                v-if="!connectionStatus.includes('Connected')"
                @click="connect"
                class="bg-[#2d3238] hover:bg-[#3b4148] text-white px-6 py-2 rounded shadow transition-colors text-sm font-medium"
              >
                Connect
              </button>

              <!-- Show Disconnect if connected -->
              <button
                v-else
                @click="disconnect"
                class="bg-red-500 hover:bg-red-600 text-white px-6 py-2 rounded shadow transition-colors text-sm font-medium"
              >
                Disconnect
              </button>
            </div>
          </div>

          <!-- DEBUGGING / SERIAL MONITOR -->
          <div class="bg-gray-50 p-6 rounded-lg border">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-gray-700 font-medium">Serial Monitor</h3>
              <button
                @click="showSerialMonitor = !showSerialMonitor"
                class="text-sm text-green-600 hover:text-green-700 font-medium transition-colors"
              >
                {{ showSerialMonitor ? "Hide Monitor" : "Show Monitor" }}
              </button>
            </div>

            <div v-if="showSerialMonitor">
              <div
                class="bg-[#1e1e1e] text-[#d4d4d4] font-mono text-xs p-4 rounded-lg h-64 overflow-y-auto shadow-inner break-all flex flex-col-reverse"
              >
                <!-- flex-col-reverse makes new messages appear at the bottom naturally if we reverse the array -->
                <div class="flex flex-col">
                  <div
                    v-if="serialLogs.length === 0"
                    class="text-gray-500 italic py-2"
                  >
                    Waiting for incoming serial data...
                  </div>
                  <div
                    v-for="(log, index) in serialLogs"
                    :key="index"
                    class="py-0.5 border-b border-[#333] last:border-0"
                  >
                    <span class="text-gray-500 mr-2"
                      >[{{ new Date().toLocaleTimeString() }}]</span
                    >
                    {{ log }}
                  </div>
                </div>
              </div>
              <div class="mt-3 flex justify-end">
                <button
                  @click="clearLogs"
                  class="text-xs text-gray-500 hover:text-red-500 transition-colors border border-gray-300 px-3 py-1 rounded"
                >
                  Clear Output
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- BOTTOM STATUS BAR -->
      <footer
        class="absolute bottom-0 w-full h-10 bg-[#363b41] flex items-center px-6 border-t border-[#464c54]"
      >
        <div class="flex items-center space-x-2">
          <svg
            v-if="connectionStatus.includes('Connected')"
            class="w-4 h-4 text-green-400"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
              clip-rule="evenodd"
            ></path>
          </svg>
          <svg
            v-else
            class="w-4 h-4 text-gray-400"
            fill="currentColor"
            viewBox="0 0 20 20"
          >
            <path
              fill-rule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
              clip-rule="evenodd"
            ></path>
          </svg>
          <span class="text-xs text-gray-300">{{
            connectionStatus === "Disconnected"
              ? "Device disconnected"
              : connectionStatus
          }}</span>
        </div>
      </footer>
    </main>
  </div>
</template>
