<template>
  <div class="pod-page">
    <!-- 页面主题 -->
    <h2 class="title">
      <i class="fa fa-cube icon"></i> 查看 Pod 资源
    </h2>

    <!-- 筛选区域 -->
    <div class="filters">
      <div class="filter-item">
        <label for="ns">命名空间：</label>
        <select id="ns" v-model="selectedNs">
          <option value="">-- 全部 --</option>
          <option v-for="ns in namespaces" :key="ns" :value="ns">{{ ns }}</option>
        </select>
      </div>
      <div class="filter-item">
        <label for="node">节点：</label>
        <select id="node" v-model="selectedNode">
          <option value="">-- 全部 --</option>
          <option v-for="node in nodes" :key="node" :value="node">{{ node }}</option>
        </select>
      </div>
      <div class="filter-item">
        <label for="podName">Pod 名称：</label>
        <input id="podName" type="text" placeholder="请输入 Pod 名称" v-model="podName" />
      </div>
      <div class="filter-item">
        <label for="status">状态：</label>
        <select id="status" v-model="selectedStatus">
          <option value="">-- 全部 --</option>
          <option value="Running">Running</option>
          <option value="Pending">Pending</option>
          <option value="Succeeded">Succeeded</option>
          <option value="Failed">Failed</option>
          <option value="Unknown">Unknown</option>
        </select>
      </div>
    </div>

    <!-- 展示框 -->
    <div class="display-box">
      <div v-if="pagedPods.length === 0" class="placeholder">
        <i class="fa fa-info-circle"></i> 暂无符合条件的 Pod
      </div>
      <div v-else class="data-list">
        <div v-for="(pod, index) in pagedPods" :key="index" class="pod-card-wrapper">
          <div class="pod-card">
            <div class="pod-info">
              <div><i class="fa fa-cube pod-icon"></i> <b>{{ pod.name }}</b></div>
              <div>状态：
                <span :class="pod.status === 'Running' ? 'status-running' : 'status-error'">
                  {{ pod.status }}
                </span>
              </div>
              <div>命名空间：{{ pod.namespace }}</div>
              <div>节点：{{ pod.node }}</div>
            </div>
            <div class="pod-actions">
              <button class="btn btn-danger" @click="deletePod(pod)" :disabled="!isAdmin">
                          <i class="fa fa-trash"></i> 删除
              </button>
              <button class="btn btn-warning" @click="editPod(pod)" :disabled="!isAdmin"><i class="fa fa-edit"></i> 修改</button>
              <button class="btn btn-dark" @click="openShell(pod)" :disabled="!isAdmin">
                <i class="fa fa-terminal"></i> Shell
              </button>
              <button class="btn btn-info" @click="openLogs(pod)" :disabled="!isAdmin">
                <i class="fa fa-file-alt"></i> 日志
              </button>
              <button class="btn btn-detail" @click="showDetails(pod)">
                <i class="fa fa-info-circle"></i> 详情
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination" v-if="totalPages > 1">
      <button :disabled="currentPage === 1" @click="prevPage">上一页</button>
      <span>第 {{ currentPage }} / {{ totalPages }} 页</span>
      <button :disabled="currentPage === totalPages" @click="nextPage">下一页</button>
    </div>

    <!-- Pod 详情模态框 -->
    <div v-if="modalVisible" class="modal-mask">
      <div class="modal-wrapper">
        <div class="modal-container">
          <div class="modal-header">
            <h3>Pod 详情 - {{ modalPod.name }}</h3>
            <button class="modal-close" @click="closeModal">&times;</button>
          </div>
          <p><span class="label">状态：</span>
            <span :class="modalPod.status === 'Running' ? 'status-running' : 'status-error'">
              {{ modalPod.status }}
            </span>
          </p>
          <p><span class="label">命名空间：</span>{{ modalPod.namespace }}</p>
          <p><span class="label">节点：</span>{{ modalPod.node }}</p>
          <div v-if="modalPod.labels">
            <h4>标签</h4>
            <div class="labels">
              <span v-for="(val, key) in modalPod.labels" :key="key" class="label-badge">
                {{ key }} = {{ val }}
              </span>
            </div>
          </div>
          <div v-if="modalPod.containers">
            <h4>容器信息</h4>
            <div v-for="(c, idx) in modalPod.containers" :key="idx" class="container-card">
              <p><span class="label">容器名：</span>{{ c.name }}</p>
              <p><span class="label">镜像：</span><span class="image">{{ c.image }}</span></p>
              <p><span class="label">拉取策略：</span>{{ c.imagePullPolicy }}</p>
              <p><span class="label">端口：</span>{{ c.ports }}</p>
            </div>
          </div>
          <p class="modal-note">提示：弹窗开启时主页面不可操作</p>
        </div>
      </div>
    </div>

    <!-- Shell / Logs 模态框 -->
    <div v-if="terminalVisible" class="modal-mask">
      <div class="terminal-modal-wrapper">
        <div class="terminal-modal-header">
          <h3>
            <i :class="terminalType === 'shell' ? 'fa fa-terminal' : 'fa fa-file-alt'"></i>
            {{ terminalType === 'shell' ? 'Pod Shell' : 'Pod 日志' }} - {{ terminalPod.name }}
          </h3>
          <button class="modal-close" @click="closeTerminal">&times;</button>
        </div>

        <!-- 容器列表操作 -->
        <div class="terminal-info">
          <p><b>Namespace:</b> {{ terminalPod.namespace }}</p>
          <p><b>Pod:</b> {{ terminalPod.name }}</p>
          <div class="container-list">
            <div v-for="(c, idx) in terminalPod.containers" :key="idx" class="container-item">
              <span>{{ c.name }}</span>
              <button class="connect-btn" @click="connectContainer(c.name)">连接</button>
              <button class="btn-clear" @click="clearShell">清空</button>
              <button class="btn-disconnect" @click="disconnectContainer(c.name)">断开</button>
              <span :class="containerStatusClass(c.name)">{{ containerStatusText(c.name) }}</span>
            </div>
          </div>
        </div>

        <!-- Terminal / Logs 显示区 -->
        <div class="terminal-content">
          <div v-for="(c, idx) in terminalPod.containers" :key="'term-'+idx" class="terminal-container-wrapper">
            <div v-show="activeContainer === c.name" class="terminal-container-inner">
              <div v-if="terminalType === 'shell'" class="xterm-container" :ref="'xterm_'+c.name"></div>
              <div v-if="terminalType === 'logs'" class="logs-container">
                <pre>{{ logData }}</pre>
              </div>
            </div>
          </div>
        </div>

      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, computed, onMounted, watch, nextTick } from "vue";
