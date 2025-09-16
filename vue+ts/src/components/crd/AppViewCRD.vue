<template>
  <div class="bookapp-page">
    <!-- 页面主题 -->
    <h2 class="title">
      <i class="fa fa-book icon"></i> 查看 Bookapp 资源
    </h2>

    <!-- 筛选区域 -->
    <div class="filters">
      <div class="filter-item">
        <label for="name">名称：</label>
        <input id="name" type="text" placeholder="请输入名称" v-model="searchName" />
      </div>
      <div class="filter-item">
        <label for="namespace">命名空间：</label>
        <select id="namespace" v-model="selectedNs">
          <option value="">-- 全部 --</option>
          <option v-for="ns in namespaces" :key="ns" :value="ns">{{ ns }}</option>
        </select>
      </div>
      <div class="filter-item">
        <label for="serverName">域名：</label>
        <input id="serverName" type="text" placeholder="请输入域名" v-model="searchServerName" />
      </div>
    </div>

    <!-- 展示框 -->
    <div class="display-box">
      <div v-if="pagedBookapps.length === 0" class="placeholder">
        <i class="fa fa-info-circle"></i> 暂无符合条件的 Bookapp 资源
      </div>
      <div v-else class="data-list">
        <div v-for="(item, index) in pagedBookapps" :key="index" class="bookapp-card-wrapper">
          <div class="bookapp-card">
            <div class="bookapp-info">
              <div><i class="fa fa-cube bookapp-icon"></i> <b>{{ item.name }}</b></div>
              <div>命名空间：{{ item.namespace }}</div>
              <div>镜像：<span class="image">{{ item.spec?.image }}</span></div>
              <div>副本数：{{ item.spec?.size }}</div>
              <div>域名：{{ item.spec?.serverName }}</div>
              <div>
                状态：
                <span class="status-running">运行中 {{ item.status?.runing }}</span> /
                <span class="status-error">未运行 {{ item.status?.notRuning }}</span>
              </div>
            </div>
            <div class="bookapp-actions">
              <button class="btn btn-danger" @click="deleteBookapp(item)" :disabled="!isAdmin">
                  <i class="fa fa-trash"></i> 删除
              </button>
              <button class="btn btn-warning" @click="editBookapp(item)" :disabled="!isAdmin">
                   <i class="fa fa-edit"></i> 修改
              </button>
              <button class="btn btn-detail" @click="showDetails(item)">
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

    <!-- 详情模态框 -->
    <div v-if="modalVisible" class="modal-mask">
      <div class="modal-wrapper">
        <div class="modal-container">
          <div class="modal-header">
            <h3>Bookapp 详情 - {{ modalData.name }}</h3>
            <button class="modal-close" @click="closeModal">&times;</button>
          </div>
          <p><span class="label">命名空间：</span>{{ modalData.namespace }}</p>
          <p><span class="label">镜像：</span>{{ modalData.spec?.image }}</p>
          <p><span class="label">副本数：</span>{{ modalData.spec?.size }}</p>
          <p><span class="label">域名：</span>{{ modalData.spec?.serverName }}</p>
          <div v-if="modalData.labels">
            <h4>标签</h4>
            <div class="labels">
              <span v-for="(val, key) in modalData.labels" :key="key" class="label-badge">
                {{ key }} = {{ val }}
              </span>
            </div>
          </div>
          <div v-if="modalData.annotations">
            <h4>注解</h4>
            <pre class="annotations">{{ modalData.annotations }}</pre>
          </div>
          <div>
            <h4>端口信息</h4>
            <p><span class="label">端口名：</span>{{ modalData.spec?.port?.name }}</p>
            <p><span class="label">协议：</span>{{ modalData.spec?.port?.protocol }}</p>
            <p><span class="label">端口：</span>{{ modalData.spec?.port?.port }}</p>
            <p><span class="label">目标端口：</span>{{ modalData.spec?.port?.targetPort }}</p>
          </div>
          <div>
            <h4>状态</h4>
            <p><span class="label">运行中：</span>{{ modalData.status?.runing }}</p>
            <p><span class="label">未运行：</span>{{ modalData.status?.notRuning }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, onMounted, computed, watch } from "vue";
