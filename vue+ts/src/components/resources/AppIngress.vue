<template>
  <div class="ingress-page">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-random icon"></i> 查看 Ingress 资源
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
          <label>Ingress 名称：</label>
          <input type="text" placeholder="请输入 Ingress 名称" v-model="ingName"/>
        </div>
      </div>

      <!-- 列表展示 -->
      <div class="display-box">
        <div v-if="pagedIngress.length === 0" class="placeholder">
          <i class="fa fa-info-circle"></i> 暂无符合条件的 Ingress
        </div>
        <div v-else class="data-list">
          <div v-for="(ing, index) in pagedIngress" :key="index" class="resource-card">
            <div class="resource-info">
              <div><i class="fa fa-random ingress-icon"></i> <b>{{ ing.name }}</b></div>
              <div>命名空间：{{ ing.namespace }}</div>
              <div>IngressClass：{{ ing.ingressClassName || '无' }}</div>
              <div v-if="ing.rules && ing.rules.length">
                <div v-for="(rule, rIdx) in ing.rules" :key="rIdx">
                  主机：{{ rule.host || '无' }}
                  <div v-for="(path, pIdx) in rule.paths" :key="pIdx">
                    路径：{{ path.path }} → {{ path.service }}:{{ path.port }}
                  </div>
                </div>
              </div>
            </div>
            <div class="resource-actions">
              <button class="btn btn-danger" @click="deleteIngress(ing)" :disabled="!isAdmin">
                <i class="fa fa-trash"></i> 删除
              </button>
              <button class="btn btn-warning" @click="editIngress(ing)" :disabled="!isAdmin"><i class="fa fa-edit"></i> 修改</button>
              <button class="btn btn-detail" @click="openModal(ing)">
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

    <!-- 详情弹窗 -->
    <div v-if="showModal" class="modal-mask">
      <div class="modal-wrapper">
        <div class="modal-container">
          <div class="modal-header">
            <h3>Ingress 详情 - {{ modalIngress.name }}</h3>
            <button class="modal-close" @click="closeModal">&times;</button>
          </div>

          <p><span class="label">API 版本：</span>{{ modalIngress.apiVersion }}</p>
          <p><span class="label">类型：</span>{{ modalIngress.kind }}</p>
          <p><span class="label">名称：</span>{{ modalIngress.name }}</p>
          <p><span class="label">命名空间：</span>{{ modalIngress.namespace }}</p>
          <p><span class="label">IngressClass：</span>{{ modalIngress.ingressClassName || '无' }}</p>

          <div v-if="modalIngress.labels && Object.keys(modalIngress.labels).length">
            <h4>标签</h4>
            <div class="labels">
              <span v-for="(val, key) in modalIngress.labels" :key="key" class="label-badge">
                {{ key }} = {{ val }}
              </span>
            </div>
          </div>

          <div v-if="modalIngress.annotations && Object.keys(modalIngress.annotations).length">
            <h4>注解</h4>
            <div class="labels">
              <span v-for="(val, key) in modalIngress.annotations" :key="key" class="label-badge">
                {{ key }} = {{ val }}
              </span>
            </div>
          </div>

          <!-- 路由规则 -->
          <div v-if="modalIngress.rules && modalIngress.rules.length">
            <h4>规则</h4>
            <div v-for="(rule, rIdx) in modalIngress.rules" :key="rIdx" class="container-card">
              <p><span class="label">Host：</span>{{ rule.host || '无' }}</p>
              <div v-for="(path, pIdx) in rule.paths" :key="pIdx" class="labels">
                <span class="label-badge">路径：{{ path.path }}</span>
                <span class="label-badge">服务：{{ path.service }}</span>
                <span class="label-badge">端口：{{ path.port }}</span>
              </div>
            </div>
          </div>

          <!-- TLS -->
          <div v-if="modalIngress.tls && modalIngress.tls.length">
            <h4>TLS</h4>
            <div v-for="(tls, tIdx) in modalIngress.tls" :key="tIdx" class="labels">
              <span class="label-badge">Secret：{{ tls.secretName }}</span>
              <span v-for="host in tls.hosts || []" :key="host" class="label-badge">{{ host }}</span>
            </div>
          </div>

          <!-- 状态 -->
          <div v-if="modalIngress.status && modalIngress.status.length">
            <h4>状态</h4>
            <div v-for="(st, sIdx) in modalIngress.status" :key="sIdx" class="labels">
              <span class="label-badge">IP：{{ st.ip || '-' }}</span>
              <span class="label-badge">Hostname：{{ st.hostname || '-' }}</span>
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
  name: "AppIngressView",
  setup() {
    const username = localStorage.getItem("username") || "";
    const isAdmin = computed(() => username === "admin");
    const { namespaces } = useK8sData()
    const selectedNs = ref("")
    const ingName = ref("")
    const currentPage = ref(1)
    const pageSize = 5
    const ingress = ref([])

    const showModal = ref(false)
    const modalIngress = ref({})
    //新增路由跳转
    const router = useRouter();
const editIngress = (ingress) => {
  const event = new CustomEvent('clear-welcome-dropdowns');
  window.dispatchEvent(event);
  router.push({
    name: 'UpdateIngress',
    query: { ingress: encodeURIComponent(JSON.stringify(ingress)) }
  });
};
    // 获取 Ingress 数据
    const fetchIngress = async () => {
      try {
        const res = await axios.get("http://192.168.216.50:8090/ingest/GetIngressAll")
        if (res.data && res.data.date) {
          ingress.value = res.data.date.map(i => ({
            apiVersion: i.apiVersion,
            kind: i.kind,
            name: i.name,
            namespace: i.namespace,
            labels: i.labels || {},
            annotations: i.annotations || {},
            ingressClassName: i.ingressClassName || "",
            rules: i.rules || [],
            tls: i.tls || [],
            status: i.status || []
          }))
        }
      } catch (err) {
        console.error("获取 Ingress 数据失败:", err)
      }
    }

    onMounted(fetchIngress)

    // 删除 Ingress
    const deleteIngress = async (ing) => {
      if (!confirm(`确定删除 Ingress "${ing.name}" 吗？`)) return
      try {
        const res = await axios.post("http://192.168.216.50:8090/ingest/IngressDelete", {
          namespace: ing.namespace,
          name: ing.name
        })
        if (res.data && res.data.code === 0) {
          alert(`Ingress "${ing.name}" 删除成功！`)

          // 前端立即移除
          ingress.value = ingress.value.filter(item => !(item.namespace === ing.namespace && item.name === ing.name))

          // 如果当前页没有数据且不是第一页，自动跳到上一页
          if (pagedIngress.value.length === 0 && currentPage.value > 1) {
            currentPage.value--
          }

        } else {
          alert(`删除失败: ${res.data.message || '未知错误'}`)
        }
      } catch (err) {
        console.error("删除 Ingress 失败:", err)
        alert("删除失败，请检查控制台信息")
      }
    }

    // 过滤
    const filtered = computed(() =>
      ingress.value.filter(i => {
        const matchNs = selectedNs.value ? i.namespace === selectedNs.value : true
        const matchName = ingName.value ? i.name.includes(ingName.value) : true
        return matchNs && matchName
      })
    )

    // 筛选条件变化时重置页码
    watch([selectedNs, ingName], () => { currentPage.value = 1 })

    // 分页
    const totalPages = computed(() => Math.ceil(filtered.value.length / pageSize))
    const pagedIngress = computed(() => {
      const start = (currentPage.value - 1) * pageSize
      return filtered.value.slice(start, start + pageSize)
    })

    const prevPage = () => { if (currentPage.value > 1) currentPage.value-- }
    const nextPage = () => { if (currentPage.value < totalPages.value) currentPage.value++ }

    const openModal = (ing) => { modalIngress.value = ing; showModal.value = true }
    const closeModal = () => { showModal.value = false; modalIngress.value = {} }

    return {
      namespaces,
      selectedNs,
      ingName,
      pagedIngress,
      currentPage,
      totalPages,
      prevPage,
      nextPage,
      showModal,
      modalIngress,
      openModal,
      closeModal,
      deleteIngress,
      editIngress,
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
.title>.icon{ color:#f66d6d }
.ingress-icon{ color:#f66d6d }

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
.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.modal-header h3 { margin-bottom: 12px; }
.modal-close {
  position: absolute; top: 15px; right: 15px;
  font-size: 20px;
  background: transparent; border: none; cursor: pointer; color: #999;
  transition: color 0.2s;
}
.modal-close:hover { color: #333; }
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
.modal-note { font-size: 12px; color: #999; margin-top: 10px; }

.btn-detail { background: #eef6ff; color: #3178c6; border-color: #b6d4fe; }
</style>

<style src="@/css/ResourceCommon.css"></style>