import { useK8sData } from "@/js/useK8sData";
// import axios from "axios";
import axios from "@/main" 
import { Terminal } from 'xterm';
import 'xterm/css/xterm.css';
import { useRouter } from "vue-router";
export default defineComponent({
  name: "AppPodView",
  setup() {
    const { namespaces, nodes } = useK8sData();
    const selectedNs = ref("");
    const selectedNode = ref("");
    const podName = ref("");
    const selectedStatus = ref("");
    const currentPage = ref(1);
    const pageSize = 5;

    const modalVisible = ref(false);
    const modalPod = ref({});
    const pods = ref([]);

    const terminalVisible = ref(false);
    const terminalPod = ref({});
    const terminalType = ref('shell');
    const logData = ref('');
    const activeContainer = ref('');
    const containerSockets = ref({});
    const containerTerms = ref({});
    const containerStatus = ref({});
    
    //新增路由跳转
    const router = useRouter();
    //路由跳转的方法
   const editPod = (pod) => {
  // 1. 通知 WelcomeUser 清空下拉框
  const event = new CustomEvent('clear-welcome-dropdowns');
  window.dispatchEvent(event);

  // 2. 跳转到修改页面
  router.push({
    name: 'UpdatePod',
    query: {
      pod: encodeURIComponent(JSON.stringify(pod))
    }
  });
};
    const fetchPods = async () => {
      try {
        const res = await axios.get("http://192.168.216.50:8090/pod/podalllist");
        if (res.data && res.data.data) pods.value = res.data.data;
      } catch (err) {
        console.error("获取 Pod 数据失败:", err);
      }
    };

    onMounted(fetchPods);
    watch([selectedNs, selectedNode, podName, selectedStatus], () => {
      currentPage.value = 1;
    });

    const filteredPods = computed(() => {
      return pods.value.filter(pod => {
        const ns = (pod.namespace || "").toLowerCase();
        const node = (pod.node || "").toLowerCase();
        const name = (pod.name || "").toLowerCase();
        const status = (pod.status || "").toLowerCase();
        return (selectedNs.value ? ns === selectedNs.value.toLowerCase() : true) &&
               (selectedNode.value ? node === selectedNode.value.toLowerCase() : true) &&
               (podName.value ? name.includes(podName.value.toLowerCase()) : true) &&
               (selectedStatus.value ? status === selectedStatus.value.toLowerCase() : true);
      });
    });

    const totalPages = computed(() => Math.ceil(filteredPods.value.length / pageSize));
    const pagedPods = computed(() => {
      const start = (currentPage.value - 1) * pageSize;
      return filteredPods.value.slice(start, start + pageSize);
    });

    const prevPage = () => { if (currentPage.value > 1) currentPage.value--; };
    const nextPage = () => { if (currentPage.value < totalPages.value) currentPage.value++; };

    const showDetails = (pod) => { modalPod.value = pod; modalVisible.value = true; document.body.style.overflow = 'hidden'; };
    const closeModal = () => { modalVisible.value = false; document.body.style.overflow = ''; };

    const openShell = (pod) => { terminalType.value = 'shell'; openTerminal(pod); };
    const openLogs = (pod) => { terminalType.value = 'logs'; openTerminal(pod); };

    const openTerminal = (pod) => {
      terminalPod.value = pod;
      terminalVisible.value = true;
      document.body.style.overflow = 'hidden';
      logData.value = '';
      activeContainer.value = '';
      containerSockets.value = {};
      containerTerms.value = {};
      containerStatus.value = {};
    };

    const closeTerminal = () => {
      terminalVisible.value = false;
      document.body.style.overflow = '';
      Object.values(containerSockets.value).forEach(ws => ws.close());
      Object.values(containerTerms.value).forEach(t => t.dispose());
      containerSockets.value = {};
      containerTerms.value = {};
      containerStatus.value = {};
      activeContainer.value = '';
      logData.value = '';
    };
    //获取token
const token = localStorage.getItem('token')
if (!token) {
  console.error('没有 token，无法连接 WebSocket')
  return
}
//验证用户名
const username = localStorage.getItem("username") || "";
const isAdmin = computed(() => username === "admin");
    const connectContainer = async (containerName) => {
      Object.keys(containerSockets.value).forEach(name => { if (name !== containerName) disconnectContainer(name); });
      activeContainer.value = containerName;

      await nextTick();
      if (terminalType.value === 'shell') {
        // 初始化 xterm
        if (containerTerms.value[containerName]) containerTerms.value[containerName].dispose();
        const t = new Terminal({ cursorBlink: true, fontSize: 14, theme: { background: "#000", foreground: "#0f0" } });
        containerTerms.value[containerName] = t;
        const containerRef = document.querySelector(`.xterm-container`);
        t.open(containerRef);

        // const ws = new WebSocket(`ws://192.168.216.50:8090/pod/shell/${terminalPod.value.namespace}/${terminalPod.value.name}/${containerName}`);
        const ws = new WebSocket(`ws://192.168.216.50:8090/pod/shell/${terminalPod.value.namespace}/${terminalPod.value.name}/${containerName}?token=${token}`);
        containerSockets.value[containerName] = ws;
        containerStatus.value[containerName] = '连接中...';
        ws.onopen = () => { containerStatus.value[containerName] = '已连接'; };
        ws.onclose = () => { containerStatus.value[containerName] = '已断开'; };
        ws.onerror = () => { containerStatus.value[containerName] = '连接失败'; };
        t.onData(data => { if (ws.readyState === WebSocket.OPEN) ws.send(JSON.stringify({ operation: 'stdin', data })); });
        ws.onmessage = (event) => {
          let msg;
          try { msg = JSON.parse(event.data); } catch { t.write("\r\n*** 返回数据解析错误 ***\r\n"); return; }
          if (msg.operation==='stdout' || msg.operation==='stderr') t.write(msg.data);
        };
      } else if (terminalType.value === 'logs') {
        logData.value = '';
        // const ws = new WebSocket(`ws://192.168.216.50:8090/pod/log/${terminalPod.value.namespace}/${terminalPod.value.name}/${containerName}`);
        const ws = new WebSocket(`ws://192.168.216.50:8090/pod/log/${terminalPod.value.namespace}/${terminalPod.value.name}/${containerName}?token=${token}`);
        containerSockets.value[containerName] = ws;
        containerStatus.value[containerName] = '连接中...';
        ws.onopen = () => { containerStatus.value[containerName] = '已连接'; };
        ws.onclose = () => { containerStatus.value[containerName] = '已断开'; };
        ws.onmessage = (event) => {
          logData.value += event.data + '\n';
        };
      }
    };

    const disconnectContainer = (containerName) => {
      if (containerSockets.value[containerName]) { containerSockets.value[containerName].close(); delete containerSockets.value[containerName]; }
      if (containerTerms.value[containerName]) { containerTerms.value[containerName].dispose(); delete containerTerms.value[containerName]; }
      containerStatus.value[containerName] = '已断开';
    };

    const clearShell = () => {
      if (terminalType.value==='shell' && activeContainer.value && containerTerms.value[activeContainer.value]) { containerTerms.value[activeContainer.value].clear(); }
      if (terminalType.value==='logs') logData.value='';
    };

    const containerStatusText = (name) => containerStatus.value[name] || '未连接';
    const containerStatusClass = (name) => {
      const s = containerStatus.value[name];
      if (s==='已连接') return 'status-connected';
      if (s==='连接中...') return 'status-connecting';
      return 'status-disconnect';
    };
const deletePod = async (pod) => {
  if (!pod || !pod.name || !pod.namespace) return;

  const confirmed = confirm(`确定要删除 Pod: ${pod.name} 吗？`);
  if (!confirmed) return;

  try {
    const res = await axios.post("http://192.168.216.50:8090/pod/Delete", {
      name: pod.name,
      namespace: pod.namespace
    });

    if (res.data && res.data.message) {
      alert(res.data.message);
      // ✅ 删除成功后直接从 pods 数组中移除该 Pod
      const index = pods.value.findIndex(p => p.name === pod.name && p.namespace === pod.namespace);
      if (index !== -1) {
        pods.value.splice(index, 1);
      }

      // 如果当前页已经没有数据，自动翻到上一页
      if (pagedPods.value.length === 0 && currentPage.value > 1) {
        currentPage.value--;
      }
    }
  } catch (err) {
    console.error("删除 Pod 失败:", err);
    alert("删除 Pod 失败：" + (err.response?.data?.message || err.message));
  }
};

    return {
      namespaces, nodes, selectedNs, selectedNode, podName, selectedStatus,
      pagedPods, currentPage, totalPages, prevPage, nextPage,
      modalVisible, modalPod, showDetails, closeModal,
      terminalVisible, terminalPod, terminalType, logData, activeContainer,
      openShell, openLogs, closeTerminal, connectContainer, disconnectContainer,
      clearShell, containerStatusText, containerStatusClass,
       deletePod,editPod,isAdmin  // ✅ 新增
    };
  }
});
</script>
<style scoped>
.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
/* ======================= 页面整体 ======================= */
.pod-page {
  background: #fff;
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
  font-family: "Segoe UI", "Arial", sans-serif;
  color: #333;
}

