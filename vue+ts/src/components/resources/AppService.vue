<template>
  <div class="svc-page">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-plug icon"></i> 查看 Service 资源
      </h2>

      <!-- 筛选条件 -->
      <div class="filters">
        <div class="filter-item">
          <label>命名空间：</label>
          <select v-model="selectedNs">
            <option value="">-- 全部 --</option>
            <option v-for="ns in namespaces" :key="ns" :value="ns">{{ ns }}</option>
          </select>
        </div>

        <div class="filter-item">
          <label>Service 名称：</label>
          <input type="text" placeholder="请输入 Service 名称" v-model="svcName"/>
        </div>
      </div>

      <!-- 数据列表 -->
      <div class="display-box">
        <div v-if="pagedServices.length === 0" class="placeholder">
          <i class="fa fa-info-circle"></i> 暂无符合条件的 Service
        </div>
        <div v-else class="data-list">
          <div v-for="(svc, index) in pagedServices" :key="index" class="resource-card">
            <div class="resource-info">
              <div><i class="fa fa-plug service-icon"></i> <b>{{ svc.name }}</b></div>
              <div>命名空间：{{ svc.namespace }}</div>
              <div>类型：{{ svc.type }}</div>
              <div v-if="svc.clusterIP">ClusterIP：{{ svc.clusterIP }}</div>
            </div>
            <div class="resource-actions">
              <!-- 删除按钮 -->
              <button class="btn btn-danger" @click="deleteService(svc)" :disabled="!isAdmin">
                <i class="fa fa-trash"></i> 删除
              </button>
              <button class="btn btn-warning" @click="editService(svc)" :disabled="!isAdmin"><i class="fa fa-edit"></i> 修改</button>
              <button class="btn btn-detail" @click="openModal(svc)">
                <i class="fa fa-info-circle"></i> 详情
              </button>
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
    </div>

    <!-- Service 详情弹窗 -->
    <div v-if="showModal" class="modal-mask">
      <div class="modal-wrapper">
        <div class="modal-container">
          <!-- 弹窗头部 -->
          <div class="modal-header">
            <h3>Service 详情 - {{ modalSvc.name }}</h3>
            <button class="modal-close" @click="closeModal">&times;</button>
          </div>

          <!-- 基本信息 -->
          <p><span class="label">API 版本：</span>{{ modalSvc.apiVersion }}</p>
          <p><span class="label">类型：</span>{{ modalSvc.kind }}</p>
          <p><span class="label">名称：</span>{{ modalSvc.name }}</p>
          <p><span class="label">命名空间：</span>{{ modalSvc.namespace }}</p>
          <p><span class="label">Service 类型：</span>{{ modalSvc.type }}</p>
          <p v-if="modalSvc.clusterIP"><span class="label">ClusterIP：</span>{{ modalSvc.clusterIP }}</p>

          <!-- Selector -->
          <div v-if="modalSvc.selector && Object.keys(modalSvc.selector).length">
            <h4>Selector</h4>
            <div class="labels">
              <span v-for="(val, key) in modalSvc.selector" :key="key" class="label-badge">
                {{ key }} = {{ val }}
              </span>
            </div>
          </div>
<!-- Labels -->
<div v-if="modalSvc.labels && Object.keys(modalSvc.labels).length">
  <h4>Labels</h4>
  <div class="labels">
    <span v-for="(val, key) in modalSvc.labels" :key="key" class="label-badge">
      {{ key }} = {{ val }}
    </span>
  </div>
</div>
          <!-- Ports -->
          <div v-if="modalSvc.ports && modalSvc.ports.length">
            <h4>端口信息</h4>
            <div v-for="(p, idx) in modalSvc.ports" :key="idx" class="container-card">
              <p><span class="label">Port：</span>{{ p.port }}</p>
              <p><span class="label">TargetPort：</span>{{ p.targetPort }}</p>
              <p><span class="label">Protocol：</span>{{ p.protocol }}</p>
              <p v-if="p.nodePort"><span class="label">NodePort：</span>{{ p.nodePort }}</p>
            </div>
          </div>

          <p class="modal-note">提示：弹窗开启时主页面不可操作</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, computed, onMounted, watch } from "vue"