// import axios from "axios";
import axios from "@/main" 
import { useK8sData } from "@/js/useK8sData"
import { useRouter } from "vue-router";
export default defineComponent({
  name: "AppBookappView",
  setup() {
    const username = ref(localStorage.getItem("username") || " ");
const isAdmin = computed(() => username.value === "admin");
    const router = useRouter();
    const bookapps = ref([])
  const modalVisible = ref(false)
  const modalData = ref({})
   // ✅ 使用 useK8sData 获取 namespaces
   const { namespaces, fetchAll } = useK8sData()
   const selectedNs = ref("")
    const searchName = ref("")
  const searchServerName = ref("")
  const currentPage = ref(1)
  const pageSize = 5

    // const fetchNamespaces = async () => {
    //   try {
    //     const res = await axios.get("http://192.168.216.50:8090/pod/GetNamespace");
    //     if (res.data && res.data.data) namespaces.value = res.data.data;
    //   } catch (err) {
    //     console.error("获取命名空间失败:", err);
    //   }
    // };
    const editBookapp = (item) => {
  if (!item) return;

  // 可选：通知 WelcomeUser 清空下拉框（和 Pod 类似）
  const event = new CustomEvent('clear-welcome-dropdowns');
  window.dispatchEvent(event);

  // 跳转到修改页面，并传递 bookapp 数据
  router.push({
    name: 'UpdateBookapp',
    query: {
      bookapp: encodeURIComponent(JSON.stringify(item))
    }
  });
};

    //删除的方法
    // 删除 Bookapp
const deleteBookapp = async (item) => {
  if (!confirm(`确定要删除 Bookapp "${item.name}" 吗？`)) return;

  try {
    await axios.post("http://192.168.216.50:8090/bookapp/deletebookapp", {
      namespace: item.namespace,
      name: item.name
    });

    // 删除成功后，从本地列表移除
    const index = bookapps.value.findIndex(b => b.name === item.name && b.namespace === item.namespace);
    if (index !== -1) bookapps.value.splice(index, 1);

    alert(`Bookapp "${item.name}" 删除成功`);
  } catch (err) {
    console.error("删除失败:", err);
    alert(`删除失败: ${err.response?.data?.message || err.message}`);
  }
};

     const fetchData = async () => {
    try {
      const res = await axios.get("http://192.168.216.50:8090/bookapp/bookapps")
      //  const res = await axios.get("/bookapp/bookapps")
      if (res.data && res.data.data) bookapps.value = res.data.data
    } catch (err) {
      console.error("获取 Bookapp 数据失败:", err)
    }
  }

    // onMounted(() => {
    //   fetchNamespaces();
    //   fetchData();
    // });
      onMounted(() => {
    fetchAll()       // ✅ 调用 useK8sData 获取 namespaces
    fetchData()      // 获取 bookapp 数据
  })

    const filteredBookapps = computed(() => {
      return bookapps.value.filter(item => {
        const nameMatch = item.name?.toLowerCase().includes(searchName.value.toLowerCase());
        const nsMatch = selectedNs.value ? item.namespace === selectedNs.value : true;
        const serverMatch = item.spec?.serverName?.toLowerCase().includes(searchServerName.value.toLowerCase());
        return nameMatch && nsMatch && serverMatch;
      });
    });

    const totalPages = computed(() => Math.ceil(filteredBookapps.value.length / pageSize));
    const pagedBookapps = computed(() => {
      const start = (currentPage.value - 1) * pageSize;
      return filteredBookapps.value.slice(start, start + pageSize);
    });

    const prevPage = () => { if (currentPage.value > 1) currentPage.value--; };
    const nextPage = () => { if (currentPage.value < totalPages.value) currentPage.value++; };

    watch([searchName, searchServerName, selectedNs], () => { currentPage.value = 1; });

    const showDetails = (item) => {
      modalData.value = item;
      modalVisible.value = true;
      document.body.style.overflow = "hidden";
    };
    const closeModal = () => { modalVisible.value = false; document.body.style.overflow = ""; };

    return {
      bookapps, modalVisible, modalData,
      namespaces, selectedNs, searchName, searchServerName,
      currentPage, totalPages, pagedBookapps,
      prevPage, nextPage, showDetails, closeModal,deleteBookapp,
  editBookapp, isAdmin  // ✅ 返回跳转函数
    };
  }
});
</script>