/* 标题 */
.title {
  font-size: 26px;
  font-weight: 600;
  display: flex;
  align-items: center;
  margin-bottom: 20px;
}
.icon {
  color: #61ddaa;
  margin-right: 10px;
}

/* ======================= 筛选区 ======================= */
.filters {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
  flex-wrap: nowrap;
  align-items: center;
  padding-right: 8px;
}
.filter-item {
  display: flex;
  align-items: center;
  gap: 6px;
}
.filter-item select,
.filter-item input {
  padding: 6px 12px;
  font-size: 15px;
  border-radius: 8px;
  border: 1px solid #ddd;
  background: #fff;
  box-sizing: border-box;
}
.filter-item select#ns,
.filter-item select#node { width: 160px; }
.filter-item input#podName { width: 120px; }
.filter-item select#status { width: 100px; }

/* ======================= 数据区 ======================= */
.display-box {
  border: 1px solid #f0f0f0;
  border-radius: 12px;
  padding: 20px;
  min-height: 300px;
  background: #fafafa;
}
.placeholder {
  font-size: 16px;
  color: #999;
  display: flex;
  align-items: center;
}
.placeholder i { margin-right: 6px; color: #f6c022; }
.data-list { display: flex; flex-direction: column; gap: 15px; }

/* ======================= Pod 卡片 ======================= */
.pod-card {
  background: #fff;
  border: 1px solid #f0f0f0;
  border-radius: 12px;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: all 0.2s ease-in-out;
}
.pod-card:hover {
  box-shadow: 0 2px 10px rgba(64, 158, 255, 0.2);
  transform: translateY(-2px);
}
.pod-info { font-size: 14px; line-height: 1.6; }
.pod-icon { margin-right: 6px; color: #61ddaa; }
.pod-actions { display: flex; gap: 8px; }
.btn {
  padding: 6px 12px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  transition: all 0.2s;
  border: 1px solid transparent;
}
.btn-danger { background: #ffecec; color: #f56c6c; border-color: #f5b7b1; }
.btn-warning { background: #fff4e5; color: #e6a23c; border-color: #f8d7a0; }
.btn-dark { background: #f2f6fc; color: #409eff; border-color: #c6e2ff; }
.btn-info { background: #e6f8f6; color: #17a2b8; border-color: #a9e3ef; }
.btn-detail { background: #eef6ff; color: #3178c6; border-color: #b6d4fe; }

/* 状态 */
.status-running { color: #67c23a; font-weight: 600; }
.status-error { color: #f56c6c; font-weight: 600; }

/* 详情 */
.pod-details {
  background: #f9f9f9;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 16px;
  margin: 10px 0 20px 0;
  font-size: 14px;
  color: #444;
}
.detail-section { margin-bottom: 15px; }
.detail-section h4 {
  font-size: 15px;
  font-weight: 600;
  margin-bottom: 8px;
  color: #3178c6;
}
.label { font-weight: 500; color: #555; }
.labels { display: flex; flex-wrap: wrap; gap: 6px; }
.label-badge {
  background: #eef6ff;
  color: #3178c6;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 13px;
}
.container-card {
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 10px;
  padding: 10px 14px;
  margin-bottom: 10px;
  font-size: 13px;
  color: #333;
}
.container-card .image { color: #409eff; font-weight: 500; }

/* ======================= 分页 ======================= */
.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
  gap: 12px;
  align-items: center;
}
.pagination button {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  background: #fff;
  cursor: pointer;
  transition: all 0.2s;
  color: #409eff;
}
.pagination button:hover:not(:disabled) { background: #ecf5ff; border-color: #409eff; }
.pagination button:disabled { opacity: 0.5; cursor: not-allowed; }

/* ======================= 模态弹窗 ======================= */
.modal-mask {
  position: fixed;
  z-index: 9999;
  top: 0; left: 0;
  width: 100%; height: 100%;
  background-color: rgba(0,0,0,0.6);
  display: flex;
  justify-content: center;
  align-items: center;
  overflow-y: auto;
}
.modal-wrapper {
  width: 600px;
  max-height: 80%;
  background: #fff;
  border-radius: 12px;
  padding: 20px;
  overflow-y: auto;
  box-shadow: 0 6px 20px rgba(0,0,0,0.3);
  position: relative;
}
.modal-container h3 { margin-bottom: 12px; }
.modal-note { font-size: 12px; color: #999; margin-top: 10px; }
.modal-close {
  position: absolute; top: 15px; right: 15px;
  font-size: 20px;
  background: transparent; border: none; cursor: pointer; color: #999;
  transition: color 0.2s;
}
.modal-close:hover { color: #333; }

/* ======================= Shell / Logs 弹窗 ======================= */
.terminal-modal-wrapper {
  width: 80%;
  max-height: 85%;
  min-height: 550px;
  background: #1c1c1c;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-shadow: 0 6px 20px rgba(0,0,0,0.5);
}
.terminal-modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  background: #222;
  color: #fff;
  font-size: 18px;
}
.terminal-info {
  padding: 12px 20px;
  color: #fff;
  background: #2a2a2a;
}
.container-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 6px;
}
.container-item {
  background: #333;
  padding: 6px 12px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  gap: 6px;
}
.connect-btn {
  background: #409eff;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 2px 8px;
  cursor: pointer;
}
.connect-btn:hover { background: #66b1ff; }
.btn-clear {
  background: #f2f6fc;
  border: 1px solid #409eff;
  color: #409eff;
  border-radius: 4px;
  padding: 2px 8px;
  cursor: pointer;
  margin-left: 10px;
}
.btn-clear:hover { background: #66b1ff; color: #fff; }
.btn-disconnect {
  background: #f56c6c;
  color: #fff;
  border: none;
  border-radius: 4px;
  padding: 2px 8px;
  cursor: pointer;
}
.btn-disconnect:hover { background: #ff7b7b; }

.xterm-container {
  width: 100%;
  height: 100%;
}
.logs-container {
  flex: 1;
  min-height: 400px;
  max-height: 600px;
  overflow-y: auto;
  background: #1e1e1e;
  color: #0f0;
  padding: 10px;
  font-family: monospace;
  white-space: pre-wrap;
  border-radius: 6px;
  line-height: 20px;
}
.logs-container pre { margin: 0; }

.status-disconnect { color: #f56c6c; font-weight: bold; margin-right: 8px; }
.status-connecting { color: #e6a23c; font-weight: bold; margin-right: 8px; }
.status-connected { color: #67c23a; font-weight: bold; margin-right: 8px; }

.terminal-status { display: flex; align-items: center; margin-top: 8px; }

</style>