import { useK8sData } from "@/js/useK8sData"
// import axios from "axios"
import axios from "@/main" 
import { useRouter } from "vue-router";
export default defineComponent({
  name: "AppServiceView",
  setup() {
    const username = localStorage.getItem("username") || "";
    const isAdmin = computed(() => username === "admin");
    const { namespaces } = useK8sData()
    const selectedNs = ref("")
    const svcName = ref("")
    const currentPage = ref(1)
    const pageSize = 5
    const services = ref([])

    const showModal = ref(false)
    const modalSvc = ref({})
        //新增路由跳转
    const router = useRouter();
const editService = (service) => {
  const event = new CustomEvent('clear-welcome-dropdowns');
  window.dispatchEvent(event);
  router.push({
    name: 'UpdateService',
    query: { service: encodeURIComponent(JSON.stringify(service)) }
  });
};

    // 获取 Service 数据
    const fetchServices = async () => {
      try {
        const res = await axios.get("http://192.168.216.50:8090/service/GetServiceInfo")
        if (res.data && res.data.date) {
          services.value = res.data.date.map(s => ({
            apiVersion: s.apiVersion,
            kind: s.kind,
            name: s.name,
            namespace: s.namespace,
            type: s.type,
            clusterIP: s.clusterIP,
            selector: s.selector || {},
            ports: s.ports || [],
             labels: s.labels || {}
          }))
        }
      } catch (err) {
        console.error("获取 Service 数据失败:", err)
      }
    }

    onMounted(fetchServices)

    // 删除 Service
    const deleteService = async (svc) => {
      if (!confirm(`确定要删除 Service "${svc.name}" 吗？`)) return
      try {
        await axios.post("http://192.168.216.50:8090/service/DeleteService", {
          namespace: svc.namespace,
          name: svc.name
        })
        alert(`Service "${svc.name}" 删除成功`)
        fetchServices()  // 删除成功后刷新列表
      } catch (err) {
        console.error("删除 Service 失败:", err)
        alert(`删除 Service "${svc.name}" 失败: ${err.message}`)
      }
    }

    // 过滤数据
    const filtered = computed(() =>
      services.value.filter(s => {
        const matchNs = selectedNs.value ? s.namespace === selectedNs.value : true
        const matchName = svcName.value ? s.name.includes(svcName.value) : true
        return matchNs && matchName
      })
    )

    // 筛选条件变化时重置页码
    watch([selectedNs, svcName], () => { currentPage.value = 1 })

    // 分页
    const totalPages = computed(() => Math.ceil(filtered.value.length / pageSize))
    const pagedServices = computed(() => {
      const start = (currentPage.value - 1) * pageSize
      return filtered.value.slice(start, start + pageSize)
    })

    // 分页控制
    const prevPage = () => { if (currentPage.value > 1) currentPage.value-- }
    const nextPage = () => { if (currentPage.value < totalPages.value) currentPage.value++ }

    // 弹窗控制
    const openModal = (svc) => {
      modalSvc.value = svc
      showModal.value = true
    }
    const closeModal = () => {
      showModal.value = false
      modalSvc.value = {}
    }

    return {
      namespaces,
      selectedNs,
      svcName,
      pagedServices,
      currentPage,
      totalPages,
      prevPage,
      nextPage,
      showModal,
      modalSvc,
      openModal,
      closeModal,
      deleteService,
      editService,
      isAdmin
    }
  }
})
</script>

<style scoped>
.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}
.title>.icon { color:#f6c022 }
.service-icon{ color:#f6c022 }
.btn-detail { background: #eef6ff; color: #3178c6; border-color: #b6d4fe; }
/* 遮罩层 */
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

/* 弹窗容器 */
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

/* 弹窗头部 */
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.modal-header h3 { margin-bottom: 12px; }

/* 关闭按钮 */
.modal-close {
  position: absolute; top: 15px; right: 15px;
  font-size: 20px;
  background: transparent; border: none; cursor: pointer; color: #999;
  transition: color 0.2s;
}
.modal-close:hover { color: #333; }

/* 标签样式 */
.label { font-weight: 500; color: #555; }
.labels { display: flex; flex-wrap: wrap; gap: 6px; }
.label-badge {
  background: #eef6ff;
  color: #3178c6;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 13px;
}

/* 端口信息卡片 */
.container-card {
  background: #fff;
  border: 1px solid #ddd;
  border-radius: 10px;
  padding: 10px 14px;
  margin-bottom: 10px;
  font-size: 13px;
  color: #333;
}

/* 弹窗底部提示 */
.modal-note { font-size: 12px; color: #999; margin-top: 10px; }
</style>

<style src="@/css/ResourceCommon.css"></style>
