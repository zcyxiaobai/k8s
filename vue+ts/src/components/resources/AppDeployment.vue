<template>
  <div class="deploy-page">
    <div class="page-container">
      <h2 class="title">
        <i class="fa fa-cubes icon"></i> 查看 Deployment 资源
      </h2>

      <div class="filters">
        <div class="filter-item">
          <label for="ns">命名空间：</label>
          <select id="ns" v-model="selectedNs">
            <option value="">-- 全部 --</option>
            <option v-for="ns in namespaces" :key="ns" :value="ns">{{ ns }}</option>
          </select>
        </div>

        <div class="filter-item">
          <label for="deployName">Deployment 名称：</label>
          <input id="deployName" type="text" placeholder="请输入 Deployment 名称" v-model="deployName"/>
        </div>
      </div>

      <div class="display-box">
        <div v-if="pagedDeployments.length === 0" class="placeholder">
          <i class="fa fa-info-circle"></i> 暂无符合条件的 Deployment
        </div>
        <div v-else class="data-list">
          <div v-for="(deploy, index) in pagedDeployments" :key="index" class="resource-card">
            <div class="resource-info">
              <div><i class="fa fa-cubes deployment-icon"></i> <b>{{ deploy.name }}</b></div>
              <div>命名空间：{{ deploy.namespace }}</div>
              <div>副本数：{{ deploy.replicas }}</div>
            </div>
            <div class="resource-actions">
              <button class="btn btn-danger" @click="deleteDeployment(deploy)" :disabled="!isAdmin">
                <i class="fa fa-trash"></i> 删除
              </button>
              <button class="btn btn-warning" @click="editDeployment(deploy)" :disabled="!isAdmin"><i class="fa fa-edit"></i> 修改</button>
              <button class="btn btn-detail" @click="openModal(deploy)">
                <i class="fa fa-info-circle"></i> 详情
              </button>
            </div>
          </div>
        </div>
      </div>

      <div class="pagination" v-if="totalPages > 1">
        <button :disabled="currentPage === 1" @click="prevPage">上一页</button>
        <span>第 {{ currentPage }} / {{ totalPages }} 页</span>
        <button :disabled="currentPage === totalPages" @click="nextPage">下一页</button>
      </div>
    </div>

    <!-- Deployment 详情弹窗 -->
    <div v-if="showModal" class="modal-mask">
      <div class="modal-wrapper">
        <div class="modal-container">
          <div class="modal-header">
            <h3>Deployment 详情 - {{ modalDeploy.name }}</h3>
            <button class="modal-close" @click="closeModal">&times;</button>
          </div>

          <p><span class="label">API 版本：</span>{{ modalDeploy.apiVersion }}</p>
          <p><span class="label">类型：</span>{{ modalDeploy.kind }}</p>
          <p><span class="label">名称：</span>{{ modalDeploy.name }}</p>
          <p><span class="label">命名空间：</span>{{ modalDeploy.namespace }}</p>
          <p><span class="label">副本数：</span>{{ modalDeploy.replicas }}</p>
          <p><span class="label">Image Pull Secrets：</span>{{ modalDeploy.imagePullSecrets || '无' }}</p>

          <div v-if="modalDeploy.labels && Object.keys(modalDeploy.labels).length">
            <h4>容器标签</h4>
            <div class="labels">
              <span v-for="(val, key) in modalDeploy.labels" :key="key" class="label-badge">
                {{ key }} = {{ val }}
              </span>
            </div>
          </div>

          <div v-if="modalDeploy.selector && Object.keys(modalDeploy.selector).length">
            <h4>Selector</h4>
            <div class="labels">
              <span v-for="(val, key) in modalDeploy.selector" :key="key" class="label-badge">
                {{ key }} = {{ val }}
              </span>
            </div>
          </div>

          <div v-if="modalDeploy.containers && modalDeploy.containers.length">
            <h4>容器信息</h4>
            <div v-for="(c, idx) in modalDeploy.containers" :key="idx" class="container-card">
              <p><span class="label">容器名：</span>{{ c.name }}</p>
              <p><span class="label">镜像：</span><span class="image">{{ c.image }}</span></p>
              <p><span class="label">拉取策略：</span>{{ c.imagePullPolicy || '默认' }}</p>
              <p v-if="c.ports && c.ports.length"><span class="label">端口：</span>{{ c.ports.join(', ') }}</p>
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
  name: "AppDeploymentView",
  setup() {
    const username = localStorage.getItem("username") || "";
    const isAdmin = computed(() => username === "admin");
    const { namespaces } = useK8sData()
    const selectedNs = ref("")
    const deployName = ref("")
    const currentPage = ref(1)
    const pageSize = 5
    const deployments = ref([])

    const showModal = ref(false)
    const modalDeploy = ref({})
     const router = useRouter() // ✅ 获取路由实例
const editDeployment = (deployment) => {
  // 1. 清空 WelcomeUser 左侧下拉框
  const event = new CustomEvent('clear-welcome-dropdowns');
  //验证用户名
  
  window.dispatchEvent(event);

  // 2. 跳转到对应修改页面，传递数据
  router.push({
    name: 'UpdateDeployment', // 注意路由名字
    query: { deployment: encodeURIComponent(JSON.stringify(deployment)) }
  });
};
    const fetchDeployments = async () => {
      try {
        const res = await axios.get("http://192.168.216.50:8090/dep/GetDeploymentList")
        if (res.data && res.data.date) {
          deployments.value = res.data.date.map(d => ({
            apiVersion: d.apiVersion,
            kind: d.kind,
            name: d.name,
            namespace: d.namespace,
            replicas: d.replicas,
            imagePullSecrets: d.imagePullSecrets,
            containers: d.containers || [],
            labels: d.labels || {},
            selector: d.selector || {}
          }))
        }
      } catch (err) {
        console.error("获取 Deployment 数据失败:", err)
      }
    }

    onMounted(fetchDeployments)

    const filtered = computed(() =>
      deployments.value.filter(d => {
        const matchNs = selectedNs.value ? d.namespace === selectedNs.value : true
        const matchName = deployName.value ? d.name.includes(deployName.value) : true
        return matchNs && matchName
      })
    )

    watch([selectedNs, deployName], () => { currentPage.value = 1 })

    const totalPages = computed(() => Math.ceil(filtered.value.length / pageSize))
    const pagedDeployments = computed(() => {
      const start = (currentPage.value - 1) * pageSize
      return filtered.value.slice(start, start + pageSize)
    })

    const prevPage = () => { if (currentPage.value > 1) currentPage.value-- }
    const nextPage = () => { if (currentPage.value < totalPages.value) currentPage.value++ }

    const openModal = (deploy) => {
      modalDeploy.value = deploy
      showModal.value = true
    }
    const closeModal = () => {
      showModal.value = false
      modalDeploy.value = {}
    }

    // 删除功能
    const deleteDeployment = async (deploy) => {
      if (!confirm(`确定要删除 Deployment "${deploy.name}" 吗？`)) return
      try {
        await axios.post("http://192.168.216.50:8090/dep/DeleteDev", {
          namespace: deploy.namespace,
          name: deploy.name
        })
        alert(`Deployment "${deploy.name}" 删除成功！`)
        // 删除后刷新数据
        fetchDeployments()
      } catch (err) {
        console.error("删除 Deployment 失败:", err)
        alert(`删除 Deployment "${deploy.name}" 失败！`)
      }
    }

    return {
      namespaces,
      selectedNs,
      deployName,
      pagedDeployments,
      currentPage,
      totalPages,
      prevPage,
      nextPage,
      showModal,
      modalDeploy,
      openModal,
      closeModal,
      deleteDeployment,
      editDeployment,
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
.btn-detail { background: #eef6ff; color: #3178c6; border-color: #b6d4fe; }
.title>.icon { color: #65789b; margin-right: 10px; }
.deployment-icon { color: #65789b; }

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

/* 基础标签样式 */
.label { font-weight: 500; color: #555; }
.status-running { color: green; }
.status-error { color: red; }

/* 标签容器 */
.labels { display: flex; flex-wrap: wrap; gap: 6px; }
.label-badge {
  background: #eef6ff;
  color: #3178c6;
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 13px;
}

/* 容器信息卡片 */
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
.image { color: #555; }

/* 弹窗底部提示 */
.modal-note { font-size: 12px; color: #999; margin-top: 10px; }
</style>

<style src="@/css/ResourceCommon.css"></style>