<style scoped>
button:disabled,
.btn[disabled] {
 opacity: 0.5;
  cursor: not-allowed;
 
}
.bookapp-page {
  background: #fff;
  padding: 24px;
  border-radius: 16px;
  box-shadow: 0 4px 16px rgba(0,0,0,0.05);
  font-family: "Segoe UI","Arial",sans-serif;
  color: #333;
}
.title { font-size: 26px; font-weight: 600; display: flex; align-items: center; margin-bottom: 20px; }
.icon { color: #409eff; margin-right: 10px; }
.filters {
  display: flex;           /* 横向排列 */
  gap: 20px;               /* 各筛选项间距 */
  align-items: flex-end;   /* 标签底部对齐输入框 */
  flex-wrap: wrap;         /* 屏幕窄时自动换行 */
}
.filter-item {
  display: flex;            /* 内部横向排列 label 和 input/select */
  align-items: center;
  gap: 6px;                 /* label 与输入框间距 */
  font-size: 13px;
}
.filter-item input,
.filter-item select {
  width: 160px;
  padding: 6px 12px;
  font-size: 15px;
  border-radius: 8px;
  border: 1px solid #ddd;
  outline: none;
  transition: all 0.2s;
}
.filter-item select { width: 160px; } .filter-item input { width: 160px; }
.filter-item input:focus,
.filter-item select:focus {
  border-color: #409eff;        /* 蓝色边框 */
  box-shadow: 0 0 5px rgba(64,158,255,0.4); /* 高亮阴影 */
}
.filter-item label {
  white-space: nowrap;
  font-weight: 500;
  color: #555;
}

.filter-item input,
.filter-item select {
  padding: 6px 12px;
  font-size: 15px;
  border-radius: 8px;
  border: 1px solid #ddd;
  outline: none;
  transition: all 0.2s;
}
.display-box { border: 1px solid #f0f0f0; border-radius: 12px; padding: 20px; min-height: 300px; background: #fafafa; }
.placeholder { font-size: 16px; color: #999; display: flex; align-items: center; } .placeholder i { margin-right: 6px; color: #f6c022; }
.data-list { display: flex; flex-direction: column; gap: 15px; }

.bookapp-card { background: #fff; border: 1px solid #f0f0f0; border-radius: 12px; padding: 15px 20px; display: flex; justify-content: space-between; align-items: center; transition: all 0.2s ease-in-out; }
.bookapp-card:hover { box-shadow: 0 2px 10px rgba(64,158,255,0.2); transform: translateY(-2px); }
.bookapp-info { font-size: 14px; line-height: 1.6; }
.bookapp-icon { margin-right: 6px; color: #409eff; }
.bookapp-actions { display: flex; gap: 8px; }
.btn { padding: 6px 12px; border-radius: 6px; cursor: pointer; font-size: 13px; border: 1px solid transparent; transition: all 0.2s; }
.btn-danger { background: #ffecec; color: #f56c6c; border-color: #f5b7b1; }
.btn-warning { background: #fff4e5; color: #e6a23c; border-color: #f8d7a0; }
.btn-detail { background: #eef6ff; color: #3178c6; border-color: #b6d4fe; }

.status-running { color: #67c23a; font-weight: 600; margin-right: 6px; }
.status-error { color: #f56c6c; font-weight: 600; }

.pagination { display: flex; justify-content: center; align-items: center; gap: 12px; margin-top: 20px; }
.pagination button { padding: 6px 12px; border: 1px solid #ddd; border-radius: 6px; background: #fff; cursor: pointer; transition: all 0.2s; color: #409eff; }
.pagination button:hover:not(:disabled) { background: #ecf5ff; border-color: #409eff; }
.pagination button:disabled { opacity: 0.5; cursor: not-allowed; }

.modal-mask { position: fixed; z-index: 9999; top: 0; left: 0; width: 100%; height: 100%; background-color: rgba(0,0,0,0.6); display: flex; justify-content: center; align-items: center; overflow-y: auto; }
.modal-wrapper { width: 600px; max-height: 80%; background: #fff; border-radius: 12px; padding: 20px; overflow-y: auto; box-shadow: 0 6px 20px rgba(0,0,0,0.3); position: relative; }
.modal-container h3 { margin-bottom: 12px; }
.modal-close { position: absolute; top: 15px; right: 15px; font-size: 20px; background: transparent; border: none; cursor: pointer; color: #999; transition: color 0.2s; }
.modal-close:hover { color: #333; }
.label { font-weight: 500; color: #555; }
.labels { display: flex; flex-wrap: wrap; gap: 6px; }
.label-badge { background: #eef6ff; color: #3178c6; padding: 4px 10px; border-radius: 12px; font-size: 13px; }
.annotations { background: #f9f9f9; padding: 10px; border-radius: 6px; white-space: pre-wrap; font-size: 13px; }
</style>
